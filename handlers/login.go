package handlers

import (
	"CollegeAdministration/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) login(ctx *gin.Context) {

	req := &models.InstructorLogin{}

	err := ctx.BindJSON(req)
	if err != nil {
		ctx.JSON(http.StatusMethodNotAllowed, err)
		return
	}
	err = h.service.ValidateLogin(req.EmailId, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err1 := h.service.CheckCredentials(req.EmailId, req.Password)
	if err1 != nil {
		ctx.JSON(http.StatusInternalServerError, err1.Error())
		return
	}

	instructor_id, err3 := h.service.GetInstructorIDWithEmail(req.EmailId)
	if err3 != nil {
		ctx.JSON(http.StatusInternalServerError, err3.Error())
		return
	}
	account_status, err := h.service.VerifyAccountStatusById(instructor_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !account_status {
		ctx.JSON(http.StatusNotAcceptable, "Account is inactive! Please Verify your account")
		return
	}

	token, err2 := h.service.GetTokenAfterLogging(instructor_id)
	if err2 != nil {
		ctx.JSON(http.StatusNotAcceptable, err2.Error())
		return
	}

	ctx.Writer.Header().Set("token", token.String())
	ctx.Writer.Header().Set("account_id", instructor_id)
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Value:    token.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "account_id",
		Value:    instructor_id,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})
	ctx.JSON(http.StatusOK, "successfully logged-in")

}

func (h *Handler) logout(ctx *gin.Context) {

	token := ctx.Query("token")
	err := h.service.DisableToken(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

}

func (h *Handler) SendOTPEmail(ctx *gin.Context) {
	emailId := ctx.Query("email_id")
	err := h.service.GenerateOTPAndStore(emailId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "OTP email sent successfully")
}

func (h *Handler) VerifyOTP(ctx *gin.Context) {
	email := ctx.Query("email_id")
	otp := ctx.Query("otp")

	err := h.service.VerifyAccountWithOTP(email, otp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "OTP verified successfully")
}
