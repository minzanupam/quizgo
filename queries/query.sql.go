// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createQuiz = `-- name: CreateQuiz :one
INSERT INTO quizzes (title, owner_id) VALUES ($1, $2) RETURNING ID
`

type CreateQuizParams struct {
	Title   string
	OwnerID int32
}

func (q *Queries) CreateQuiz(ctx context.Context, arg CreateQuizParams) (int32, error) {
	row := q.db.QueryRow(ctx, createQuiz, arg.Title, arg.OwnerID)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getQuestion = `-- name: GetQuestion :many
SELECT
	questions.ID, questions.quiz_id, questions.body, options.ID, options.Body
FROM
	questions
INNER JOIN
	options
ON
	options.question_id = questions.ID
WHERE
	questions.ID = $1
`

type GetQuestionRow struct {
	ID     int32
	QuizID int32
	Body   string
	ID_2   int32
	Body_2 string
}

func (q *Queries) GetQuestion(ctx context.Context, id int32) ([]GetQuestionRow, error) {
	rows, err := q.db.Query(ctx, getQuestion, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetQuestionRow
	for rows.Next() {
		var i GetQuestionRow
		if err := rows.Scan(
			&i.ID,
			&i.QuizID,
			&i.Body,
			&i.ID_2,
			&i.Body_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQuiz = `-- name: GetQuiz :many
SELECT
	quiz.ID,
	quiz.title,
	quiz.created_at,
	quiz.updated_at,
	quiz.status,
	ques.ID,
	ques.body,
	opt.ID,
	opt.body
FROM
	quizzes AS quiz
LEFT JOIN
	questions AS ques
ON
	quiz.ID = ques.quiz_id
LEFT JOIN
	options AS opt
ON
	ques.ID = opt.question_id
WHERE
	quiz.ID = $1
AND
	owner_id = $2
`

type GetQuizParams struct {
	ID      int32
	OwnerID int32
}

type GetQuizRow struct {
	ID        int32
	Title     string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	Status    QuizStatus
	ID_2      pgtype.Int4
	Body      pgtype.Text
	ID_3      pgtype.Int4
	Body_2    pgtype.Text
}

func (q *Queries) GetQuiz(ctx context.Context, arg GetQuizParams) ([]GetQuizRow, error) {
	rows, err := q.db.Query(ctx, getQuiz, arg.ID, arg.OwnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetQuizRow
	for rows.Next() {
		var i GetQuizRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Status,
			&i.ID_2,
			&i.Body,
			&i.ID_3,
			&i.Body_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQuizStatus = `-- name: GetQuizStatus :one
SELECT status FROM quizzes WHERE ID = $1 AND owner_id = $2
`

type GetQuizStatusParams struct {
	ID      int32
	OwnerID int32
}

func (q *Queries) GetQuizStatus(ctx context.Context, arg GetQuizStatusParams) (QuizStatus, error) {
	row := q.db.QueryRow(ctx, getQuizStatus, arg.ID, arg.OwnerID)
	var status QuizStatus
	err := row.Scan(&status)
	return status, err
}

const updateQuizStatusPublish = `-- name: UpdateQuizStatusPublish :exec
UPDATE quizzes SET status = 'published' WHERE ID = $1
`

func (q *Queries) UpdateQuizStatusPublish(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, updateQuizStatusPublish, id)
	return err
}
