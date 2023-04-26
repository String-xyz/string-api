package unit21

import (
	"context"
	"encoding/json"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/internal/common"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/rs/zerolog/log"
)

type Transaction interface {
	Evaluate(ctx context.Context, transaction model.Transaction) (pass bool, err error)
	Create(ctx context.Context, transaction model.Transaction) (unit21Id string, err error)
	Update(ctx context.Context, transaction model.Transaction) (unit21Id string, err error)
}

type TransactionRepos struct {
	User   repository.User
	TxLeg  repository.TxLeg
	Asset  repository.Asset
	Device repository.Device
}

type transaction struct {
	repos TransactionRepos
}

func NewTransaction(r TransactionRepos) Transaction {
	return &transaction{repos: r}
}

func (t transaction) Evaluate(ctx context.Context, transaction model.Transaction) (pass bool, err error) {
	transactionData, err := t.getTransactionData(ctx, transaction)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 transaction source")
		return false, libcommon.StringError(err)
	}

	digitalData, err := t.getEventDigitalData(ctx, transaction)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 digital data")
		return false, libcommon.StringError(err)
	}

	url := config.Var.UNIT21_RTR_URL
	if url == "" {
		url = "https://rtr.sandbox2.unit21.com/evaluate"
	}

	body, err := u21Post(url, mapToUnit21TransactionEvent(transaction, transactionData, digitalData))
	if err != nil {
		log.Err(err).Msg("Unit21 Transaction evaluate failed")
		return false, libcommon.StringError(err)
	}

	// var u21Response *createEventResponse
	var response evaluateEventResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return false, libcommon.StringError(err)
	}

	for _, rule := range *response.RuleExecutions {
		if rule.Status != "PASS" {
			return false, nil
		}
	}

	return true, nil
}

func (t transaction) Create(ctx context.Context, transaction model.Transaction) (unit21Id string, err error) {
	transactionData, err := t.getTransactionData(ctx, transaction)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 transaction source")
		return "", libcommon.StringError(err)
	}

	digitalData, err := t.getEventDigitalData(ctx, transaction)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 digital data")
		return "", libcommon.StringError(err)
	}

	url := "https://" + config.Var.UNIT21_ENV + ".unit21.com/v1/events/create"
	body, err := u21Post(url, mapToUnit21TransactionEvent(transaction, transactionData, digitalData))
	if err != nil {
		log.Err(err).Msg("Unit21 Transaction create failed")
		return "", libcommon.StringError(err)
	}

	var u21Response *createEventResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", libcommon.StringError(err)
	}

	log.Info().Str("unit21Id", u21Response.Unit21Id).Send()
	return u21Response.Unit21Id, nil
}

func (t transaction) Update(ctx context.Context, transaction model.Transaction) (unit21Id string, err error) {
	transactionData, err := t.getTransactionData(ctx, transaction)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 transaction source")
		return "", libcommon.StringError(err)
	}

	digitalData, err := t.getEventDigitalData(ctx, transaction)
	if err != nil {
		log.Err(err).Msg("Failed to gather Unit21 digital data")
		return "", libcommon.StringError(err)
	}

	orgName := config.Var.UNIT21_ORG_NAME
	url := "https://" + config.Var.UNIT21_ENV + ".unit21.com/v1/" + orgName + "/events/" + transaction.Id + "/update"
	body, err := u21Put(url, mapToUnit21TransactionEvent(transaction, transactionData, digitalData))

	if err != nil {
		log.Err(err).Msg("Unit21 Transaction create failed:")
		return "", libcommon.StringError(err)
	}

	var u21Response *updateEventResponse
	err = json.Unmarshal(body, &u21Response)
	if err != nil {
		log.Err(err).Msg("Reading body failed")
		return "", libcommon.StringError(err)
	}
	log.Info().Str("unit21Id", u21Response.Unit21Id).Send()
	return u21Response.Unit21Id, nil
}

