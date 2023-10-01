package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/dao"
	"github.com/UTDNebula/nebula-api/api/responses"

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

	filter := dao.CourseFilter{
		CourseNumber:         c.Query("course_number"),
		SubjectPrefix:        c.Query("subject_prefix"),
		Title:                c.Query("title"),
		Description:          c.Query("description"),
		School:               c.Query("school"),
		CreditHours:          c.Query("credit_hours"),
		ActivityType:         c.Query("activity_type"),
		Grading:              c.Query("grading"),
		InternalCourseNumber: c.Query("internal_course_number"),
		LectureContactHours:  c.Query("lecture_contact_hours"),
		OfferingFrequency:    c.Query("offering_frequency"),
	}

	courses, err := api.dao.Filter(ctx, &filter)
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
