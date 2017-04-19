package myerror

import (
	"fmt"
	"os"
)

func CheckError(err error, message string) {
	if err != nil {
		panic(err)
		os.Exit(-1)
	}
}
func CheckErrorConsole(err error, message string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%s:%s", err, message))
		os.Exit(-1)
	}
}
func CheckErrorJson(err error, message string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%s:%s", err, message))
		os.Exit(-1)
	}
}
