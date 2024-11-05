package response

import "github.com/ysfgrl/gerror"

type Api[CType any] struct {
	Code    int           `json:"code"`
	Content CType         `json:"content"`
	Error   *gerror.Error `json:"error"`
}
