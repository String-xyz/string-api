package repository

type Transaction interface {
}

type transaction struct {
}

func NewTransaction(db any) Transaction {
	return &transaction{}
}