func (t transaction) getTransactionData(ctx context.Context, transaction model.Transaction) (txData transactionData, err error) {
	senderData, err := t.repos.TxLeg.GetById(ctx, transaction.OriginTxLegId)
	if err != nil {
		log.Err(err).Msg("Failed go get origin transaction leg")
		err = libcommon.StringError(err)
		return
	}

	receiverData, err := t.repos.TxLeg.GetById(ctx, transaction.DestinationTxLegId)
	if err != nil {
		log.Err(err).Msg("Failed go get origin transaction leg")
		err = libcommon.StringError(err)
		return
	}

	senderAsset, err := t.repos.Asset.GetById(ctx, senderData.AssetId)
	if err != nil {
		log.Err(err).Msg("Failed go get transaction sender asset")
		err = libcommon.StringError(err)
		return
	}

	receiverAsset, err := t.repos.Asset.GetById(ctx, receiverData.AssetId)
	if err != nil {
		log.Err(err).Msg("Failed go get transaction receiver asset")
		err = libcommon.StringError(err)
		return
	}

	amount, err := common.BigNumberToFloat(senderData.Value, 6)
	if err != nil {
		log.Err(err).Msg("Failed to convert amount")
		err = libcommon.StringError(err)
		return
	}

	senderAmount, err := common.BigNumberToFloat(senderData.Amount, senderAsset.Decimals)
	if err != nil {
		log.Err(err).Msg("Failed to convert senderAmount")
		err = libcommon.StringError(err)
		return
	}

	receiverAmount, err := common.BigNumberToFloat(receiverData.Amount, receiverAsset.Decimals)
	if err != nil {
		log.Err(err).Msg("Failed to convert receiverAmount")
		err = libcommon.StringError(err)
		return
	}
	var stringFee float64
	if transaction.StringFee != "" {
		stringFee, err = common.BigNumberToFloat(transaction.StringFee, 6)
		if err != nil {
			log.Err(err).Msg("Failed to convert stringFee")
			err = libcommon.StringError(err)
			return
		}
	}

	var processingFee float64
	if transaction.ProcessingFee != "" {
		processingFee, err = common.BigNumberToFloat(transaction.ProcessingFee, 6)
		if err != nil {
			log.Err(err).Msg("Failed to convert processingFee")
			err = libcommon.StringError(err)
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
		SenderEntityId:       senderData.UserId,
		SenderEntityType:     "user",
		SenderInstrumentId:   senderData.InstrumentId,
		ReceivedAmount:       receiverAmount,
		ReceivedCurrency:     receiverAsset.Name,
		ReceiverEntityId:     receiverData.UserId,
		ReceiverEntityType:   "user",
		ReceiverInstrumentId: receiverData.InstrumentId,
		ExchangeRate:         exchangeRate,
		TransactionHash:      transaction.TransactionHash,
		USDConversionNotes:   "",
		InternalFee:          stringFee,
		ExternalFee:          processingFee,
	}

	return
}

func (t transaction) getEventDigitalData(ctx context.Context, transaction model.Transaction) (digitalData eventDigitalData, err error) {
	if transaction.DeviceId == "" {
		return
	}

	device, err := t.repos.Device.GetById(ctx, transaction.DeviceId)
	if err != nil {
		log.Err(err).Msg("Failed to get transaction device")
		err = libcommon.StringError(err)
		return
	}

	digitalData = eventDigitalData{
		IPAddress:         transaction.IPAddress,
		ClientFingerprint: device.Fingerprint,
	}
	return
}

func mapToUnit21TransactionEvent(transaction model.Transaction, transactionData transactionData, digitalData eventDigitalData) *u21Event {
	var transactionTagArr []string
	if transaction.Tags != nil {
		for key, value := range transaction.Tags {
			transactionTagArr = append(transactionTagArr, key+":"+value)
		}
	}

	jsonBody := &u21Event{
		GeneralData: &eventGeneral{
			EventId:      transaction.Id,                    //required
			EventType:    "transaction",                     //required
			EventTime:    int(transaction.CreatedAt.Unix()), //required
			EventSubtype: "Fiat to Crypto",                  //required for RTR
			Status:       transaction.Status,
			Parents:      nil,
			Tags:         transactionTagArr,
		},
		TransactionData: &transactionData,
		ActionData:      nil,
		DigitalData:     &digitalData,
		LocationData:    nil,
		CustomData:      nil,
	}

	return jsonBody
}
