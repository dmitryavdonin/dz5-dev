package tokenManager

import (
	"errors"
	"os"
)

func init() {
	//create folder for cert if not exists
	if _, err := os.Stat(jwksPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(jwksPath, os.ModePerm)
		if err != nil {
			return
		}
	}
}
