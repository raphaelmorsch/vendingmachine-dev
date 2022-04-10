package domains

type Purchase struct {
	TotalSpent int `json:"totalSpent"`

	ProductId int `json:"productId"`

	Change map[int]int `json:"change"`
}
