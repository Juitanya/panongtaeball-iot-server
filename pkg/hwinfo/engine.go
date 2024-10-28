package hwinfo

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"net"
	"strings"
)

type Report struct {
	Topic   string `json:"topic"`
	Content string `json:"content"`
}

// HostInfo represents host system information
type HostInfo struct {
	Hostname    string `json:"hostname"`
	OS          string `json:"os"`
	Platform    string `json:"platform"`
	PlatformVer string `json:"platform_ver"`
	KernelVer   string `json:"kernel_ver"`
	UptimeHours uint64 `json:"uptime_hours"`
}

// CPUInfo represents CPU information
type CPUInfo struct {
	Model string  `json:"model"`
	Cores int32   `json:"cores"`
	Mhz   float64 `json:"mhz"`
}

// MemoryInfo represents memory information
type MemoryInfo struct {
	TotalGB     float64 `json:"total_gb"`
	FreeGB      float64 `json:"free_gb"`
	UsedGB      float64 `json:"used_gb"`
	UsedPercent float64 `json:"used_percent"`
}

// DiskInfo represents disk partition information
type DiskInfo struct {
	MountPoint  string  `json:"mount_point"`
	TotalGB     float64 `json:"total_gb"`
	FreeGB      float64 `json:"free_gb"`
	UsedGB      float64 `json:"used_gb"`
	UsedPercent float64 `json:"used_percent"`
}

// NetworkInterface represents network interface information
type NetworkInterface struct {
	Name          string    `json:"name"`
	HardwareAddr  string    `json:"hardware_addr"`
	MTU           int       `json:"mtu"`
	Flags         net.Flags `json:"flags"`
	IPv4Addresses []string  `json:"ip_v_4_addresses"`
	IPv6Addresses []string  `json:"ipv_6_addresses"`
}

type SystemInfo struct {
	Host    HostInfo           `json:"host"`
	CPU     []CPUInfo          `json:"cpu"`
	Memory  MemoryInfo         `json:"memory"`
	Disks   []DiskInfo         `json:"disks"`
	Network []NetworkInterface `json:"network"`
}

func NewSystemInfo() (SystemInfo, error) {
	var sysInfo SystemInfo
	var err error
	err = sysInfo.FetchData()
	if err != nil {
		return SystemInfo{}, err
	}
	return sysInfo, err
}

func (si *SystemInfo) FetchData() error {
	var err error

	// Get Host Information
	if hostInfo, err := host.Info(); err == nil {
		si.Host = HostInfo{
			Hostname:    hostInfo.Hostname,
			OS:          hostInfo.OS,
			Platform:    hostInfo.Platform,
			PlatformVer: hostInfo.PlatformVersion,
			KernelVer:   hostInfo.KernelVersion,
			UptimeHours: hostInfo.Uptime / 3600,
		}
	}

	// Get CPU Information
	if cpuInfo, err := cpu.Info(); err == nil {
		for _, cpu := range cpuInfo {
			si.CPU = append(si.CPU, CPUInfo{
				Model: cpu.ModelName,
				Cores: cpu.Cores,
				Mhz:   cpu.Mhz,
			})
		}
	}

	// Get Memory Information
	if memInfo, err := mem.VirtualMemory(); err == nil {
		si.Memory = MemoryInfo{
			TotalGB:     float64(memInfo.Total) / (1024 * 1024 * 1024),
			FreeGB:      float64(memInfo.Free) / (1024 * 1024 * 1024),
			UsedGB:      float64(memInfo.Total-memInfo.Free) / (1024 * 1024 * 1024),
			UsedPercent: memInfo.UsedPercent,
		}
	}

	// Get Disk Information
	if partitions, err := disk.Partitions(false); err == nil {
		for _, partition := range partitions {
			if diskUsage, err := disk.Usage(partition.Mountpoint); err == nil {
				si.Disks = append(si.Disks, DiskInfo{
					MountPoint:  partition.Mountpoint,
					TotalGB:     float64(diskUsage.Total) / (1024 * 1024 * 1024),
					FreeGB:      float64(diskUsage.Free) / (1024 * 1024 * 1024),
					UsedGB:      float64(diskUsage.Used) / (1024 * 1024 * 1024),
					UsedPercent: diskUsage.UsedPercent,
				})
			}
		}
	}

	// Get Network Interface Information
	if interfaces, err := net.Interfaces(); err == nil {
		for _, iface := range interfaces {
			netInterface := NetworkInterface{
				Name:         iface.Name,
				HardwareAddr: iface.HardwareAddr.String(),
				MTU:          iface.MTU,
				Flags:        iface.Flags,
			}

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
						netInterface.IPv4Addresses = append(netInterface.IPv4Addresses, ip.String())
					} else {
						netInterface.IPv6Addresses = append(netInterface.IPv6Addresses, ip.String())
					}
				}
			}
			si.Network = append(si.Network, netInterface)
		}
	}
	return err
}

