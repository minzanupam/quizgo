-- name: GetQuiz :many
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
	owner_id = $2;


-- name: CreateQuiz :one
INSERT INTO quizzes (title, owner_id) VALUES ($1, $2) RETURNING ID;

-- name: GetQuizStatus :one
SELECT status FROM quizzes WHERE ID = $1 AND owner_id = $2;

-- name: UpdateQuizStatusPublish :exec
UPDATE quizzes SET status = 'published' WHERE ID = $1;

-- name: GetQuestion :many
SELECT
	questions.ID, questions.quiz_id, questions.body, options.ID, options.Body
FROM
	questions
INNER JOIN
	options
ON
	options.question_id = questions.ID
WHERE
	questions.ID = $1;
