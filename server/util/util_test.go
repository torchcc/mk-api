package util

import (
	"testing"
)

func TestIsMobile(t *testing.T) {

	s := []string{"18505921256", "18330823069", "61915902321", ""}
	for _, v := range s {
		t.Log(IsMobile(v))
	}
}
