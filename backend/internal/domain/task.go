package domain

// Concrete task to run
type Task struct {
	Code  string
	Lang  string
	ReqId string
	Resp  chan ExecuteResponse
}

type ExecuteResponse struct {
	Resp string // TODO подумать как парсить аутпут выполнения кода в контейнере
}
