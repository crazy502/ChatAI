package response

import "server/pkg/code"

type Response struct {
	StatusCode code.Code `json:"status_code"`
	StatusMsg  string    `json:"status_msg,omitempty"`
}

func (r *Response) CodeOf(c code.Code) Response {
	if r == nil {
		r = new(Response)
	}
	r.StatusCode = c
	r.StatusMsg = c.Msg()
	return *r
}

func (r *Response) Success() {
	r.CodeOf(code.CodeSuccess)
}
