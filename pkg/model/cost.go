package model

type CostType string

const (
	Gas   CostType = "Gas"
	Token CostType = "Token"
)

type Cost struct {
	CostType
	Timestamp int `json:"timestamp" db:"timestamp"`
	Value     int `json:"value" db:"value"`
}
