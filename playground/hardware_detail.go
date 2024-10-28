package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Report struct {
	Topic   string
	Content string
}

func getSystemInfo() []Report {
	info := make([]Report, 5) // 5 sections: Host, CPU, Memory, Disk, Network
	info[0].Topic = "Host Information"
	info[1].Topic = "CPU Information"
	info[2].Topic = "Memory Information"
	info[3].Topic = "Disk Information"
	info[4].Topic = "Network Interfaces"
	// Get Host Information
	if hostInfo, err := host.Info(); err == nil {
		info[0].Content = fmt.Sprintf(
			"Hostname: %v\n"+
				"OS: %v\n"+
				"Platform: %v\n"+
				"Platform Version: %v\n"+
				"Kernel Version: %v\n"+
				"Uptime: %v hours",
			hostInfo.Hostname,
			hostInfo.OS,
			hostInfo.Platform,
			hostInfo.PlatformVersion,
			hostInfo.KernelVersion,
			hostInfo.Uptime/3600)
	} else {
		info[0].Content = "Failed to get host information"
	}

	// Get CPU Information
	if cpuInfo, err := cpu.Info(); err == nil {
		var cpuDetails strings.Builder
		for _, cpu := range cpuInfo {
			fmt.Fprintf(&cpuDetails, "CPU Model: %v\n"+
				"Cores: %v\n"+
				"MHz: %v\n",
				cpu.ModelName,
				cpu.Cores,
				cpu.Mhz)
		}
		info[1].Content = cpuDetails.String()
	} else {
		info[1].Content = "Failed to get CPU information"
	}

	// Get Memory Information
	if memInfo, err := mem.VirtualMemory(); err == nil {
		info[2].Content = fmt.Sprintf(
			"Total: %v GB\n"+
				"Free: %v GB\n"+
				"Used: %v GB\n"+
				"Usage: %.2f%%",
			memInfo.Total/1024/1024/1024,
			memInfo.Free/1024/1024/1024,
			(memInfo.Total-memInfo.Free)/1024/1024/1024,
			memInfo.UsedPercent)
	} else {
		info[2].Content = "Failed to get memory information"
	}

	// Get Disk Information
	if partitions, err := disk.Partitions(false); err == nil {
		var diskDetails strings.Builder
		for _, partition := range partitions {
			if diskUsage, err := disk.Usage(partition.Mountpoint); err == nil {
				fmt.Fprintf(&diskDetails, "\nMount Point: %v\n"+
					"Total: %v GB\n"+
					"Free: %v GB\n"+
					"Used: %v GB\n"+
					"Usage: %.2f%%\n",
					partition.Mountpoint,
					diskUsage.Total/1024/1024/1024,
					diskUsage.Free/1024/1024/1024,
					diskUsage.Used/1024/1024/1024,
					diskUsage.UsedPercent)
			}
		}
		info[3].Content = diskDetails.String()
	} else {
		info[3].Content = "Failed to get disk information"
	}

	// Get Network Interfaces
	if interfaces, err := net.Interfaces(); err == nil {
		var netDetails strings.Builder

		for _, iface := range interfaces {
			fmt.Fprintf(&netDetails, "\nInterface: %v\n"+
				"Hardware Address: %v\n"+
				"MTU: %v\n"+
				"Flags: %v\n",
				iface.Name,
				iface.HardwareAddr,
				iface.MTU,
				iface.Flags)

			if addrs, err := iface.Addrs(); err == nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}

					if ip.To4() != nil {
						fmt.Fprintf(&netDetails, "IPv4 Address: %v\n", ip)
					} else {
						fmt.Fprintf(&netDetails, "IPv6 Address: %v\n", ip)
					}
				}
			}
		}
		info[4].Content = netDetails.String()
	} else {
		info[4].Content = "Failed to get network information"
	}

	return info
}

func main() {
	systemInfo := getSystemInfo()

	// Print all information with index numbers
	for i, info := range systemInfo {
		fmt.Println(fmt.Sprintf("INDEX %d TOPIC %s", i, info.Topic))
	}

	// Example of accessing specific information
	fmt.Println("\nAccessing specific information:")
	fmt.Printf("Host Information :\n%s\n", systemInfo[0].Content)
}
