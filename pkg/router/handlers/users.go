package handlers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/Iliad/vacancytest/pkg/db"
	"github.com/Iliad/vacancytest/pkg/models"
	"github.com/Iliad/vacancytest/pkg/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetUserHandler(ctx *gin.Context) {
	db := ctx.MustGet(middleware.DBContext).(db.DB)
	user, err := db.GetUser(ctx.Request.Context(), ctx.Param("login"))
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	user.Password = ""
	ctx.JSON(http.StatusOK, user)
}

func GetUsersHandler(ctx *gin.Context) {
	db := ctx.MustGet(middleware.DBContext).(db.DB)
	users, err := db.GetUsers(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	ctx.JSON(http.StatusOK, users)
}

func CreateUserHandler(ctx *gin.Context) {
	db := ctx.MustGet(middleware.DBContext).(db.DB)

	var request models.User
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 8)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	request.Password = string(hashedPassword)

	usersCount, err := db.GetUsersCount(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	//First user would be created with editor permissions.
	if *usersCount == 0 {
		request.Role = models.Editor
	} else {
		request.Role = models.Viewer
	}

	if err := db.CreateUser(ctx.Request.Context(), &request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	request.Password = ""
	ctx.JSON(http.StatusAccepted, request)
}

func SetUserRoleHandler(ctx *gin.Context) {
	db := ctx.MustGet(middleware.DBContext).(db.DB)

	var request models.ChangeUserRole
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	authLogin, _, _ := ctx.Request.BasicAuth()
	if request.Login == authLogin {
		ctx.AbortWithStatusJSON(http.StatusForbidden, "unable to change own role")
		return
	}

	if err := db.ChangeUserRole(ctx.Request.Context(), request.Login, request.NewRole); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusAccepted)
}
