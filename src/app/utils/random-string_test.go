package utils

import (
	"strings"
	"testing"
)

func Test_Rendom_String_Length(t *testing.T) {
	str := RandomString(8)

	if len(str) != 8 {
		t.Errorf("Random string length must be 8 characters")
	}
}

func Test_Random_String_Contained_Chars(t *testing.T) {
	str := RandomString(10)

	if !strings.ContainsAny(str, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTYVWXYZ0123456789") {
		t.Fail()
	}
}
