package handlers

import (
	"net/http"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/pdf"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/generate", func(c *gin.Context) {
		var payload map[string]string
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pdfPath, err := pdf.FillForm("templates/delegacny_list_ligy.pdf", payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.FileAttachment(pdfPath, "delegacny.pdf")
	})

	r.GET("/list-fields", func(c *gin.Context) {
		err := pdf.ListFillableFields("templates/delegacny_list_ligy.pdf")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Fields listed to console. Check server logs for details."})
	})
}
