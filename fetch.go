package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const VERSION = "1.0.0"

// ANSI color codes
const (
    Reset  = "\033[0m"
    Black  = "\033[30m"
    Red    = "\033[31m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Blue   = "\033[34m"
    Magenta= "\033[35m"
    Cyan   = "\033[36m"
    White  = "\033[37m"
)

func main() {
	verbose := false
	minimal := false

	// Check command line arguments
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-v":
			verbose = true
		case "--minimal":
			minimal = true
		}
	}
	if verbose {
		fmt.Printf("\n%sSystem Information - Version %s%s\n\n", Green, VERSION, Reset)
	}
	printSystemInfo(minimal)
}

func printSystemInfo(minimal bool) {
	fmt.Printf("%sOS:%s %s\n", Magenta, Reset, getOSInfo())
	fmt.Printf("%sPackages:%s %s\n", Yellow, Reset, getPackageInfo())
	fmt.Printf("%sKernel:%s %s\n", Yellow, Reset, getKernelInfo())
	fmt.Printf("%sUptime:%s %s\n", Yellow, Reset, getUptimeInfo())
	fmt.Printf("%sShell:%s %s\n", Yellow, Reset, getShellInfo())
	fmt.Printf("%sCPU:%s %s\n", Blue, Reset, getCPUInfo())
	fmt.Printf("%sGPU:%s %s\n", Yellow, Reset, getGPUInfo())
	fmt.Printf("%sStorage:%s %s\n", Yellow, Reset, getStorageInfo())
	fmt.Printf("%sRAM:%s %s\n", Yellow, Reset, getRAMInfo())
	fmt.Printf("%sWM:%s %s\n", Cyan, Reset, getWMInfo())
        fmt.Printf("%sDesktop Environment:%s %s\n", Red, Reset, getDesktopEnvironment())
}

func getOSInfo() string {
	cmd := exec.Command("lsb_release", "-ds")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return strings.TrimSpace(string(out))
}

func getPackageInfo() string {
	var cmd *exec.Cmd

	// Check for available package managers in order of priority
	if _, err := exec.LookPath("dpkg"); err == nil {
		cmd = exec.Command("dpkg", "--list")
	} else if _, err := exec.LookPath("rpm"); err == nil {
		cmd = exec.Command("rpm", "-qa")
	} else if _, err := exec.LookPath("dnf"); err == nil {
		cmd = exec.Command("dnf", "list", "--installed")
	} else if _, err := exec.LookPath("pacman"); err == nil {
		cmd = exec.Command("pacman", "-Q")
	} else if _, err := exec.LookPath("xbps-query"); err == nil {
		cmd = exec.Command("xbps-query", "-l")
	} else if _, err := exec.LookPath("nix-env"); err == nil {
		cmd = exec.Command("nix-env", "--query", "--installed", "--out-path")
	} else if _, err := exec.LookPath("apk"); err == nil {
		cmd = exec.Command("apk", "info")
	} else if _, err := exec.LookPath("pkg"); err == nil {
		cmd = exec.Command("pkg", "info")
	} else {
		// If no recognized package manager is found, return "Unknown"
		return fmt.Sprintf("%sPackages:%s Unknown", Yellow, Reset)
	}

	out, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("%sPackages:%s Unknown", Yellow, Reset)
	}

	// Count the number of lines in the output to get the package count
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	packageCount := len(lines)

	return fmt.Sprintf("%sPackages:%s %d", Yellow, Reset, packageCount)
}

func getKernelInfo() string {
	cmd := exec.Command("uname", "-r")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return strings.TrimSpace(string(out))
}

func getUptimeInfo() string {
	return fmt.Sprintf("Uptime: %s", getUptime())
}

func getUptime() string {
	cmd := exec.Command("uptime", "-p")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return strings.TrimSpace(string(out))
}

func getShellInfo() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "Unknown"
	}

	return "Shell: " + shell
}

func getCPUInfo() string {
	cmd := exec.Command("sh", "-c", "lscpu | grep 'Model name' | awk -F':' '{print $2}'")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return "CPU: " + strings.TrimSpace(string(out))
}

func getGPUInfo() string {
	cmd := exec.Command("sh", "-c", "lspci | grep 'VGA' | awk -F':' '{print $3}'")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return "GPU: " + strings.TrimSpace(string(out))
}

func getStorageInfo() string {
	cmd := exec.Command("df", "-h", "--output=size,used,avail,pcent,target", "/")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return "Storage: " + strings.TrimSpace(string(out))
}

func getRAMInfo() string {
	cmd := exec.Command("sh", "-c", "free -h | grep 'Mem' | awk '{print $2, $3, $4, $7}'")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return "RAM: " + strings.TrimSpace(string(out))
}

func getWMInfo() string {
	// Try to get window manager information using XDG_SESSION_TYPE
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	if sessionType != "" {
		return "Window Manager: " + sessionType
	}

	// Try to get window manager information using xprop
	cmd := exec.Command("xprop", "-root", "_NET_WM_NAME")
	out, err := cmd.Output()
	if err == nil {
		// Parse the output to extract the window manager name
		if parts := strings.SplitN(string(out), "=", 2); len(parts) == 2 {
			return "Window Manager: " + strings.TrimSpace(parts[1])
		}
	}

	return "Unknown"
}

func getDesktopEnvironment() string {
    // For Linux
    desktopEnv := os.Getenv("XDG_CURRENT_DESKTOP")
    if desktopEnv != "" {
	return desktopEnv
    }
    // Try an alternative environment variable
    altDesktopEnv := os.Getenv("DESKTOP_SESSION")
    if altDesktopEnv != "" {
	return altDesktopEnv
    }
    return "Unknown"
}
