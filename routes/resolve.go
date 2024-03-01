package routes

import (
	"fmt"

	"github.com/Makinde1034/url-shortner/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveUrl(c *fiber.Ctx ) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	fmt.Println(url,"url")

	value, err := r.Get(database.Ctx,url).Result()

	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error" : "Key not found in the database."})
	}


	return c.Redirect(value,301)
}