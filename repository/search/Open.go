package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/ysfgrl/fiber-pkg/response"
	"github.com/ysfgrl/gerror"
	"strings"
	"time"
)

type Open[IType any] struct {
	OpenClient *opensearchapi.Client
	IndexName  string
	FilterKeys []string
}

func (es Open[IType]) Add(ctx context.Context, document IType) (bool, *gerror.Error) {
	data, err := json.Marshal(document)
	if err != nil {
		//es.OpenClient.Index(es.IndexName, bytes.NewReader(data))
		return false, gerror.GetError(err)
	}
	_, err = es.OpenClient.Index(ctx, opensearchapi.IndexReq{
		Index: es.IndexName,
		Body:  bytes.NewReader(data),
	})
	if err != nil {
		return false, gerror.GetError(err)
	}
	return true, nil
}

func (es Open[IType]) List(ctx context.Context, request response.ListRequest) (*response.ListResponse[IType], json.RawMessage, *gerror.Error) {

	var dateRange = map[string]interface{}{
		"gt": time.Now().AddDate(-5, 0, 0).UTC(),
		"lt": time.Now().UTC(),
	}
	var terms []map[string]interface{}
	var filter []map[string]interface{}
	var dna map[string]interface{}
	for key, val := range request.Filters {
		if key == "keyword" {
			for _, filterKey := range es.FilterKeys {
				filter = append(filter, map[string]interface{}{
					"wildcard": map[string]string{
						filterKey: "*" + val.(string) + "*",
					},
				})
			}

		} else if key == "lte" {
			dateRange["lt"] = val
		} else if key == "gte" {
			dateRange["gt"] = val
		} else if key == "dna" {
			dna = map[string]interface{}{
				"value": val,
			}
			continue
		} else {
			terms = append(terms, map[string]interface{}{
				"term": map[string]string{
					key: val.(string),
				},
			})
		}
	}

	must := []map[string]interface{}{
		{
			"range": map[string]interface{}{
				"createdAt": dateRange,
			},
		},
	}
	if dna != nil {
		must = append(must, map[string]interface{}{
			"prefix": map[string]interface{}{
				"dna.keyword": dna,
			},
		})
	}
	if len(filter) > 0 {
		must = append(must, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": filter,
			},
		})
	}
	for _, term := range terms {
		must = append(must, term)
	}

	searchMap := map[string]interface{}{
		"size": 0,
		"from": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"sort": []map[string]interface{}{
			{
				"createdAt": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}
	if request.Aggs != nil {
		searchMap["aggs"] = request.Aggs
	} else {
		searchMap["size"] = request.PageSize
		searchMap["from"] = (request.Page - 1) * request.PageSize
	}
	searchByte, err := json.Marshal(searchMap)
	if err != nil {
		return nil, nil, gerror.GetError(err)
	}
	searchResp, err := es.OpenClient.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{es.IndexName},
		Body:    strings.NewReader(string(searchByte)),
	})
	if err != nil {
		return nil, nil, gerror.GetError(err)
	}
	if searchResp.Aggregations != nil {
		return nil, searchResp.Aggregations, nil
	}
	var list []IType
	for _, hit := range searchResp.Hits.Hits {
		var item IType
		err1 := json.Unmarshal(hit.Source, &item)
		if err1 != nil {
			return nil, nil, gerror.GetError(err)
		}
		json.Unmarshal([]byte(fmt.Sprintf(`{"id":"%s"}`, hit.ID)), &item)
		list = append(list, item)
	}
	if len(searchResp.Hits.Hits) == 0 {
		list = make([]IType, 0)
	}
	return &response.ListResponse[IType]{
		Page:     request.Page,
		PageSize: len(searchResp.Hits.Hits),
		Total:    int64(searchResp.Hits.Total.Value),
		List:     list,
	}, nil, nil
}
