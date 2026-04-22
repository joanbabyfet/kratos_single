package i18n

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func InitI18n() {
	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	//Linux / Mac / Windows 都可用
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "configs", "lang", "zh.json")

	_, err := bundle.LoadMessageFile(path)
	if err != nil {
		log.Fatal(err)
	}

	path = filepath.Join(wd, "configs", "lang", "en.json")
	_, err = bundle.LoadMessageFile(path)
	if err != nil {
		log.Fatal(err)
	}
}

func T(lang string, key string) string {
	localizer := i18n.NewLocalizer(bundle, lang)
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	return msg
}