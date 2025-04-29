CREATE TABLE users (
	ID SERIAL PRIMARY KEY,
	fullname VARCHAR NOT NULL,
	email VARCHAR NOT NULL UNIQUE,
	password VARCHAR NOT NULL
);

CREATE TYPE quiz_status AS ENUM ('published', 'unpublished', 'expired');

CREATE TABLE quizzes (
	ID SERIAL PRIMARY KEY,
	title VARCHAR NOT NULL,
	owner_id INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	status QUIZ_STATUS NOT NULL,
	FOREIGN KEY (owner_id) REFERENCES users(ID)
);

CREATE TABLE questions (
	ID SERIAL PRIMARY KEY,
	body VARCHAR NOT NULL,
	quiz_id INT NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	FOREIGN KEY (quiz_id) REFERENCES quizzes(ID)
);

CREATE TABLE options (
	ID SERIAL PRIMARY KEY,
	body VARCHAR NOT NULL,
	question_id INT NOT NULL,
	FOREIGN KEY (question_id) REFERENCES questions(ID)
);
