package routes

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Makinde1034/url-shortner/database"
	"github.com/Makinde1034/url-shortner/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL  string `json:"url"`
	CustomShort string `json:"customShort"`
	Expiry time.Duration
}

type response struct {
	URL string `json:"url"`
	CustomShort string `json:"short"`
	Expiry time.Duration `json:"expiry"`
	XRateRemaining int `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_rest"`                     
}

func ShortenUrl(c *fiber.Ctx) error {
	body := new (request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"error parsing request"})
		
	}


	// implement rate limiting
	r2 := database.CreateClient(1)
	defer  r2.Close()

	value, err := r2.Get(database.Ctx,c.IP()).Result()
	fmt.Println(err,"valeuee")
	fmt.Println(c.IP(),"valeuee")

	if err == redis.Nil {
    // sets new value or updates  in redis DB
		r2.Set(database.Ctx,c.IP(),os.Getenv("API_QUOTA"),1800 * time.Second)
	}else{
		valInt, _ := strconv.Atoi(value)

		if valInt <= 0 {
			//  Returns the remaining time to live of a key that has a timeout
			limit, _ := r2.TTL(database.Ctx,c.IP()).Result()
			fmt.Println(valInt)
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error" : fmt.Sprintf("Rate limit exceeded. Try again in %d minutes",limit / time.Nanosecond / time.Minute),
				"rate_limit" : limit / time.Nanosecond / time.Minute,

			})
		}
	}


	if !govalidator.IsURL(body.URL){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid url"})
	}

	if !helpers.RemoveDomainError(body.URL){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error" : "error"})
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	

	r := database.CreateClient(0)
	defer r.Close()

	customShortExists, _ := r.Get( database.Ctx, body.CustomShort).Result()

	if customShortExists != "" { 
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error" : "Custom url already exists"})
	}

	saveError := r.Set(database.Ctx,body.CustomShort,body.URL,24*36000*time.Second).Err()

	if saveError != nil {
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error" : "Unable to complete request"})
	}




	response := response{
		URL:  body.URL,
		CustomShort: body.CustomShort,
		Expiry:  body.Expiry,
		XRateRemaining: 10,
		XRateLimitReset: 30,
	}

  r2.Decr(database.Ctx,c.IP())

  val, _ := r2.Get(database.Ctx,c.IP()).Result()

  response.XRateRemaining,_ = strconv.Atoi(val)

  ttl, _ := r2.TTL(database.Ctx,c.IP()).Result()
  response.XRateLimitReset = ttl / time.Nanosecond /time.Minute

  response.CustomShort = os.Getenv("DOMAIN") + "/" + body.CustomShort

  return c.Status(fiber.StatusOK).JSON(response)

}