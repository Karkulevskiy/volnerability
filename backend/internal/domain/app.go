package domain

type User struct {
	ID            int64
	Email         string
	PassHash      []byte
	TotalAttempts int
	PassLevels    int
}

type Level struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Hints         []string `json:"hints"`
	ExpectedInput string   `json:"expectedInput"`
}

type Hint struct {
	Id      int    `json:"id"`
	LevelId int    `json:"levelId"`
	Text    string `json:"text"`
}
