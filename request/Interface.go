package request

import (
	"github.com/ysfgrl/fiber-pkg/response"
	"github.com/ysfgrl/gerror"
)

type Interface[CType any] interface {
	GetById(id string) (*CType, *gerror.Error)
	List(schema response.ListRequest) (*response.ListResponse[CType], *gerror.Error)
}
