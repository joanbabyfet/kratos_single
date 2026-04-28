package data

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"kratos_single/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

//用途：注册成功、订单成立、付款成功 推送 Telegram
type TelegramRepo struct {
	log    *log.Helper
	conf   *conf.Data
}

func NewTelegramRepo(c *conf.Data, logger log.Logger) *TelegramRepo {
	return &TelegramRepo{
		conf: c,
		log:  log.NewHelper(logger),
	}
}

//Telegram API 回传结构
type TelegramResp struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

//发送通知
func (r *TelegramRepo) Send(ctx context.Context, message string) error {

	token := r.conf.Telegram.Token
	chatID := r.conf.Telegram.ChatId

	api := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		token,
	)

	form := url.Values{}
	form.Set("chat_id", chatID)
	form.Set("text", message)
	form.Set("parse_mode", "HTML")

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		api,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		r.log.Errorf("telegram request fail: %v", err)
		return err
	}
	defer resp.Body.Close()

	var result TelegramResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Ok {
		return fmt.Errorf("telegram api fail: %s", result.Description)
	}

	r.log.Infof("telegram send success")

	return nil
}