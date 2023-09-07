package file

import (
	"os"
)

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsExist(err)
}
