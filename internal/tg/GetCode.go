package tg

import (
	"log"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

func GetCode(phone string) (string, error) {
	log.Print("Get code")
	resp, err := client.Auth().SendCode(*clientCtx, phone, auth.SendCodeOptions{})
	if err != nil {
		return "", err
	}

	phoneCodeHash := resp.(*tg.AuthSentCode).PhoneCodeHash

	return phoneCodeHash, nil
}
