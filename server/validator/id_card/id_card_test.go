package id_card

import (
	"fmt"
	"testing"
)

func TestIdCard(t *testing.T) {
	citizenNo := []byte("430528196402101068")
	// citizenNo := []byte("340321199001234560")
	// citizenNo := []byte("340321900123456")
	birthday, isMale, addrMask, err := GetCitizenNoInfo(citizenNo)
	if err != nil {
		fmt.Println("Invalid citizen number.")
	} else {
		fmt.Println("Valid citizen number.")
		fmt.Printf("Information from citizen: birthday=%v, ismale=%v, addrmask=%d\n", birthday, isMale, addrMask)
	}
}
