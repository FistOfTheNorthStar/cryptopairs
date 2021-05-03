package fetcher

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/FistOfTheNorthStar/cryptopairs/internal/config"
)

type Binance struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol                     string        `json:"symbol"`
		Status                     string        `json:"status"`
		BaseAsset                  string        `json:"baseAsset"`
		BaseAssetPrecision         int           `json:"baseAssetPrecision"`
		QuoteAsset                 string        `json:"quoteAsset"`
		QuotePrecision             int           `json:"quotePrecision"`
		QuoteAssetPrecision        int           `json:"quoteAssetPrecision"`
		BaseCommissionPrecision    int           `json:"baseCommissionPrecision"`
		QuoteCommissionPrecision   int           `json:"quoteCommissionPrecision"`
		OrderTypes                 []string      `json:"orderTypes"`
		IcebergAllowed             bool          `json:"icebergAllowed"`
		OcoAllowed                 bool          `json:"ocoAllowed"`
		QuoteOrderQtyMarketAllowed bool          `json:"quoteOrderQtyMarketAllowed"`
		IsSpotTradingAllowed       bool          `json:"isSpotTradingAllowed"`
		IsMarginTradingAllowed     bool          `json:"isMarginTradingAllowed"`
		Filters                    []interface{} `json:"filters"`
		Permissions                []string      `json:"permissions"`
	} `json:"symbols"`
}

type Coinbase struct {
	Level  string  `json:"level"`
	Ts     float64 `json:"ts"`
	Caller string  `json:"caller"`
	Msg    string  `json:"msg"`
}

func FetchSymbols(conf *config.ConfigYaml) error {

	var wg sync.WaitGroup
	cErr := make(chan error)

	for _, endPoint := range conf.ApiEndPoints {
		wg.Add(1)
		go worker(cErr, endPoint, conf.CryptoFileDir, &wg)

	}

	wg.Wait()

	// this could be given some exit code and run on separate thread, but ran out of time

	for e := range cErr {
		if e != nil {
			return e
		}
	}

	return nil
}

func worker(c chan error, site string, cryptDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(site)
	if err != nil {
		c <- err
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c <- err
		return
	}

	switch {

	case strings.Contains(site, "binance"):

		err = binanceWriter(body, cryptDir)
		if err != nil {
			c <- err
			return
		}

	case strings.Contains(site, "coinbase"):

		err = coinBaseWriter(body, cryptDir)
		if err != nil {
			c <- err
			return
		}

	case strings.Contains(site, "ftx"):

	default:
		c <- errors.New("site was not matched")
	}

}

func binanceWriter(bodyData []byte, cryptDir string) error {

	if _, err := os.Stat(cryptDir + "/binance.txt"); err == nil {
		return nil
	}

	f, err := os.OpenFile(cryptDir+"/binance.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var binanceVar Binance
	err = json.Unmarshal(bodyData, &binanceVar)
	if err != nil {

		return err
	}

	for _, val := range binanceVar.Symbols {
		if _, err := f.WriteString(val.Symbol + "/" + val.QuoteAsset + "\n"); err != nil {
			return err
		}
	}

	return nil

}

func coinBaseWriter(bodyData []byte, cryptDir string) error {

	f, err := os.OpenFile(cryptDir+"/coinbase.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var decoded []interface{}

	err = json.Unmarshal(bodyData, &decoded)

	if err != nil {
		return err
	}

	for _, s := range decoded {
		sp, _ := s.(map[string]interface{})
		if _, err := f.WriteString(sp["display_name"].(string) + "\n"); err != nil {
			return err
		}
	}

	return nil

}
