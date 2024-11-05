package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/ysfgrl/fiber-pkg/response"
	"github.com/ysfgrl/gerror"
)

type Elastic[IType any] struct {
	ElasticClient *elasticsearch.TypedClient
	IndexName     string
}

func (es Elastic[IType]) Add(ctx context.Context, document IType) (bool, *gerror.Error) {

	_, err := es.ElasticClient.Index(es.IndexName).Request(document).Do(ctx)
	if err != nil {
		return false, gerror.GetError(err)
	}
	return true, nil
}

func (es Elastic[IType]) List(ctx context.Context, request response.ListRequest) (*response.ListResponse[IType], *gerror.Error) {
	resRequest := search.Request{
		Query: &types.Query{
			MatchAll: &types.MatchAllQuery{},
		},
	}
	res, err := es.ElasticClient.Search().Index(es.IndexName).Request(&resRequest).Do(ctx)

	if err != nil {
		return nil, gerror.GetError(err)
	}
	var list []IType
	for _, hit := range res.Hits.Hits {
		var item IType
		err1 := json.Unmarshal(hit.Source_, &item)
		if err1 != nil {
			return nil, gerror.GetError(err)
		}
		json.Unmarshal([]byte(fmt.Sprintf(`{"id":"%s"}`, hit.Id_)), &item)
		list = append(list, item)
	}
	return &response.ListResponse[IType]{
		Page:     0,
		PageSize: len(res.Hits.Hits),
		Total:    res.Hits.Total.Value,
		List:     list,
	}, nil
}
