package course

import (
	"context"
	"fmt"

	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func NewFilterFromQueryParams(m map[string][]string) (*Filter, error) {
	var filter Filter
	err := schema.NewDecoder().Decode(&filter, m)
	return &filter, errors.Wrap(err, "dao.course.Filter: could not decode struct")
}

func (filter *Filter) toDocument() ([]byte, error) {
	b, err := bson.Marshal(filter)
	return b, errors.Wrap(err, "dao.course.Filter: could not marshal to document")
}

type Dao interface {
	Filter(ctx context.Context, filter *Filter) ([]Course, error)
	FindById(ctx context.Context, objId string) (Course, error)
}

type daoImpl struct {
	coll     *mongo.Collection
	pageSize int64
}

func NewDao(client *mongo.Client, courseCollName string, pageSize int64) Dao {
	return &daoImpl{
		coll: configs.GetCollection(client, courseCollName),
	}
}

func (dao *daoImpl) Filter(ctx context.Context, filter *Filter) ([]Course, error) {
	query, err := filter.toDocument()
	if err != nil {
		return nil, err
	}

	limitOpt := options.Find.SetSkip(filter.Offset).SetLimit(dao.pageSize)
	cursor, err := dao.coll.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query courses: %v", err)
	}

	var courses []Course
	err = cursor.All(ctx, &courses)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal courses from query: %v", err)
	}

	return courses, nil
}

func (dao *daoImpl) FindById(ctx context.Context, courseId string) (Course, error) {
	objId, err := primitive.ObjectIDFromHex(courseId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse object id: %v", courseId)
	}

	var course Course
	err = dao.coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&course)
	if err != nil {
		return nil, fmt.Errorf("failed to find and decode course by id=%v : %v", objId, err)
	}

	return course, nil
}
