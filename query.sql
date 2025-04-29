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

