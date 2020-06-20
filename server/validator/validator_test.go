package validator

import (
	"regexp"
	"testing"
)

func isMobile(mobile string) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(mobile)
}

func isIdCardNo(no string) bool {
	reg := "^[1-9]\\d{7}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}$|^[1-9]\\d{5}[1-9]\\d{3}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}([0-9]|X)$"
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(no)
}

func TestAll(t *testing.T) {
	testCheckMobile(t)
	testCheckIdCardNo(t)
}

func testCheckMobile(t *testing.T) {
	m1 := "18218871255"
	m2 := "11122299912"
	if isMobile(m1) {
		t.Log("done!")
	} else {
		t.Error("failed")
	}
	if !isMobile(m2) {
		t.Log("done!")
	} else {
		t.Error("failed")
	}

}

func testCheckIdCardNo(t *testing.T) {
	no1 := "440883199201211152"
	no2 := "12345678945613"
	if isIdCardNo(no1) {
		t.Log("done!")
	} else {
		t.Error("failed")
	}
	if !isIdCardNo(no2) {
		t.Log("done!")
	} else {
		t.Error("failed")
	}
}
