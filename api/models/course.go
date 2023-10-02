package models

// This wrapper is necessary to properly represent the ObjectID structure when serializing to JSON instead of BSON
/*type string struct {
	Id primitive.ObjectID `json:"$oid"`
}*/

type Course struct {
	Id                       string                 `json:"_id"`
	Subject_prefix           string                 `json:"subject_prefix"`
	Course_number            string                 `json:"course_number"`
	Title                    string                 `json:"title"`
	Description              string                 `json:"description"`
	Enrollment_reqs          string                 `json:"enrollment_reqs"`
	School                   string                 `json:"school"`
	Credit_hours             string                 `json:"credit_hours"`
	Class_level              string                 `json:"class_level"`
	Activity_type            string                 `json:"activity_type"`
	Grading                  string                 `json:"grading"`
	Internal_course_number   string                 `json:"internal_course_number"`
	Prerequisites            *CollectionRequirement `json:"prerequisites"`
	Corequisites             *CollectionRequirement `json:"corequisites"`
	Co_or_pre_requisites     *CollectionRequirement `json:"co_or_pre_requisites"`
	Sections                 []string               `json:"sections"`
	Lecture_contact_hours    string                 `json:"lecture_contact_hours"`
	Laboratory_contact_hours string                 `json:"laboratory_contact_hours"`
	Offering_frequency       string                 `json:"offering_frequency"`
	Catalog_year             string                 `json:"catalog_year"`
	Attributes               interface{}            `json:"attributes"`
}
