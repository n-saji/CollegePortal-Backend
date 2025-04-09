package handlers

import (
	"CollegeAdministration/models"
	"CollegeAdministration/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) InsertCourse(ctx *gin.Context) {

	token := ctx.GetHeader("Token")
	var err1 error
	if token == "" {
		token, err1 = ctx.Cookie("token")
		if err1 != nil {
			ctx.JSON(http.StatusInternalServerError, fmt.Sprint("no token found -", err1.Error()))
			return
		}
	}

	err2 := h.service.CheckTokenWithCookie(token)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, err2.Error())
		return
	}

	account_id,err := h.service.GetAccountByToken(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var ca models.CourseInfo
	err = ctx.BindJSON(&ca)

	if err != nil {
		log.Println("not able to store values")
	}
	response := h.service.InsertValuesToCA(account_id,&ca)
	if response != nil {
		ctx.JSON(http.StatusInternalServerError, response.Error())
		return
	} else {
		ctx.JSON(http.StatusCreated, "successfully inserted to table")
	}

}

func (h *Handler) RetrieveValuesCourse(ctx *gin.Context) {

	ok := ctx.GetHeader("Internal-call")

	if ok != "true" {
		token := ctx.GetHeader("Token")
		var err1 error
		if token == "" {
			token, err1 = ctx.Cookie("token")
			if err1 != nil {
				ctx.JSON(http.StatusInternalServerError, fmt.Sprint("no token found -", err1.Error()))
				return
			}
		}

		err2 := h.service.CheckTokenWithCookie(token)
		if err2 != nil {
			ctx.JSON(http.StatusInternalServerError, err2.Error())
			return
		}
	}

	response, err := h.service.RetrieveCA()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	} else {
		ctx.JSON(http.StatusOK, response)
	}

}
func (h *Handler) UpdateValuesCourse(ctx *gin.Context) {

	token := ctx.GetHeader("Token")
	var err1 error
	if token == "" {
		token, err1 = ctx.Cookie("token")
		if err1 != nil {
			ctx.JSON(http.StatusInternalServerError, fmt.Sprint("no token found -", err1.Error()))
			return
		}
	}

	err2 := h.service.CheckTokenWithCookie(token)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, err2.Error())
		return
	}
	var rc models.CourseInfo
	var name = ctx.Param("name")
	if name == ":name" {
		ctx.JSON(http.StatusInternalServerError, "new course can't be empty")
		return
	}
	ctx.BindJSON(&rc)
	if name == rc.CourseName {
		ctx.JSON(http.StatusOK, "no changes")
		return
	}
	err := h.service.UpdateCA(name, &rc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error.Error(err))
		return
	} else {
		ctx.JSON(http.StatusOK, "successfully updated")
	}

}

func (h *Handler) DeleteCourse(ctx *gin.Context) {

	token := ctx.GetHeader("Token")
	var err1 error
	if token == "" {
		token, err1 = ctx.Cookie("token")
		if err1 != nil {
			ctx.JSON(http.StatusInternalServerError, fmt.Sprint("no token found -", err1.Error()))
			return
		}
	}

	err2 := h.service.CheckTokenWithCookie(token)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, err2.Error())
		return
	}

	var CourseName string = ctx.Param("courseName")
	err := h.service.DeleteCA(CourseName)
	if err != nil {
		res := models.DeleteResponse{}
		_, err1 := utils.MakeRequest(http.MethodGet, "https://dolphin-app-2zya2.ondigitalocean.app/retrieve-all-courses", "Fetching course", nil, &res.Courses)
		if err1 != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err1)
			return
		}
		res.Message = err.Error() + "!"
		ctx.IndentedJSON(http.StatusInternalServerError, res)
	} else {
		ctx.JSON(http.StatusOK, "successfully deleted")
	}
}
