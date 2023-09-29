package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseAPI struct {
	coll *mongo.Collection
}

func NewCourseAPI(client *mongo.Client) *CourseAPI {
	return &CourseAPI{
		coll: configs.GetCollection(client, "courses"),
	}
}
func (api *CourseAPI) CourseSearch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

	// @TODO: Fix with model - There is NO typechecking!
	// var courses []models.Course
	var courses []map[string]interface{}

	// build query key value pairs (only one value per key)
	query := bson.M{}
	for key := range queryParams {
		query[key] = c.Query(key)
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		c.JSON(http.StatusConflict, responses.CourseResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	// get cursor for query results
	cursor, err := api.coll.Find(ctx, query, optionLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.CourseResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &courses); err != nil {
		panic(err)
	}

	// return result
	c.JSON(http.StatusOK, responses.CourseResponse{Status: http.StatusOK, Message: "success", Data: courses})
}

func (api *CourseAPI) CourseById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := c.Param("id")

	// @TODO: Fix with model - There is NO typechecking!
	// var course models.Course
	var course map[string]interface{}

	// parse object id from id parameter
	objId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.CourseResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
		return
	}

	// find and parse matching course
	err = api.coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.CourseResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// return result
	c.JSON(http.StatusOK, responses.CourseResponse{Status: http.StatusOK, Message: "success", Data: course})
}
