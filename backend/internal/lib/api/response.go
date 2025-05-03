package api

import "net/http"

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"error,omitempty"`
	Model  any    `json:"model,omitempty"`
}

type opt func(*Response)

func OK() Response {
	return Response{Status: http.StatusOK}
}

func BadRequest(msg string) Response {
	return Response{Status: http.StatusBadRequest, Msg: msg}
}

func InternalError(opts ...opt) Response {
	r := &Response{Status: http.StatusInternalServerError, Msg: "internal error"}
	for _, opt := range opts {
		opt(r)
	}
	return *r
}

func New(opts ...opt) Response {
	resp := &Response{}
	for _, opt := range opts {
		opt(resp)
	}
	return *resp
}

func WithCode(statusCode int) opt {
	return func(r *Response) {
		r.Status = statusCode
	}
}

func WithMsg(msg string) opt {
	return func(r *Response) {
		r.Msg = msg
	}
}

func WithModel(a any) opt {
	return func(r *Response) {
		r.Model = a
	}
}
