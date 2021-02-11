package sysinfo_go

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetNetworkInterfaces(t *testing.T) {
	interfaces, err := GetNetworkInterface()
	if nil != err {
		t.Error(err)
	}
	data, err := json.MarshalIndent(interfaces, "", "    ")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(string(data))
}