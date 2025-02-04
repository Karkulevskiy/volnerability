package domain

type Task struct {
	Code  string
	Lang  string
	ReqId string
}

type ExecuteResponse struct {
	Resp string // TODO подумать как парсить аутпут выполнения кода в контейнере
}
