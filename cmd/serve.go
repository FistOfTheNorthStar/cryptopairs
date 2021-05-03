package cmd

import (
	"path"
	"path/filepath"
	"runtime"

	"github.com/FistOfTheNorthStar/cryptopairs/internal/config"
	"github.com/FistOfTheNorthStar/cryptopairs/internal/fetcher"
	"github.com/FistOfTheNorthStar/cryptopairs/internal/global"
)

type Application struct {
	config *config.ConfigYaml
}

func Initialize() {

	app := Application{}

	var b string
	_, b, _, _ = runtime.Caller(0)
	dirCalled := filepath.Dir(b)
	projectDir := path.Clean(path.Join(path.Dir(dirCalled), "data"))

	if conf, err := config.GetConfigYaml(projectDir); err == nil {
		app.config = conf
	} else {
		global.Log.Fatalf("Application conf to start: ", err.Error())
	}

	err := config.CreateCryptoDir(app.config.CryptoFileDir)

	if err != nil {
		global.Log.Fatalf("Failed to create directory for pair files: ", err.Error())
	}

	err = fetcher.FetchSymbols(app.config)

	if err != nil {
		global.Log.Fatalf("Failed to fetch symbols: ", err.Error())
	}

}
