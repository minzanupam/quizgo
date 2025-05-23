// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package queries

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type QuizStatus string

const (
	QuizStatusPublished   QuizStatus = "published"
	QuizStatusUnpublished QuizStatus = "unpublished"
	QuizStatusExpired     QuizStatus = "expired"
)

func (e *QuizStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = QuizStatus(s)
	case string:
		*e = QuizStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for QuizStatus: %T", src)
	}
	return nil
}

type NullQuizStatus struct {
	QuizStatus QuizStatus
	Valid      bool // Valid is true if QuizStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullQuizStatus) Scan(value interface{}) error {
	if value == nil {
		ns.QuizStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.QuizStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullQuizStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.QuizStatus), nil
}

type Option struct {
	ID         int32
	Body       string
	QuestionID int32
}

type Question struct {
	ID     int32
	Body   string
	QuizID int32
}

type Quiz struct {
	ID        int32
	Title     string
	OwnerID   int32
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Status    QuizStatus
}

type User struct {
	ID       int32
	Fullname string
	Email    string
	Password string
}
