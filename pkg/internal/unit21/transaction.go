package unit21

import (
	"encoding/json"
	"os"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/rs/zerolog/log"
)

type Transaction interface {
	Evaluate(transaction model.Transaction) (pass bool, err error)
	Create(transaction model.Transaction) (unit21Id string, err error)
	Update(transaction model.Transaction) (unit21Id string, err error)
}

type TransactionRepos struct {
	User  repository.User
	TxLeg repository.TxLeg
	Asset repository.Asset
}

type transaction struct {
	repos TransactionRepos
}

func NewTransaction(r TransactionRepos) Transaction {
	return &transaction{repos: r}
}

func (t transaction) Evaluate(transaction model.Transaction) (pass bool, err error) {
	transactionData, err := t.getTransactionData(transaction)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 transaction source")
		return false, common.StringError(err)
	}

	url := os.Getenv("UNIT21_RTR_URL")
	if url == "" {
		url = "https://rtr.sandbox2.unit21.com/evaluate"
	}

	body, err := u21Post(url, mapToUnit21TransactionEvent(transaction, transactionData))
	if err != nil {
		log.Err(err).Msg("Unit21 Transaction evaluate failed")
		return false, common.StringError(err)
	}

	// var u21Response *createEventResponse
	var response evaluateEventResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return false, common.StringError(err)
	}

	for _, rule := range *response.RuleExecutions {
		if rule.Status != "PASS" {
			return false, nil
		}
	}

	return true, nil
}

func (t transaction) Create(transaction model.Transaction) (unit21Id string, err error) {
	transactionData, err := t.getTransactionData(transaction)

	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 transaction source")
		return "", common.StringError(err)
	}

	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/events/create"
	body, err := u21Post(url, mapToUnit21TransactionEvent(transaction, transactionData))
	if err != nil {
		log.Err(err).Msg("Unit21 Transaction create failed")
		return "", common.StringError(err)
	}

	var u21Response *createEventResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", common.StringError(err)
	}

	log.Info().Str("unit21Id", u21Response.Unit21Id).Send()
	return u21Response.Unit21Id, nil
}

func (t transaction) Update(transaction model.Transaction) (unit21Id string, err error) {
	transactionData, err := t.getTransactionData(transaction)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 transaction source")
		return "", common.StringError(err)
	}

	orgName := os.Getenv("UNIT21_ORG_NAME")
	url := "https://" + os.Getenv("UNIT21_ENV") + ".unit21.com/v1/" + orgName + "/events/" + transaction.ID + "/update"
	body, err := u21Put(url, mapToUnit21TransactionEvent(transaction, transactionData))

	if err != nil {
		log.Err(err).Msg("Unit21 Transaction create failed:")
		return "", common.StringError(err)
	}

	var u21Response *updateEventResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", common.StringError(err)
	}
	log.Info().Str("unit21Id", u21Response.Unit21Id).Send()
	return u21Response.Unit21Id, nil
}

func (t transaction) getTransactionData(transaction model.Transaction) (txData transactionData, err error) {
	senderData, err := t.repos.TxLeg.GetById(transaction.OriginTxLegID)
	if err != nil {
		log.Err(err).Msg("Failed go get origin transaction leg")
		err = common.StringError(err)
		return
	}

	receiverData, err := t.repos.TxLeg.GetById(transaction.DestinationTxLegID)
	if err != nil {
		log.Err(err).Msg("Failed go get origin transaction leg")
		err = common.StringError(err)
		return
	}

	senderAsset, err := t.repos.Asset.GetById(senderData.AssetID)
	if err != nil {
		log.Err(err).Msg("Failed go get transaction sender asset")
		err = common.StringError(err)
		return
	}

	receiverAsset, err := t.repos.Asset.GetById(receiverData.AssetID)
	if err != nil {
		log.Err(err).Msg("Failed go get transaction receiver asset")
		err = common.StringError(err)
		return
	}

	amount, err := common.BigNumberToFloat(senderData.Value, 6)
	if err != nil {
		log.Err(err).Msg("Failed to convert amount")
		err = common.StringError(err)
		return
	}

	senderAmount, err := common.BigNumberToFloat(senderData.Amount, senderAsset.Decimals)
	if err != nil {
		log.Err(err).Msg("Failed to convert senderAmount")
		err = common.StringError(err)
		return
	}

	receiverAmount, err := common.BigNumberToFloat(receiverData.Amount, receiverAsset.Decimals)
	if err != nil {
		log.Err(err).Msg("Failed to convert receiverAmount")
		err = common.StringError(err)
		return
	}
	var stringFee float64
	if transaction.StringFee != "" {
		stringFee, err = common.BigNumberToFloat(transaction.StringFee, 6)
		if err != nil {
			log.Err(err).Msg("Failed to convert stringFee")
			err = common.StringError(err)
			return
		}
	}

	var processingFee float64
	if transaction.ProcessingFee != "" {
		processingFee, err = common.BigNumberToFloat(transaction.ProcessingFee, 6)
		if err != nil {
			log.Err(err).Msg("Failed to convert processingFee")
			err = common.StringError(err)
			return
		}
	}

	var exchangeRate float64
	if receiverAmount > 0 {
		exchangeRate = senderAmount / receiverAmount
	}

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
		ExchangeRate:         exchangeRate,
		TransactionHash:      transaction.TransactionHash,
		USDConversionNotes:   "",
		InternalFee:          stringFee,
		ExternalFee:          processingFee,
	}

	return
}

func mapToUnit21TransactionEvent(transaction model.Transaction, transactionData transactionData) *u21Event {
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
			EventSubtype: "Card Payment",                    //required for RTR
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
