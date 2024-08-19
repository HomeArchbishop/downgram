package loginview

import (
	"log"
	"net/url"
	"strconv"

	"github.com/homearchbishop/downgram/internal/cache"
	"github.com/homearchbishop/downgram/internal/tg"
	"golang.org/x/net/proxy"
)

type loginInfo struct {
	appID    string
	apphash  string
	phone    string
	proxyStr string
	newLogin bool
}

var (
	isLogining = false
)

func hasScheme(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != ""
}

func login(gci loginInfo, askForCode func() string, successCb func()) {
	if isLogining {
		return
	}
	isLogining = true

	if tg.IsClientRunning() {
		isLogining = false
		return
	}

	address := ""
	auth := &proxy.Auth{}
	func() {
		proxyStr := gci.proxyStr
		if !hasScheme(proxyStr) {
			proxyStr = "socks5://" + proxyStr
		}
		parsedURL, err := url.Parse(proxyStr)
		if err != nil {
			return
		}

		address = parsedURL.Host

		if parsedURL.User != nil {
			auth.User = parsedURL.User.Username()
			password, _ := parsedURL.User.Password()
			auth.Password = password
		}
	}()

	intAppID, _ := strconv.Atoi(gci.appID)

	// Start telegram client, nothing to do with the Authentification
	successChan := make(chan bool)
	go func() {
		if gci.newLogin {
			tg.RemoveSession()
		}
		err := tg.StartTelegramClient(intAppID, gci.apphash, address, auth, successChan)
		if err != nil {
			isLogining = false
			return
		}
	}()

	go func() {
		if <-successChan { // Wait for the client to start
			tg.TaskQueue <- func() bool {
				isAuth, err := tg.GetStatus()
				if err != nil {
					log.Print(err)
					isLogining = false
					return tg.ClientStop
				}
				if !isAuth { // use code-only login
					phoneCodeHash, err := tg.GetCode(gci.phone)
					if err != nil {
						log.Print(err)
						isLogining = false
						return tg.ClientStop
					}
					code := askForCode()
					if err := tg.LoginWithCode(gci.phone, code, phoneCodeHash); err != nil {
						return tg.ClientStop
					}
				}
				isLogining = false
				cache.Update("appid", gci.appID)
				cache.Update("apphash", gci.apphash)
				cache.Update("phone", gci.phone)
				cache.Update("proxy", gci.proxyStr)

				successCb()
				return tg.ClientContinue
			}
		}
	}()
}
