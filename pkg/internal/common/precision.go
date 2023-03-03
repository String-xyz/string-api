package common

import (
	"strconv"

	"github.com/String-xyz/string-api/pkg/model"
)

func QuoteToPrecise(imprecise model.Quote) model.PrecisionSafeQuote {
	res := model.PrecisionSafeQuote{
		Timestamp:  imprecise.Timestamp,
		BaseUSD:    strconv.FormatFloat(imprecise.BaseUSD, 'G', -1, 64),
		GasUSD:     strconv.FormatFloat(imprecise.GasUSD, 'G', -1, 64),
		TokenUSD:   strconv.FormatFloat(imprecise.TokenUSD, 'G', -1, 64),
		ServiceUSD: strconv.FormatFloat(imprecise.ServiceUSD, 'G', -1, 64),
		TotalUSD:   strconv.FormatFloat(imprecise.TotalUSD, 'G', -1, 64),
	}
	return res
}

func QuoteToImprecise(precise model.PrecisionSafeQuote) model.Quote {
	res := model.Quote{
		Timestamp: precise.Timestamp,
	}
	res.BaseUSD, _ = strconv.ParseFloat(precise.BaseUSD, 64)
	res.GasUSD, _ = strconv.ParseFloat(precise.GasUSD, 64)
	res.TokenUSD, _ = strconv.ParseFloat(precise.TokenUSD, 64)
	res.ServiceUSD, _ = strconv.ParseFloat(precise.ServiceUSD, 64)
	res.TotalUSD, _ = strconv.ParseFloat(precise.TotalUSD, 64)
	return res
}

func ExecutionRequestToPrecise(imprecise model.ExecutionRequest) model.PrecisionSafeExecutionRequest {
	res := model.PrecisionSafeExecutionRequest{
		TransactionRequest: imprecise.TransactionRequest,
		PrecisionSafeQuote: QuoteToPrecise(imprecise.Quote),
		Signature:          imprecise.Signature,
		CardToken:          imprecise.CardToken,
	}
	return res
}

func ExecutionRequestToImprecise(precise model.PrecisionSafeExecutionRequest) model.ExecutionRequest {
	res := model.ExecutionRequest{
		TransactionRequest: precise.TransactionRequest,
		Quote:              QuoteToImprecise(precise.PrecisionSafeQuote),
		Signature:          precise.Signature,
		CardToken:          precise.CardToken,
	}
	return res
}
