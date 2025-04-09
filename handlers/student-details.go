package handlers

import (
	"CollegeAdministration/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) InsertStudentDetails(ctx *gin.Context) {

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

	var cad models.StudentInfo
	err := ctx.BindJSON(&cad)
	if err != nil {
		log.Println("not able to store values")
	}
	account_id, err := h.service.GetAccountByToken(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	response := h.service.InsertValuesToCAd(&cad, account_id)
	if response != nil {
		ctx.JSON(http.StatusInternalServerError, response.Error())
		return
	} else {
		ctx.JSON(http.StatusCreated, "successfully inserted to table")
	}
}

func (h *Handler) RetrieveValuesForStudent(ctx *gin.Context) {

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

	order_by := ctx.Query("order")

	if order_by == "" {
		response, err := h.service.Retrieve_student_details()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, response)
		}
	} else {
		response, err := h.service.Retrieve_student_detailsbyOrder(order_by)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, response)
		}
	}

}

func (h *Handler) UpdateValuesForStudent(ctx *gin.Context) {

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

	oldCourse := ctx.Param("coursename")
	oldName := ctx.Param("student_name")
	oldRollNumber := ctx.Param("roll_number")

	if oldCourse == ":coursename" || oldName == ":student_name" || oldRollNumber == ":roll_number" {
		ctx.JSON(http.StatusBadRequest, "Empty Inputs from Client")
		return
	}
	var rcd models.StudentInfo
	ctx.BindJSON(&rcd)
	err := h.service.Update_Student_Details(&rcd, oldCourse, oldName, oldRollNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error.Error(err))
		return
	} else {
		ctx.JSON(http.StatusOK, "Success")
	}

}

func (h *Handler) DeleteStudentDetails(ctx *gin.Context) {

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

	ID := ctx.Param("id")
	if ID == ":id" {
		ctx.JSON(http.StatusBadRequest, "Empty Inputs from Client")
		return
	}

	IDUUID, errParsingUUID := uuid.Parse(ID)
	if errParsingUUID != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprint("error in parsing uuid -", errParsingUUID.Error()))
		return
	}

	err := h.service.DeleteStudent(IDUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error.Error(err))
	} else {
		ctx.JSON(http.StatusOK, "Success")
	}

}
func (h *Handler) UpdateStudentNameAndAge(ctx *gin.Context) {

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

	var rcd *models.StudentInfo
	existing_name := ctx.Param("name")
	ctx.BindJSON(&rcd)
	err := h.service.UpdateStudentNameAge(existing_name, rcd.Name, rcd.Age)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error.Error(err))
	} else {
		ctx.JSON(http.StatusOK, "Success")
	}

}

func (h *Handler) FetchAllCourseForAStudent(ctx *gin.Context) {

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
	student_name := ctx.Param("name")
	res, err := h.service.FetchStudentCourse(student_name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error.Error(err))
	} else {
		ctx.JSON(http.StatusOK, res)
	}

}

func (h *Handler) DeleteStudentCourse(ctx *gin.Context) {

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

	parameter := ctx.Params
	student_name := parameter.ByName("name")
	course_name := parameter.ByName("course")
	if student_name == ":name" || course_name == ":course" {
		ctx.IndentedJSON(http.StatusInternalServerError, "parameter is empty")
		return
	}
	err := h.service.DeleteStudentCourseService(student_name, course_name)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	} else {
		ctx.IndentedJSON(http.StatusOK, "deleted")
	}
}

func (h *Handler) GetRankingForACourse(ctx *gin.Context) {

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
	course_name := ctx.Param("coursename")
	model, err := h.service.GetAllStudentsMarksForGivenCourse(course_name)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
	} else {
		ctx.IndentedJSON(http.StatusOK, model)
	}

}

func (h *Handler) GetSelectedFieldsAllStudent(ctx *gin.Context) {
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

	response, err4 := h.service.GetAllStudentSelectiveData()
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, err4.Error())
		return
	}
	ctx.IndentedJSON(http.StatusAccepted, response)

}
func (h *Handler) DeleteStudentWithSpecifics(ctx *gin.Context) {

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

	st_req := &models.StudentInfo{}
	err := ctx.BindJSON(st_req)
	if err != nil {
		err = fmt.Errorf("unable to store values")
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err3 := h.service.DeleteStudentSpecifics(st_req)
	if err3 != nil {
		ctx.JSON(http.StatusInternalServerError, err3.Error())
		return
	}
	ctx.IndentedJSON(http.StatusAccepted, "Success")

}

func (h *Handler) UpdateValuesForStudentV2(ctx *gin.Context) {

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

	var rcd models.StudentInfo
	ctx.BindJSON(&rcd)
	err := h.service.UpdateStudentDetailsV2(&rcd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error.Error(err))
		return
	} else {
		ctx.JSON(http.StatusOK, "Success")
	}

}
