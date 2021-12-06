package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"go.elastic.co/apm/module/apmfiber"
	"go.elastic.co/apm/module/apmgorm"
	_ "go.elastic.co/apm/module/apmgorm/dialects/postgres"
	"gorm.io/gorm"

	_ "github.com/fahminlb33/devoria1-wtc-backend/docs"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/util"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// @title MEWS API
// @version 1.0
// @description MEWS API for Devoria's WTC
// @termsOfService http://swagger.io/terms/
// @contact.name Fahmi Noor Fiqri
// @license.name MIT License
// @license.url http://www.opensource.org/licenses/MIT
// @host :9000
// @BasePath /
func main() {
	util.LoadConfig()

	app := fiber.New()

	// enable apm
	app.Use(apmfiber.Middleware())

	// cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	// basic auth
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			util.GlobalConfig.Authentication.BasicUsername: util.GlobalConfig.Authentication.BasicPassword,
		},
	}))

	// swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// GET /api/register
	app.Get("/api/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})

	// GET /flights/LAX-SFO
	app.Get("/flights/:from-:to", func(c *fiber.Ctx) error {
		db, err := apmgorm.Open("postgres", "")
		if err != nil {
			log.Fatal(err)
		}

		db = apmgorm.WithContext(c.Context(), db)
		msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
		return c.SendString(msg) // => 💸 From: LAX, To: SFO
	})

	// GET /dictionary.txt
	app.Get("/:file.:ext", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("📃 %s.%s", c.Params("file"), c.Params("ext"))
		return c.SendString(msg) // => 📃 dictionary.txt
	})

	// GET /john/75
	app.Get("/:name/:age/:gender?", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("👴 %s is %s years old", c.Params("name"), c.Params("age"))
		return c.SendString(msg) // => 👴 john is 75 years old
	})

	// GET /john
	app.Get("/:name", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", util.GlobalConfig.Server.Host, util.GlobalConfig.Server.Port)))
}

// login godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Router       /accounts/{id} [get]
func login() {

}
