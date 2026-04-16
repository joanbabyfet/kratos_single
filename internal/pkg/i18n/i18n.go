package i18n

import (
	"encoding/json"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func InitI18n() {
	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	_, err := bundle.LoadMessageFile("../../configs/lang/zh.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = bundle.LoadMessageFile("../../configs/lang/en.json")
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