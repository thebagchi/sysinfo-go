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

	info, err := GetSystemInformation()
	if nil != err {
		t.Error(err)
	}
	data, err = json.MarshalIndent(info, "", "    ")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(string(data))

	load, err := GetLoadAvg()
	if nil != err {
		t.Error(err)
	}

	data, err = json.MarshalIndent(load, "", "    ")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(string(data))

	cpu, err := GetCPUInfo()
	if nil != err {
		t.Error(err)
	}

	data, err = json.MarshalIndent(cpu, "", "    ")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(string(data))

	mem, err := GetMemInfo()
	if nil != err {
		t.Error(err)
	}

	data, err = json.MarshalIndent(mem, "", "    ")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(string(data))

	stat, err := GetStat()
	if nil != err {
		t.Error(err)
	}

	data, err = json.MarshalIndent(stat, "", "    ")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(string(data))

	vm, err := GetVmStat()
	if nil != err {
		t.Error(err)
	}

	data, err = json.MarshalIndent(vm, "", "    ")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(string(data))

	uptime, err := GetUptime()
	if nil != err {
		t.Error(err)
	}

	data, err = json.MarshalIndent(uptime, "", "    ")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(string(data))

	network, err := GetNetworkStats()
	if nil != err {
		t.Error(err)
	}

	data, err = json.MarshalIndent(network, "", "    ")
	if nil != err {
		t.Error(err)
	}
	fmt.Println(string(data))
}
