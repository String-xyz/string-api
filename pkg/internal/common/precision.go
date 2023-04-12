package common

import (
	"strconv"

	"github.com/String-xyz/string-api/pkg/model"
)

func EstimateToPrecise(imprecise model.Estimate[float64]) model.Estimate[string] {
	res := model.Estimate[string]{
		Timestamp:  imprecise.Timestamp,
		BaseUSD:    strconv.FormatFloat(imprecise.BaseUSD, 'f', 2, 64),
		GasUSD:     strconv.FormatFloat(imprecise.GasUSD, 'f', 2, 64),
		TokenUSD:   strconv.FormatFloat(imprecise.TokenUSD, 'f', 2, 64),
		ServiceUSD: strconv.FormatFloat(imprecise.ServiceUSD, 'f', 2, 64),
		TotalUSD:   strconv.FormatFloat(imprecise.TotalUSD, 'f', 2, 64),
	}
	return res
}

func EstimateToImprecise(precise model.Estimate[string]) model.Estimate[float64] {
	res := model.Estimate[float64]{
		Timestamp: precise.Timestamp,
	}
	res.BaseUSD, _ = strconv.ParseFloat(precise.BaseUSD, 64)
	res.GasUSD, _ = strconv.ParseFloat(precise.GasUSD, 64)
	res.TokenUSD, _ = strconv.ParseFloat(precise.TokenUSD, 64)
	res.ServiceUSD, _ = strconv.ParseFloat(precise.ServiceUSD, 64)
	res.TotalUSD, _ = strconv.ParseFloat(precise.TotalUSD, 64)
	return res
}

// type Quote struct {
// 	TransactionRequest TransactionRequest `json:"request"`
// 	Estimate           Estimate[string]   `json:"estimate"`
// 	Signature          string             `json:"signature"`
// }

func QuoteToPrecise(imprecise model.Quote) model.Quote {
	res := model.Quote{
		TransactionRequest: imprecise.TransactionRequest,
		Estimate:           EstimateToPrecise(imprecise.Estimate),
		Signature:          imprecise.Signature,
	}
	return res
}

func QuoteToImprecise(precise model.Quote) model.Quote {
	res := model.Quote{
		TransactionRequest: precise.TransactionRequest,
		Estimate:           EstimateToImprecise(precise.Estimate),
		Signature:          precise.Signature,
	}
	return res
}