func (si *SystemInfo) ToReports(refresh bool) []Report {
	if refresh {
		si.FetchData()
	}

	info := make([]Report, 5) // 5 sections: Host, CPU, Memory, Disk, Network
	info[0].Topic = "Host Information"
	info[1].Topic = "CPU Information"
	info[2].Topic = "Memory Information"
	info[3].Topic = "Disk Information"
	info[4].Topic = "Network Interfaces"

	info[0].Content = fmt.Sprintf(
		"Hostname: %v\n"+
			"OS: %v\n"+
			"Platform: %v\n"+
			"Platform Version: %v\n"+
			"Kernel Version: %v\n"+
			"Uptime: %v hours",
		si.Host.Hostname,
		si.Host.OS,
		si.Host.Platform,
		si.Host.PlatformVer,
		si.Host.KernelVer,
		si.Host.UptimeHours)

	var cpuDetails strings.Builder

	fmt.Fprintf(&cpuDetails, "CPU Model (%d Threads): %v\n"+
		"Cores: %v\n"+
		"MHz: %v\n",
		len(si.CPU),
		si.CPU[0].Model,
		si.CPU[0].Cores,
		si.CPU[0].Mhz)

	info[1].Content = cpuDetails.String()

	info[2].Content = fmt.Sprintf(
		"Total: %.2f GB\n"+
			"Free: %.2f GB\n"+
			"Used: %.2f GB\n"+
			"Usage: %.2f%%",
		si.Memory.TotalGB,
		si.Memory.FreeGB,
		si.Memory.UsedGB,
		si.Memory.UsedPercent)

	var diskDetails strings.Builder
	for _, partition := range si.Disks {
		fmt.Fprintf(&diskDetails, "\nMount Point: %v\n"+
			"Total: %.2f GB\n"+
			"Free: %.2f GB\n"+
			"Used: %.2f GB\n"+
			"Usage: %.2f%%\n",
			partition.MountPoint,
			partition.TotalGB,
			partition.FreeGB,
			partition.UsedGB,
			partition.UsedPercent)

	}
	info[3].Content = diskDetails.String()

	var netDetails strings.Builder

	for _, iface := range si.Network {
		fmt.Fprintf(&netDetails, "\nInterface: %v\n"+
			"Hardware Address: %v\n"+
			"MTU: %v\n"+
			"Flags: %v\n",
			iface.Name,
			iface.HardwareAddr,
			iface.MTU,
			iface.Flags)

		for i := range iface.IPv4Addresses {
			fmt.Fprintf(&netDetails, "IPv4 Address: %v\n", iface.IPv4Addresses[i])
		}

		for i := range iface.IPv6Addresses {
			fmt.Fprintf(&netDetails, "IPv6 Address: %v\n", iface.IPv6Addresses[i])
		}
		info[4].Content = netDetails.String()
	}

	return info
}

func (si SystemInfo) PrintHostInfo(refresh bool) {
	if refresh {
		si.FetchData()
	}

	fmt.Printf("%s\n", fmt.Sprintf(
		"Hostname: %v\n"+
			"OS: %v\n"+
			"Platform: %v\n"+
			"Platform Version: %v\n"+
			"Kernel Version: %v\n"+
			"Uptime: %v hours",
		si.Host.Hostname,
		si.Host.OS,
		si.Host.Platform,
		si.Host.PlatformVer,
		si.Host.KernelVer,
		si.Host.UptimeHours))
}
