package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	docs "github.com/fahminlb33/devoria1-wtc-backend/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/fahminlb33/devoria1-wtc-backend/domain/articles"
	"github.com/fahminlb33/devoria1-wtc-backend/domain/users"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/authentication"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/config"
	"github.com/fahminlb33/devoria1-wtc-backend/infrastructure/database"
)

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
	// initialize services
	config.LoadConfig()
	authentication.InitializeJwtAuth()

	// open DB conection
	db, err := database.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	// database initialization
	err = database.MigrateIfNeeded(db)
	if err != nil {
		log.Fatal(err)
	}

	err = database.SeedIfNeeded(db)
	if err != nil {
		log.Fatal(err)
	}

	// create new router
	router := gin.New()

	router.Use(gin.Logger())
	//router.Use(gin.Recovery())

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// register custom validator
		v.RegisterValidation("isarticlepublishstatus", articles.IsArticlePublishStatus)
	}

	// cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// enable apm
	//router.Use(apmgin.Middleware(router))

	// swagger
	docs.SwaggerInfo.BasePath = "/api"
	router.GET("", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// users API
	userUsecase := users.ConstructUserUseCase(db)
	users.ConstructUserHandler(router, userUsecase)

	// articles API
	articleUsecase := articles.ConstructArticlesUseCase(db)
	articles.ConstructArticlesHandler(router, articleUsecase)

	// create server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.GlobalConfig.Server.Host, config.GlobalConfig.Server.Port),
		Handler: router,
	}

	// bootstrap app
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	// --- Graceful shutdown
	// create channel to allow graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// wait for termination signal
	<-quit
	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// gracefully shutdown app
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown:", err)
	}

	log.Println("Server exiting")
}
