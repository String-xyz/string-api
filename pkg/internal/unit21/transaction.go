package unit21

import (
	"encoding/json"
	"log"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Transaction interface {
	Create(transaction model.Transaction) (unit21Id string, err error)
	Update(transaction model.Transaction) (unit21Id string, err error)
}

type TransactionRepo struct {
	txLeg repository.TxLeg
	user  repository.User
	asset repository.Asset
}

type transaction struct {
	repo TransactionRepo
}

func NewTransaction(r TransactionRepo) Transaction {
	return &transaction{repo: r}
}

func (i transaction) Create(transaction model.Transaction) (unit21Id string, err error) {

	transactionData, err := getTransactionData(transaction, i.repo.user, i.repo.asset, i.repo.txLeg)
	if err != nil {
		log.Printf("Failed to gather Unit21 transaction source: %s", err)
		return "", common.StringError(err)
	}

	body, err := create("events", mapToUnit21Event(transaction, transactionData))
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

	transactionData, err := getTransactionData(transaction, i.repo.user, i.repo.asset, i.repo.txLeg)
	if err != nil {
		log.Printf("Failed to gather Unit21 transaction source: %s", err)
		return "", common.StringError(err)
	}

	body, err := update("events", transaction.ID, mapToUnit21Event(transaction, transactionData))
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

func getTransactionData(transaction model.Transaction, userRepo repository.User, assetRepo repository.Asset, txLegRepo repository.TxLeg) (txData transactionData, err error) {
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

	senderAsset, err := assetRepo.GetById(senderData.AssetID)
	if err != nil {
		log.Printf("Failed go get transaction sender asset: %s", err)
		err = common.StringError(err)
		return
	}

	receiverAsset, err := assetRepo.GetById(receiverData.AssetID)
	if err != nil {
		log.Printf("Failed go get transaction receiver asset: %s", err)
		err = common.StringError(err)
		return
	}

	amount, err := common.BigNumberToFloat(senderData.Value, 6)
	if err != nil {
		log.Printf("Failed to convert amount: %s", err)
		err = common.StringError(err)
		return
	}

	senderAmount, err := common.BigNumberToFloat(senderData.Amount, senderAsset.Decimals)
	if err != nil {
		log.Printf("Failed to convert senderAmount: %s", err)
		err = common.StringError(err)
		return
	}

	receiverAmount, err := common.BigNumberToFloat(receiverData.Amount, receiverAsset.Decimals)
	if err != nil {
		log.Printf("Failed to convert receiverAmount: %s", err)
		err = common.StringError(err)
		return
	}

	stringFee, err := common.BigNumberToFloat(transaction.StringFee, 6)
	if err != nil {
		log.Printf("Failed to convert stringFee: %s", err)
		err = common.StringError(err)
		return
	}

	processingFee, err := common.BigNumberToFloat(transaction.ProcessingFee, 6)
	if err != nil {
		log.Printf("Failed to convert processingFee: %s", err)
		err = common.StringError(err)
		return
	}

	txData = transactionData{
		Amount:               amount,
		SentAmount:           senderAmount,
		SentCurrency:         senderAsset.Name,
		SenderEntityId:       senderData.UserID,
		SenderEntityType:     senderType,
		SenderInstrumentId:   senderData.InstrumentID,
		ReceivedAmount:       receiverAmount,
		ReceivedCurrency:     receiverAsset.Name,
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

	jsonBody := &u21Event{
		GeneralData: &eventGeneral{
			EventId:      transaction.ID,
			EventType:    "transaction",
			EventTime:    int(transaction.CreatedAt.Unix()),
			EventSubtype: "",
			Status:       transaction.Status,
			Parents:      nil,
			Tags:         transactionTagArr,
		},
		TransactionData: &transactionData,
		ActionData:      nil,
		DigitalData: &eventDigitalData{
			IPAddress: transaction.IPAddress,
		},
		LocationData: nil,
		CustomData:   nil,
	}

	return jsonBody
}
