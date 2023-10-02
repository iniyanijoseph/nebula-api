package controllers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/UTDNebula/nebula-api/api/dao"
	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type CourseAPI struct {
	dao dao.CourseDao
}

func NewCourseAPI(dao dao.CourseDao) *CourseAPI {
	return &CourseAPI{dao: dao}
}

func (api *CourseAPI) Search(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter, err := dao.NewCourseFilterFromValues(c.Request.URL.Query())
	if err != nil {
		fmt.Println(reflect.TypeOf(errors.Cause(err)))
		if dao.IsUnknownFieldErr(err) {
			c.JSON(
				http.StatusBadRequest,
				responses.CourseResponse{Message: err.Error()},
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			responses.CourseResponse{Message: err.Error()},
		)
		return
	}

	courses, err := api.dao.Filter(ctx, filter)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.CourseResponse{Message: "Internal Server Error: " + err.Error()},
		)
		return
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
