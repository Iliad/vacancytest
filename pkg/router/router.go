package router

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Iliad/vacancytest/pkg/router/handlers"
	"github.com/gin-gonic/contrib/ginrus"

	"github.com/Iliad/vacancytest/pkg/db"
	"github.com/Iliad/vacancytest/pkg/router/middleware"

	"github.com/gin-gonic/gin"
)

func CreateRouter(db *db.DB) http.Handler {
	e := gin.Default()
	e.Use(gin.Recovery())
	e.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))
	e.Use(middleware.RegisterDB(db))
	initRoutes(e)
	return e
}

func initRoutes(app *gin.Engine) {
	users := app.Group("/users")
	{
		users.POST("/register", handlers.CreateUserHandler)
		users.POST("/role", middleware.CheckAuth(true), handlers.SetUserRoleHandler)
	}
	vacancies := app.Group("/vacancies")
	{
		vacancies.GET("/:id", middleware.CheckAuth(false), handlers.GetVacancyHandler)
		vacancies.GET("", middleware.CheckAuth(false), handlers.GetVacanciesHandler)
		vacancies.POST("", middleware.CheckAuth(true), handlers.AddVacancyHandler)
		vacancies.DELETE("/:id", middleware.CheckAuth(true), handlers.DeleteVacancyHandler)
	}
}
