package integrationtests

import (
	"time"
)

type User struct {
	ID        int64
	Email     string
	FirstName string
	LastName  string
}

type Event struct {
	ID          int64
	UserID      int64
	Title       string
	Status      string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	NotifiedFor int64
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UpdateEventPayload struct {
	Event Event `json:"event"`
}

type QueryEventsResponse struct {
	Events []Event `json:"events"`
}
