package xtime

import "testing"

func TestTomorrowStartAt(t *testing.T) {
	t.Logf("明天凌晨的时间戳是: %d", TomorrowStartAt())
}
