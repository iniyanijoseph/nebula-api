package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/dao/course"
	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type CourseAPI struct {
	dao course.Dao
}

func NewCourseAPI(dao course.Dao) *CourseAPI {
	return &CourseAPI{dao: dao}
}

func (api *CourseAPI) Search(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter, err := course.NewFilterFromValues(c.Request.URL.Query())
	if err != nil {
		if _, ok := errors.Cause(err).(schema.UnknownKeyError); ok {
			c.JSON(
				http.StatusBadRequest,
				responses.CourseResponse{Message: err.Error()},
			)
		} else {
			c.JSON(
				http.StatusInternalServerError,
				responses.CourseResponse{Message: err.Error()},
			)
		}
	}

	courses, err := api.dao.Filter(ctx, filter)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.CourseResponse{Message: "Internal Server Error: " + err.Error()},
		)
	}

	c.JSON(http.StatusOK, responses.CourseResponse{
		Message: "success",
		Data:    courses,
	})
}

func (api *CourseAPI) ById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	courseId := c.Param("id")
	course, err := api.dao.FindById(ctx, courseId)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.CourseResponse{Message: "Internal Server Error: " + err.Error()},
		)
	}

	c.JSON(http.StatusOK, course)
}
