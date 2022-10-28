package unit21

import (
	"encoding/json"
	"log"
	"math"
	"strconv"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Transaction interface {
	Create(transaction model.Transaction) (unit21Id string, err error)
	Update(transaction model.Transaction) (unit21Id string, err error)
}

type transaction struct {
	transactionRepo repository.Transaction
	txLegRepo       repository.TxLeg
	userRepo        repository.User
	assetRepo       repository.Asset
}

func NewTransaction(tx repository.Transaction, txLeg repository.TxLeg, user repository.User, asset repository.Asset) Transaction {
	return &transaction{transactionRepo: tx, txLegRepo: txLeg, userRepo: user, assetRepo: asset}
}

func (i transaction) Create(transaction model.Transaction) (unit21Id string, err error) {

	transactionData, err := getTransactionData(transaction, i.userRepo, i.assetRepo, i.txLegRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 transaction source: %s", err)
		return "", common.StringError(err)
	}

	body, err := create("transactions", mapToUnit21Event(transaction, transactionData))
	if err != nil {
		log.Printf("Unit21 Transaction create failed: %s", err)
		return "", common.StringError(err)
	}

	var u21Response *createEventResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "", common.StringError(err)
	}

	log.Printf("Unit21Id: %s", u21Response.Unit21Id)
	return u21Response.Unit21Id, nil
}

func (i transaction) Update(transaction model.Transaction) (unit21Id string, err error) {

	transactionData, err := getTransactionData(transaction, i.userRepo, i.assetRepo, i.txLegRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 transaction source: %s", err)
		return "", common.StringError(err)
	}

	body, err := update("transactions", transaction.ID, mapToUnit21Event(transaction, transactionData))
	if err != nil {
		log.Printf("Unit21 Transaction create failed: %s", err)
		return "", common.StringError(err)
	}

	var u21Response *updateEventResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "", common.StringError(err)
	}

	log.Printf("Unit21Id: %s", u21Response.Unit21Id)
	return u21Response.Unit21Id, nil
}

func getTransactionData(transaction model.Transaction, userRepo repository.User, assetRepo repository.Asset, txLegRepo repository.TxLeg) (transactionData transactionData, err error) {
	senderData, err := txLegRepo.GetById(transaction.OriginTxLegID)
	if err != nil {
		log.Printf("Failed go get origin transaction leg: %s", err)
		err = common.StringError(err)
		return
	}

	receiverData, err := txLegRepo.GetById(transaction.DestinationTxLegID)
	if err != nil {
		log.Printf("Failed go get origin transaction leg: %s", err)
		err = common.StringError(err)
		return
	}

	receiverType, err := getSource(receiverData.UserID, userRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 transaction receiver user source: %s", err)
		err = common.StringError(err)
		return
	}

	senderType, err := getSource(senderData.UserID, userRepo)
	if err != nil {
		log.Printf("Failed to gather Unit21 transaction sender user source: %s", err)
		err = common.StringError(err)
		return
	}

	amount, err := strconv.ParseFloat(senderData.Value, 64)
	if err != nil {
		log.Printf("Failed to convert transaction value to float: %s", err)
		err = common.StringError(err)
		return
	}
	amount = amount * math.Pow(10, -6)

	senderAsset, err := assetRepo.GetById(senderData.AssetID)
	if err != nil {
		log.Printf("Failed go get transaction sender asset: %s", err)
		err = common.StringError(err)
		return
	}

	senderAmount, err := strconv.ParseFloat(senderData.Amount, 64)
	if err != nil {
		log.Printf("Failed to convert transaction amount to float: %s", err)
		err = common.StringError(err)
		return
	}
	senderAmount = senderAmount * math.Pow(10, -float64(senderAsset.Decimals))

	receiverAsset, err := assetRepo.GetById(receiverData.AssetID)
	if err != nil {
		log.Printf("Failed go get transaction receiver asset: %s", err)
		err = common.StringError(err)
		return
	}

	receiverAmount, err := strconv.ParseFloat(receiverData.Amount, 64)
	if err != nil {
		log.Printf("Failed to convert transaction amount to float: %s", err)
		err = common.StringError(err)
		return
	}
	receiverAmount = receiverAmount * math.Pow(10, -float64(receiverAsset.Decimals))

	stringFee, err := strconv.ParseFloat(transaction.StringFee, 64)
	if err != nil {
		log.Printf("Failed to convert transaction string fee to float: %s", err)
		err = common.StringError(err)
		return
	}
	stringFee = stringFee * math.Pow(10, float64(-6))

	processingFee, err := strconv.ParseFloat(transaction.ProcessingFee, 64)
	if err != nil {
		log.Printf("Failed to convert transaction processing fee to float: %s", err)
		err = common.StringError(err)
		return
	}
	processingFee = processingFee * math.Pow(10, float64(-6))

	transactionData = &transactionData{
		Amount:               amount,
		SenderAmount:         senderAmount,
		SenderCurrency:       senderAsset.Name,
		SenderEntityId:       senderData.UserID,
		SenderEntityType:     senderType,
		SenderInstrumentId:   senderData.InstrumentID,
		ReceiverAmount:       receiverAmount,
		ReceiverCurrency:     receiverAsset.Name,
		ReceiverEntityId:     receiverData.UserID,
		ReceiverEntityType:   receiverType,
		ReceiverInstrumentId: receiverData.InstrumentID,
		ExchangeRate:         senderAmount / receiverAmount,
		TransactionHash:      transaction.TransactionHash,
		USDConversionNotes:   "",
		InternalFee:          stringFee,
		ExternalFee:          processingFee,
	}

	return
}

func mapToUnit21Event(transaction model.Transaction, transactionData transactionData) *u21Event {
	var transactionTagArr []string
	if transaction.Tags != nil {
		for key, value := range transaction.Tags {
			transactionTagArr = append(transactionTagArr, key+":"+value)
		}
	}

	//transaction.IPAddress

	// var entityArray []transactionEntity
	// entityArray = append(entityArray, entityData)

	jsonBody := &u21Event{
		GeneralData: &eventGeneral{
			EventId:      transaction.ID,
			EventType:    transaction.Type,
			EventTime:    int(transaction.CreatedAt.Unix()),
			EventSubtype: "",
			Status:       transaction.Status,
			Parents:      nil,
			Tags:         transactionTagArr,
		},
		TransactionData: &transactionData,
		ActionData:      nil,
		DigitalData:     nil,
		LocationData:    nil,
		CustomData:      nil,
	}

	return jsonBody
}
