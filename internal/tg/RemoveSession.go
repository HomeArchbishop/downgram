package tg

import "os"

func RemoveSession() {
	os.Remove(sessionStoragePath)
}
