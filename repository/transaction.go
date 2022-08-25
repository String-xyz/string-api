package repository

type Transaction interface {
}

type transaction struct {
}

func NewTransaction() Transaction {
	return &transaction{}
}
