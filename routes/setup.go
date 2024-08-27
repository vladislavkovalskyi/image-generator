package routes

import (
	"image-generator/routes/post"

	"github.com/gin-gonic/gin"
)


func SetupRoutes(router *gin.Engine) {
	router.POST("/generate_text", post.GenerateText)
}