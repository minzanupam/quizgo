package views

type DBUser struct {
	ID       string
	FullName string
	Email    string
}

type DBQuiz struct {
	ID        string
	Title     string
	Owner     DBUser
	Questions []DBQuestion
	CreatedAt string
	UpdatedAt string
	Status    string
}

type DBQuestion struct {
	ID      string
	Body    string
	Options []DBOption
}

type DBOption struct {
	ID   string
	Body string
}
