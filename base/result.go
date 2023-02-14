package base

type CollectionResult struct {
	List interface{} `json:"list"`
	Page *PageInfo   `json:"page"`
}

type PageInfo struct {
	Count    int64 `json:"count"`
	PageNO   int   `json:"pageNo"`
	PageSize int   `json:"pageSize"`
}

func NewCollectionResult(list interface{}, pageNO, pageSize int, count int64) *CollectionResult {
	return &CollectionResult{
		List: list,
		Page: &PageInfo{
			PageNO:   pageNO,
			PageSize: pageSize,
			Count:    count,
		},
	}
}
