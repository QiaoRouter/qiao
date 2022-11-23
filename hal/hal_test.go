package hal

import (
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	Init()
}

func TestIfHandle_GetNeighborMacAddr(t *testing.T) {
	Init()

}

func TestIpv6(t *testing.T) {
	if strings.Compare(prefix("fd00::3:2/112"), "fd00::3:2") != 0 {
		t.Fatal("TestIpv6 fail.")
	}
}
