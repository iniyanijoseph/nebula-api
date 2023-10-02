package course

import (
	"context"

	"github.com/UTDNebula/nebula-api/api/dao"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type Course map[string]interface{}

// Filter represents the filter parameters for the MongoDB query.
type Filter struct {
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
	Offset                 int64  `schema:"offset,omitempty"`
}

func NewFilterFromValues(m map[string][]string) (*Filter, error) {
	var filter Filter
	err := schema.NewDecoder().Decode(&filter, m)
	return &filter, errors.Wrap(err, "dao.course.Filter: could not decode struct")
}

func (filter *Filter) ToDocument() ([]byte, error) {
	b, err := bson.Marshal(filter)
	return b, errors.Wrap(err, "dao.course.Filter: could not marshal to document")
}

func (filter *Filter) GetOffset() int64 {
	return filter.Offset
}

type Dao interface {
	Filter(ctx context.Context, filter *Filter) ([]Course, error)
	FindById(ctx context.Context, objId string) (Course, error)
}

type daoImpl dao.CollectionHelper[*Filter, Course]

func (dao *daoImpl) Filter(ctx context.Context, filter *Filter) ([]Course, error) {
	return dao.Filter(ctx, filter)
}

func (dao *daoImpl) FindById(ctx context.Context, id string) (*Course, error) {
	return dao.FindById(ctx, id)
}
