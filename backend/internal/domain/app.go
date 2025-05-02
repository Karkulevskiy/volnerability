package domain

type User struct {
	ID       int64
	Email    string
	PassHash []byte
}

type Level struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Hints         []string `json:"hints"`
	ExpectedInput string   `json:"expectedInput"`
}
