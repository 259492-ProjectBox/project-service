package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/project-box/configs"
	"github.com/project-box/handlers"
	"github.com/project-box/routers"
)

func NewApp(
	projectHandler handlers.ProjectHandler,
	calendarHandler handlers.CalendarHandler,
	resourceHandler handlers.ResourceHandler,
	employeeHandler handlers.EmployeeHandler,
	configHanler handlers.ConfigHandler,
	projectconfigHandler handlers.ProjectConfigHandler,
) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowCredentials: true,
		}),
	)

	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	routers.SetupRoutes(
		r,
		projectHandler,
		resourceHandler,
		calendarHandler,
		employeeHandler,
		configHanler,
		projectconfigHandler,
	)

	return r, nil
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @SecurityDefinitions.apikey BearerAuth
// @In header
// @Name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	configs.InitialEnv(".env")
	app, cleanup, err := InitializeApp()
	if err != nil {
		log.Print(err)
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-exitChan
		log.Print("Shutting down the server...")
		cleanup()
		os.Exit(0)
	}()

	if err := app.Run(fmt.Sprintf(":%s", configs.GetPort())); err != nil {
		log.Print(err)
	}

}
