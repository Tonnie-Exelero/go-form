package models

import "github.com/shurcooL/graphql"

// Course represents the structure of the course returned by the GraphQL query.
type Course struct {
	DocumentId       graphql.String `json:"documentId"`
	Title            graphql.String `json:"title"`
	Slug             graphql.String `json:"slug"`
	LearningOutcomes graphql.String `json:"learningOutcomes"`
	Description      graphql.String `json:"description"`
	CourseCode       graphql.String `json:"courseCode"`
	Instructor       graphql.String `json:"instructor"`
	Duration         graphql.String `json:"duration"`
	StartDate        graphql.String `json:"startDate"`
	Price            Price          `json:"price"`
	Mode             Mode           `json:"mode"`
	Categories       []Category     `json:"categories"`
	Provider         Provider       `json:"provider"`
	Features         Features       `json:"features"`
}

// Price holds the pricing details for the course.
type Price struct {
	CourseFee           graphql.String  `json:"course_fee"`
	PaymentOptions      graphql.String  `json:"payment_options"`
	CertificateProvided graphql.Boolean `json:"certificateProvided"`
}

// Mode indicates the delivery methods for the course.
type Mode struct {
	Online    graphql.Boolean `json:"online"`
	HandsOn   graphql.Boolean `json:"handsOn"`
	SelfPaced graphql.Boolean `json:"selfPaced"`
	Lecture   graphql.Boolean `json:"lecture"`
	Classroom graphql.Boolean `json:"classroom"`
}

// Category represents a course category.
type Category struct {
	Name graphql.String `json:"name"`
}

// Provider holds details about the course provider.
type Provider struct {
	CompanyName  graphql.String `json:"companyName"`
	ContactEmail graphql.String `json:"contactEmail"`
	Website      graphql.String `json:"website"`
}

// Keywords wraps a list of keyword strings.
type Keywords struct {
	Keywords []graphql.String `json:"keywords"`
}

// Features lists extra benefits available with the course.
type Features struct {
	LifetimeAccess     graphql.Boolean `json:"lifetimeAccess"`
	OneOnOneMentorship graphql.Boolean `json:"oneOnOneMentorship"`
	ReferralProgram    graphql.Boolean `json:"referralProgram"`
	LiveSupport        graphql.Boolean `json:"liveSupport"`
}
