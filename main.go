package main

import (
	"fmt"
	"log"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.elastic.co/apm/module/apmfiber"
	"go.elastic.co/apm/module/apmgorm"
	_ "go.elastic.co/apm/module/apmgorm/dialects/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
			"john":  "doe",
			"admin": "123456",
		},
	}))

	// swagger
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:          "http://example.com/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
	}))

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

	log.Fatal(app.Listen(":9000"))
}

// login godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts/{id} [get]
func login() {

}
