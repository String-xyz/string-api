package unit21

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Transaction interface {
	Evaluate(transaction model.Transaction) (pass bool, err error)
	Create(transaction model.Transaction) (unit21Id string, err error)
	Update(transaction model.Transaction) (unit21Id string, err error)
}

type TransactionRepo struct {
	TxLeg repository.TxLeg
	User  repository.User
	Asset repository.Asset
}

type transaction struct {
	repo TransactionRepo
}

func NewTransaction(r TransactionRepo) Transaction {
	return &transaction{repo: r}
}

func (t transaction) Evaluate(transaction model.Transaction) (pass bool, err error) {
	transactionData, err := t.getTransactionData(transaction)
	if err != nil {
		log.Printf("Failed to gather Unit21 transaction source: %s", err)
		return false, common.StringError(err)
	}
	url := "https://rtr.sandbox2.unit21.com/evaluate" // will need to be hardcoded for production
	body, err := u21Post(url, mapToUnit21Event(transaction, transactionData))

	if err != nil {
		log.Printf("Unit21 Transaction evaluate failed: %s", err)
		return false, common.StringError(err)
	}

	// var u21Response *createEventResponse
	var response any
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return false, common.StringError(err)
	}
	log.Printf("unit21 evaluate response: %+v", response)

	// log.Printf("Unit21Id: %s")

	return true, nil
}

func (t transaction) Create(transaction model.Transaction) (unit21Id string, err error) {
	transactionData, err := t.getTransactionData(transaction)

	if err != nil {
		log.Printf("Failed to gather Unit21 transaction source: %s", err)
		return "", common.StringError(err)
	}

	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/events/create"
	body, err := u21Post(url, mapToUnit21Event(transaction, transactionData))
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

func (t transaction) Update(transaction model.Transaction) (unit21Id string, err error) {
	transactionData, err := t.getTransactionData(transaction)
	if err != nil {
		log.Printf("Failed to gather Unit21 transaction source: %s", err)
		return "", common.StringError(err)
	}

	orgName := os.Getenv("UNIT21_ORG_NAME")
	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/" + orgName + "/events/" + transaction.ID + "/update"
	body, err := u21Put(url, mapToUnit21Event(transaction, transactionData))

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

func (t transaction) getTransactionData(transaction model.Transaction) (txData transactionData, err error) {
	senderData, err := t.repo.TxLeg.GetById(transaction.OriginTxLegID)
	if err != nil {
		log.Printf("Failed go get origin transaction leg: %s", err)
		err = common.StringError(err)
		return
	}

	receiverData, err := t.repo.TxLeg.GetById(transaction.DestinationTxLegID)
	if err != nil {
		log.Printf("Failed go get origin transaction leg: %s", err)
		err = common.StringError(err)
		return
	}

	senderAsset, err := t.repo.Asset.GetById(senderData.AssetID)
	if err != nil {
		log.Printf("Failed go get transaction sender asset: %s", err)
		err = common.StringError(err)
		return
	}

	receiverAsset, err := t.repo.Asset.GetById(receiverData.AssetID)
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

	fmt.Printf("senderAsset: %+v\n", senderAsset)
	log.Printf("senderAsset.Name: %s", senderAsset.Name)
	log.Printf("receiverAsset.Name: %s", receiverAsset.Name)

	txData = transactionData{
		Amount:               amount,
		SentAmount:           senderAmount,
		SentCurrency:         senderAsset.Name,
		SenderEntityId:       senderData.UserID,
		SenderEntityType:     "user",
		SenderInstrumentId:   senderData.InstrumentID,
		ReceivedAmount:       receiverAmount,
		ReceivedCurrency:     receiverAsset.Name,
		ReceiverEntityId:     receiverData.UserID,
		ReceiverEntityType:   "user",
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
			EventId:      transaction.ID,                    //required
			EventType:    "transaction",                     //required
			EventTime:    int(transaction.CreatedAt.Unix()), //required
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
