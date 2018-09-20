package handlers

import (
	"net/http"
	"strconv"

	"github.com/Iliad/vacancytest/pkg/models"
	"github.com/gin-gonic/gin/binding"

	"github.com/Iliad/vacancytest/pkg/db"
	"github.com/Iliad/vacancytest/pkg/router/middleware"
	"github.com/gin-gonic/gin"
)

func GetVacancyHandler(ctx *gin.Context) {
	db := ctx.MustGet(middleware.DBContext).(db.DB)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := db.GetVacancy(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func GetVacanciesHandler(ctx *gin.Context) {
	db := ctx.MustGet(middleware.DBContext).(db.DB)

	resp, err := db.GetVacancies(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func AddVacancyHandler(ctx *gin.Context) {
	db := ctx.MustGet(middleware.DBContext).(db.DB)

	var request models.Vacancy
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := db.AddVacancy(ctx.Request.Context(), &request); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, request)
}

func DeleteVacancyHandler(ctx *gin.Context) {
	db := ctx.MustGet(middleware.DBContext).(db.DB)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := db.DeleteVacancy(ctx.Request.Context(), id); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusAccepted)
}
