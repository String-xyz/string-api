package service

import (
	"github.com/String-xyz/string-api/pkg/internal/unit21"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Unit21 struct {
	Action      unit21.Action
	Entity      unit21.Entity
	Instrument  unit21.Instrument
	Transaction unit21.Transaction
}

func NewUnit21(repos repository.Repositories) Unit21 {
	action := unit21.NewAction()
	entityRepos := unit21.EntityRepos{Device: repos.Device, Contact: repos.Contact, User: repos.User}
	entity := unit21.NewEntity(entityRepos)
	instrumentRepos := unit21.InstrumentRepos{User: repos.User, Device: repos.Device, Location: repos.Location}
	instrument := unit21.NewInstrument(instrumentRepos, action)
	transactionRepos := unit21.TransactionRepos{User: repos.User, TxLeg: repos.TxLeg, Asset: repos.Asset, Device: repos.Device}
	transaction := unit21.NewTransaction(transactionRepos)
	return Unit21{
		Action:      action,
		Entity:      entity,
		Instrument:  instrument,
		Transaction: transaction,
	}
}
