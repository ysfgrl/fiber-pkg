package request

import (
	"encoding/json"
	"github.com/ysfgrl/fiber-pkg/response"
	"github.com/ysfgrl/gerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Base[CType any] struct {
	Client HttpService
}

func (b *Base[CType]) Create(item CType) (bool, *gerror.Error) {
	bodyByte, err1 := json.Marshal(&item)
	if err1 != nil {
		return false, gerror.GetError(err1)
	}
	res, err := b.Client.Post("/create", bodyByte)
	if err != nil {
		return false, err
	}
	resS := new(response.Api[response.Ok])
	if err1 := json.Unmarshal(res, resS); err1 != nil {
		return false, gerror.GetError(err1)
	}
	if resS.Error != nil {
		return false, resS.Error
	}
	return resS.Content.IsOk, nil
}

func (b *Base[CType]) GetById(id *primitive.ObjectID) (*CType, *gerror.Error) {
	if id == nil {
		return nil, &gerror.Error{Code: ""}
	}
	res, err := b.Client.Get("/get/" + id.Hex())
	if err != nil {
		return nil, err
	}
	resS := new(response.Api[CType])
	if err1 := json.Unmarshal(res, resS); err1 != nil {
		return nil, gerror.GetError(err1)
	}
	if resS.Error != nil {
		return nil, resS.Error
	}
	return &resS.Content, nil
}

func (b *Base[CType]) GetDetail(id *primitive.ObjectID) (*CType, *gerror.Error) {
	if id == nil {
		return nil, &gerror.Error{Code: ""}
	}
	res, err := b.Client.Get("/detail/" + id.Hex())
	if err != nil {
		return nil, err
	}
	resS := new(response.Api[CType])
	if err1 := json.Unmarshal(res, resS); err1 != nil {
		return nil, gerror.GetError(err1)
	}
	if resS.Error != nil {
		return nil, resS.Error
	}
	return &resS.Content, nil
}

func (b *Base[CType]) SetField(id primitive.ObjectID, key string, val any) (bool, *gerror.Error) {
	body := map[string]interface{}{key: val}
	bodyByte, err1 := json.Marshal(body)
	if err1 != nil {
		return false, gerror.GetError(err1)
	}
	res, err := b.Client.Put("/setField/"+id.Hex(), bodyByte)
	if err != nil {
		return false, err
	}
	resS := new(response.Api[response.Ok])
	if err1 := json.Unmarshal(res, resS); err1 != nil {
		return false, gerror.GetError(err1)
	}
	if resS.Error != nil {
		return false, resS.Error
	}
	return resS.Content.IsOk, nil
}

func (b *Base[CType]) SetFields(id primitive.ObjectID, fields map[string]interface{}) (bool, *gerror.Error) {
	bodyByte, err1 := json.Marshal(fields)
	if err1 != nil {
		return false, gerror.GetError(err1)
	}
	res, err := b.Client.Put("/setFields/"+id.Hex(), bodyByte)
	if err != nil {
		return false, err
	}
	resS := new(response.Api[response.Ok])
	if err1 := json.Unmarshal(res, resS); err1 != nil {
		return false, gerror.GetError(err1)
	}
	if resS.Error != nil {
		return false, resS.Error
	}
	return resS.Content.IsOk, nil
}

func (b *Base[CType]) List(filter response.ListRequest) (*response.ListResponse[CType], *gerror.Error) {
	filterByte, err1 := json.Marshal(filter)
	if err1 != nil {
		return nil, gerror.GetError(err1)
	}
	res, err := b.Client.Post("/list", filterByte)
	if err != nil {
		return nil, err
	}
	resS := new(response.Api[response.ListResponse[CType]])
	if err1 = json.Unmarshal(res, resS); err1 != nil {
		return nil, gerror.GetError(err1)
	}
	if resS.Error != nil {

		return nil, resS.Error
	}
	return &resS.Content, nil
}
