package middleware

import (
	"github.com/Iliad/vacancytest/pkg/db"
	"github.com/gin-gonic/gin"
)

const (
	DBContext = "database"
)

func RegisterDB(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(DBContext, *db)
	}
}
