package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.peaut.limit/combination"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		combination.Demo()
		return c.SendString("Hello World!")
	})
	/*app.Get("/hid", func(c *fiber.Ctx) error {
		return c.SendStream()
	})*/
	app.Listen(":3000")
}

type worker interface {
	work()
}

type person struct {
	name string
	worker
}

func demo() {
	var work worker = person{}
	fmt.Println(work)
}
