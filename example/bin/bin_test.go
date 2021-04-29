package bin

import "testing"

func Test(t *testing.T) {
	aa := a() + a()
	if aa != "aa" {
		t.Error("Eror A", aa)
	}
}
