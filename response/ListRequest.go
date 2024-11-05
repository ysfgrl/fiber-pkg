package response

import "time"

type ListRequest struct {
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Filters  map[string]any `json:"filters"`
	Aggs     map[string]any `json:"aggs"`
}

func ListRequestLastDay() ListRequest {
	return ListRequest{
		Page:     1,
		PageSize: 10,
		Filters: map[string]any{
			"lte": time.Now().UTC(),
			"gte": time.Now().UTC().AddDate(0, 0, -1),
		},
		Aggs: nil,
	}
}

func ListRequestLastMonth() ListRequest {
	return ListRequest{
		Page:     1,
		PageSize: 10,
		Filters: map[string]any{
			"lte": time.Now().UTC(),
			"gte": time.Now().UTC().AddDate(0, -1, 0),
		},
		Aggs: nil,
	}
}

func ListRequestLastYear() ListRequest {
	return ListRequest{
		Page:     1,
		PageSize: 10,
		Filters: map[string]any{
			"lte": time.Now().UTC(),
			"gte": time.Now().UTC().AddDate(-1, 0, 0),
		},
		Aggs: nil,
	}
}
