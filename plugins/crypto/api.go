package crypto

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

const endpoint = "https://api.cryptowat.ch/"

func get(request string) ([]byte, error) {
	r, err := http.Get(endpoint + request)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

type Response struct {
	Allowance Allowance
	Result    interface{}
}

func Get(request string) (*Response, error) {
	b, err := get(request)
	if err != nil {
		return nil, err
	}
	var x Response
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type Allowance struct {
	Cost      int64
	Remaining int64
}

type AssetsResponse struct {
	Allowance Allowance
	Result    AssetsResult
}

func GetAssets() (*AssetsResponse, error) {
	b, err := get("assets")
	if err != nil {
		return nil, err
	}
	var x AssetsResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type AssetsResult []Asset

type Asset struct {
	Symbol string
	Name   string
	Fiat   bool
	Route  string
}

type AssetResponse struct {
	Allowance Allowance
	Result    AssetResult
}

func GetAsset(currency string) (*AssetResponse, error) {
	b, err := get("assets/" + currency)
	if err != nil {
		return nil, err
	}
	var x AssetResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type AssetResult struct {
	Symbol  string
	Name    string
	Fiat    bool
	Markets struct {
		Base  []Market
		Quote []Market
	}
}

type Market struct {
	Exchange string
	Pair     string
	Active   bool
	Route    string
}

type PairsResponse struct {
	Allowance Allowance
	Result    PairsResult
}

func GetPairs() (*PairsResponse, error) {
	b, err := get("pairs")
	if err != nil {
		return nil, err
	}
	var x PairsResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type PairsResult []struct {
	Symbol string
	Base   Asset
	Quote  Asset
	Route  string
}

type PairResponse struct {
	Allowance Allowance
	Result    PairResult
}

func GetPair(pair string) (*PairResponse, error) {
	b, err := get("pairs/" + pair)
	if err != nil {
		return nil, err
	}
	var x PairResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type PairResult struct {
	Symbol  string
	Base    Asset
	Quote   Asset
	Route   string
	Markets []Market
}

type ExchangesResponse struct {
	Allowance Allowance
	Result    ExchangesResult
}

func GetExchanges() (*ExchangesResponse, error) {
	b, err := get("exchanges")
	if err != nil {
		return nil, err
	}
	var x ExchangesResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type ExchangesResult []struct {
	Symbol string
	Name   string
	Active bool
	Route  string
}

type ExchangeResponse struct {
	Allowance Allowance
	Result    ExchangeResult
}

func GetExchange(exchange string) (*ExchangeResponse, error) {
	b, err := get("exchanges/" + exchange)
	if err != nil {
		return nil, err
	}
	var x ExchangeResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type ExchangeResult struct {
	Symbol string
	Name   string
	Active bool
	Routes struct {
		Markets string
	}
}

type MarketsResponse struct {
	Allowance Allowance
	Result    MarketsResult
}

func GetMarkets(exchange string) (*MarketsResponse, error) {
	resource := "markets"
	if exchange != "" {
		resource += "/" + exchange
	}
	b, err := get(resource)
	if err != nil {
		return nil, err
	}
	var x MarketsResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketsResult []Market

type MarketResponse struct {
	Allowance Allowance
	Result    MarketResult
}

func GetMarket(exchange, pair string) (*MarketResponse, error) {
	b, err := get("markets/" + exchange + "/" + pair)
	if err != nil {
		return nil, err
	}
	var x MarketResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketResult struct {
	Exchange string
	Pair     string
	Active   bool
	Routes   struct {
		Price     string
		Summary   string
		Orderbook string
		Trades    string
		OHLC      string
	}
}

type MarketPriceResponse struct {
	Allowance Allowance
	Result    MarketPriceResult
}

func GetMarketPrice(exchange, pair string) (*MarketPriceResponse, error) {
	b, err := get("markets/" + exchange + "/" + pair + "/price")
	if err != nil {
		return nil, err
	}
	var x MarketPriceResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketPriceResult struct {
	Price float64
}

type MarketSummaryResponse struct {
	Allowance Allowance
	Result    MarketSummaryResult
}

func GetMarketSummary(exchange, pair string) (*MarketSummaryResponse, error) {
	b, err := get("markets/" + exchange + "/" + pair + "/summary")
	if err != nil {
		return nil, err
	}
	var x MarketSummaryResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketSummaryResult Summary

type Summary struct {
	Price struct {
		Last   float64
		High   float64
		Low    float64
		Change struct {
			Percentage float64
			Absolute   float64
		}
	}
	Volume float64
}

type MarketOrderbookResponse struct {
	Allowance Allowance
	Result    MarketOrderbookResult
}

func GetMarketOrderbook(exchange, pair string) (*MarketOrderbookResponse, error) {
	b, err := get("markets/" + exchange + "/" + pair + "/orderbook")
	if err != nil {
		return nil, err
	}
	var x MarketOrderbookResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketOrderbookResult struct {
	Asks []Order
	Bids []Order
}

type Order struct {
	Price  float64
	Amount float64
}

func (o *Order) UnmarshalJSON(b []byte) error {
	var x [2]float64
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	o.Price = x[0]
	o.Amount = x[1]
	return nil
}

type MarketTradesResponse struct {
	Allowance Allowance
	Result    MarketTradesResult
}

func GetMarketTrades(exchange, pair string, options MarketTradesOptions) (*MarketTradesResponse, error) {
	b, err := get("markets/" + exchange + "/" + pair + "/trades")
	if err != nil {
		return nil, err
	}
	var x MarketTradesResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketTradesOptions struct {
	Limit int
	Since int64
}

func GetMarketTradesOptions(exchange, pair string, options MarketTradesOptions) (*MarketTradesResponse, error) {
	resource := "markets/" + exchange + "/" + pair + "/trades"
	if options.Limit != 0 {
		resource += "?limit=" + strconv.FormatInt(int64(options.Limit), 10)
	}
	if options.Since != 0 {
		resource += "?since=" + strconv.FormatInt(options.Since, 10)
	}
	b, err := get(resource)
	if err != nil {
		return nil, err
	}
	var x MarketTradesResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketTradesResult []MarketTrade

type MarketTrade struct {
	ID        int
	Timestamp int64
	Price     float64
	Amount    float64
}

func (mt *MarketTrade) UnmarshalJSON(b []byte) error {
	var x [4]float64
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	mt.ID = int(x[0])
	mt.Timestamp = int64(x[1])
	mt.Price = x[2]
	mt.Amount = x[3]
	return nil
}

type MarketOHLCResponse struct {
	Allowance Allowance
	Result    MarketOHLCResult
}

func GetMarketOHLC(exchange, pair string) (*MarketOHLCResponse, error) {
	b, err := get("markets/" + exchange + "/" + pair + "/ohcl")
	if err != nil {
		return nil, err
	}
	var x MarketOHLCResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketOHLCOptions struct {
	Before, After int64
	Periods       string
}

func GetMarketOHLCOptions(exchange, pair string, options MarketOHLCOptions) (*MarketOHLCResponse, error) {
	resource := "markets/" + exchange + "/" + pair + "/ohcl"
	if options.Before != 0 {
		resource += "?before=" + strconv.FormatInt(options.Before, 10)
	}
	if options.After != 0 {
		resource += "?after=" + strconv.FormatInt(options.After, 10)
	}
	if options.Periods != "" {
		resource += "?periods=" + options.Periods
	}
	b, err := get(resource)
	if err != nil {
		return nil, err
	}
	var x MarketOHLCResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketOHLCResult map[string][]Candle

type Candle struct {
	CloseTime  int64
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
	Volume     float64
}

func (c *Candle) UnmarshalJSON(b []byte) error {
	var x [6]float64
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	c.CloseTime = int64(x[0])
	c.OpenPrice = x[1]
	c.HighPrice = x[2]
	c.LowPrice = x[3]
	c.ClosePrice = x[4]
	c.Volume = x[5]
	return nil
}

type MarketsPricesResponse struct {
	Allowance Allowance
	Result    MarketsPricesResult
}

func GetMarketsPrices() (*MarketsPricesResponse, error) {
	b, err := get("markets/prices")
	if err != nil {
		return nil, err
	}
	var x MarketsPricesResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketsPricesResult map[string]float64

type MarketsSummariesResponse struct {
	Allowance Allowance
	Result    MarketsSummariesResult
}

func GetMarketsSummaries() (*MarketsSummariesResponse, error) {
	b, err := get("markets/summaries")
	if err != nil {
		return nil, err
	}
	var x MarketsSummariesResponse
	err = json.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

type MarketsSummariesResult map[string]Summary
