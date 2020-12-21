package internal

import "os"

func Getenv(name, _default string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}

	return _default
}
