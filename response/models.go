package response

import "github.com/ysfgrl/gerror"

type Response struct {
	Code    int           `json:"code"`
	Content any           `json:"content"`
	Error   *gerror.Error `json:"error"`
}

type Ok struct {
	IsOk bool `json:"isOk"`
}
type Str struct {
	Value string `json:"value"`
}
