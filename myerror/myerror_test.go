package myerror

import (
	"errors"
	"testing"
)

func Test_error(t *testing.T) {
	CheckError(nil, "")
	CheckErrorConsole(errors.New("====================="), "sdfsdfsdfsdf")
	CheckError(errors.New("sdfsdfsdfsdf"), "sdfsdfsdfsdfsdfsdf")
}
