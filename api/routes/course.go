package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/UTDNebula/nebula-api/api/dao"
)

func CourseRoute(router *gin.Engine, client *mongo.Client) {
	// All routes related to courses come here
	courseGroup := router.Group("/course")

	dao := dao.NewCourseDao(configs.GetCollection(client, "courses"), 1024)
	api := controllers.NewCourseAPI(dao)
	courseGroup.OPTIONS("", controllers.Preflight)
	courseGroup.GET("", api.Search)
	courseGroup.GET(":id", api.ById)
}
