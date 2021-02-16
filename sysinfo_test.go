package sysinfo_go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"syscall"
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

	processes, err := ListProcessId()
	if nil != err {
		t.Error(err)
	}

	fmt.Println("Number of processes:", len(processes))

	{
		content, err := ioutil.ReadFile(InterruptFile)
		if nil != err {
			t.Error(err)
		}
		fmt.Println(FastBytesToString(content))
	}

	{
		content, err := ioutil.ReadFile(DiskStatFile)
		if nil != err {
			t.Error(err)
		}
		fmt.Println(FastBytesToString(content))
	}

	{
		uname := new(syscall.Utsname)
		if err := syscall.Uname(uname); err != nil {
			t.Error(err)
		}
		fmt.Println(uname)
		builder := new(strings.Builder)
		for _, it := range uname.Machine {
			if it != 0 {
				builder.WriteByte(byte(it))
			}
		}
		fmt.Println(builder.String())
		builder.Reset()
		for _, it := range uname.Domainname {
			if it != 0 {
				builder.WriteByte(byte(it))
			}
		}
		fmt.Println(builder.String())
		builder.Reset()
		for _, it := range uname.Nodename {
			if it != 0 {
				builder.WriteByte(byte(it))
			}
		}
		fmt.Println(builder.String())
		builder.Reset()
		for _, it := range uname.Release {
			if it != 0 {
				builder.WriteByte(byte(it))
			}
		}
		fmt.Println(builder.String())
		builder.Reset()
		for _, it := range uname.Sysname {
			if it != 0 {
				builder.WriteByte(byte(it))
			}
		}
		fmt.Println(builder.String())
		builder.Reset()
		for _, it := range uname.Version {
			if it != 0 {
				builder.WriteByte(byte(it))
			}
		}
		fmt.Println(builder.String())
	}
}
