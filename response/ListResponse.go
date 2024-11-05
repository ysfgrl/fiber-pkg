package response

type ListResponse[CType any] struct {
	List     []CType `json:"list"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
	Total    int64   `json:"total"`
}
