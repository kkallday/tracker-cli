package application

import "os"

func SetGetenv(f func(key string) string) {
	getenv = f
}

func ResetGetenv() {
	getenv = os.Getenv
}
