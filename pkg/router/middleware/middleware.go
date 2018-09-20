package middleware

import (
	"net/http"

	"github.com/Iliad/vacancytest/pkg/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/Iliad/vacancytest/pkg/db"
	"github.com/gin-gonic/gin"
)

func CheckAuth(requireEditorRole bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login, password, _ := ctx.Request.BasicAuth()

		if login == "" || password == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "authorization required")
			return
		}

		db := ctx.MustGet(DBContext).(db.DB)
		user, err := db.GetUser(ctx.Request.Context(), login)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusForbidden, "invalid login or password")
			return
		}
		if user == nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "invalid login or password")
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusForbidden, "invalid login or password")
			return
		}
		if requireEditorRole {
			if user.Role != models.Editor {
				ctx.AbortWithStatusJSON(http.StatusForbidden, "editor permissions required")
				return
			}
		}
	}
}
