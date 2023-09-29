package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func CourseRoute(router *gin.Engine, client *mongo.Client) {
	// All routes related to courses come here
	courseGroup := router.Group("/course")
	api := controllers.NewCourseAPI(client)
	courseGroup.OPTIONS("", controllers.Preflight)
	courseGroup.GET("", api.CourseSearch)
	courseGroup.GET(":id", api.CourseById)
}
