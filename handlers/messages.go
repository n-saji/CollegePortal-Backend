package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateMessageStatusAsRead(c *gin.Context) {

	token := c.GetHeader("Token")
	var err1 error
	if token == "" {
		token, err1 = c.Cookie("token")
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprint("no token found -", err1.Error()))
			return
		}
	}

	err2 := h.service.CheckTokenWithCookie(token)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, err2.Error())
		return
	}

	account_id := c.Param("id")
	account_id, err := h.service.UpdateMessageStatusAsRead(account_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"account_id": account_id,
	})
}

func (h *Handler) SendMessage(c *gin.Context) {

	h.service.SendMessageAsBroadCast(c.Query("type"), c.Query("msg"))
	c.JSON(http.StatusOK, gin.H{
		"message": "sent",
	})
}
