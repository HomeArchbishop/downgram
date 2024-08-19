package tg

import (
	"context"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"golang.org/x/net/proxy"
)

const ClientStop = true
const ClientContinue = false

var (
	isClientRunning = false
	clientCtx       *context.Context
	client          *telegram.Client

	TaskQueue = make(chan func() bool, 100)

	sessionStoragePath string
)

func StartTelegramClient(appID int, appHash string, address string, auth *proxy.Auth, successChan chan bool) error {
	sock5, _ := proxy.SOCKS5("tcp", address, auth, proxy.Direct)
	dc := sock5.(proxy.ContextDialer)

	ctx := context.Background()

	cl := telegram.NewClient(appID, appHash, telegram.Options{
		Resolver: dcs.Plain(dcs.PlainOptions{
			Dial: dc.DialContext,
		}),
		SessionStorage: &telegram.FileSessionStorage{Path: sessionStoragePath},
	})

	// Run client
	if err := cl.Run(ctx, func(ctx context.Context) error {
		log.Print("Client running")
		isClientRunning = true
		clientCtx = &ctx
		client = cl

		successChan <- true

		for task := range TaskQueue {
			if task() == ClientStop {
				break
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// Close client
	log.Print("Close client")
	isClientRunning = false

	return nil
}

func IsClientRunning() bool {
	return isClientRunning
}

func init() {
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	executableDir := filepath.Dir(executablePath)
	sessionStoragePath = path.Join(executableDir, "session.json")
}
