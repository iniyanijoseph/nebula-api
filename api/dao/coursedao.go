package dao

import (
	"context"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Course map[string]interface{}

// CourseFilter represents the filter parameters for the MongoDB query.
type CourseFilter struct {
	CourseNumber           string `bson:"course_number,omitempty" schema:"course_number,omitempty"`
	SubjectPrefix          string `bson:"subject_prefix,omitempty" schema:"subject_prefix,omitempty"`
	School                 string `bson:"school,omitempty" schema:"school,omitempty"`
	ClassLevel             string `bson:"class_level,omitempty" schema:"class_level,omitempty"`
	CreditHours            string `bson:"credit_hours,omitempty" schema:"credit_hours,omitempty"`
	ActivityType           string `bson:"activity_type,omitempty" schema:"activity_type,omitempty"`
	Grading                string `bson:"grading,omitempty" schema:"grading,omitempty"`
	InternalCourseNumber   string `bson:"internal_course_number,omitempty" schema:"internal_course_number,omitempty"`
	LectureContactHours    string `bson:"lecture_contact_hours,omitempty" schema:"lecture_contact_hours,omitempty"`
	LaboratoryContactHours string `bson:"laboratory_contact_hours,omitempty" schema:"laboratory_contact_hours,omitempty"`
	OfferingFrequency      string `bson:"offering_frequency,omitempty" schema:"offering_frequency,omitempty"`
	Offset                 int64  `bson:"-" schema:"offset,omitempty"`
}

func NewCourseFilterFromValues(m map[string][]string) (*CourseFilter, error) {
	var filter CourseFilter
	err := schema.NewDecoder().Decode(&filter, m)
	return &filter, errors.Wrap(err, "dao.course.CourseFilter: could not decode struct")
}

func (filter *CourseFilter) ToDocument() ([]byte, error) {
	b, err := bson.Marshal(filter)
	return b, errors.Wrap(err, "dao.course.CourseFilter: could not marshal to document")
}

func (filter *CourseFilter) GetOffset() int64 {
	return filter.Offset
}

type CourseDao interface {
	Filter(ctx context.Context, filter *CourseFilter) ([]Course, error)
	FindById(ctx context.Context, objId string) (*Course, error)
}

type courseDaoImpl struct {
	helper *collectionHelper[*CourseFilter, Course]
}

func NewCourseDao(coll *mongo.Collection, pageLimit int64) CourseDao {
	return &courseDaoImpl{helper: newCollectionHelper[*CourseFilter, Course](coll, pageLimit)}
}

func (dao *courseDaoImpl) Filter(ctx context.Context, filter *CourseFilter) ([]Course, error) {
	return dao.helper.Filter(ctx, filter)
}

func (dao *courseDaoImpl) FindById(ctx context.Context, id string) (*Course, error) {
	return dao.helper.FindById(ctx, id)
}
