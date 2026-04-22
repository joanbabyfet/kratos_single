package i18n

import (
	"encoding/json"
	"kratos_single/internal/pkg/utils"
	"log"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func InitI18n() {
	bundle = i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	//Linux / Mac / Windows 都可用
	root := utils.RootPath()
	path := filepath.Join(root, "configs", "lang", "zh.json")

	_, err := bundle.LoadMessageFile(path)
	if err != nil {
		log.Fatal(err)
	}

	path = filepath.Join(root, "configs", "lang", "en.json")
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