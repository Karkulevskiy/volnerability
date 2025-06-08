package domain

import "database/sql"

type User struct {
	ID            int64
	Email         string
	PassHash      []byte
	TotalAttempts int
	PassLevels    int
	OauthID       sql.NullInt64
	IsOauth       bool
}

type UserLevel struct {
	Id              int    `json:"id"`
	LevelId         int    `json:"levelId"`
	UserId          int    `json:"userId"`
	IsCompleted     bool   `json:"isCompleted"`
	LastInput       string `json:"lastInput"`
	AttemptResponse string `json:"attemptResponse"` // last response from server on btn submit
	Attempts        int    `json:"attempts"`
}

type Level struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Hints         []string `json:"hints"`
	ExpectedInput string   `json:"expectedInput,omitempty"`
	StartInput    string   `json:"startInput,omitempty"`
}

type Hint struct {
	Id      int    `json:"id"`
	LevelId int    `json:"levelId"`
	Text    string `json:"text"`
}
