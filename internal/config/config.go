package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

type ConfigYaml struct {
	CryptoFileDir    string `yaml:"pairs_directory"`
	UpdateTim        int    `yaml:"update_time"`
	BinanceEndpoint  string `yaml:"binance_endpoint"`
	CoinBaseEndpoint string `yaml:"coinbase_endpoint"`
	FtxEndpoint      string `yaml:"ftx_enpoint"`
	KrakenEndpoint   string `yaml:"kraken_endpoint"`
	HuobiEndpoint    string `yaml:"huobi_endpoint"`

	ApiEndPoints []string
}

func GetConfigYaml(projectDir string) (yamlVar *ConfigYaml, err error) {

	var yamlBytes []byte
	if yamlBytes, err = ioutil.ReadFile(projectDir + "/default.yaml"); err != nil {
		return nil, err
	}

	yamlVar = &ConfigYaml{}
	err = yaml.Unmarshal(yamlBytes, yamlVar)
	if err != nil {
		return nil, err
	}

	GetApiEndPoints(yamlVar)

	return yamlVar, nil
}

func CreateCryptoDir(dir string) (err error) {

	_, errDir := os.Stat(dir)
	if os.IsNotExist(errDir) {

		err := os.MkdirAll(dir, 0755)
		if err != nil {

			//not possible to create dir
			return err
		}
	} else if os.IsExist(errDir) {

		//if anything needs to be done if exists
	} else {

		err = errDir
		return err
	}
	return nil
}

func GetApiEndPoints(yamlVar *ConfigYaml) {

	v := reflect.Indirect(reflect.ValueOf(yamlVar))
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if strings.Contains(typeOfS.Field(i).Name, "Endpoint") {
			if v.Field(i).Interface() != nil {
				yamlVar.ApiEndPoints = append(yamlVar.ApiEndPoints, fmt.Sprintf("%v", v.Field(i).Interface()))
			}
		}
	}
}
