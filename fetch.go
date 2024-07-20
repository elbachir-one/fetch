package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const VERSION = "1.2.0"

// ANSI color codes
const (
	Reset   = "\033[0m"
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Bold    = "\033[1m" // ANSI bold code
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
		fmt.Printf("\n%sSystem Information - Version %s%s%s\n\n", Bold, Green, VERSION, Reset)
	}
	printSystemInfo(minimal)
}

func printSystemInfo(minimal bool) {
	// Define a format string with a fixed width for labels
	format := "%-20s: %s\n"

	fmt.Printf(format, fmt.Sprintf("%sOS%s", Magenta, Reset), getOSInfo())
	fmt.Printf(format, fmt.Sprintf("%sPackages%s", Yellow, Reset), getPackageInfo())
	fmt.Printf(format, fmt.Sprintf("%sKernel%s", Yellow, Reset), getKernelInfo())
	fmt.Printf(format, fmt.Sprintf("%sUptime%s", Yellow, Reset), getUptimeInfo())
	fmt.Printf(format, fmt.Sprintf("%sShell%s", Yellow, Reset), getShellInfo())
	fmt.Printf(format, fmt.Sprintf("%sCPU%s", Blue, Reset), getCPUInfo())
	fmt.Printf(format, fmt.Sprintf("%sGPU%s", Yellow, Reset), getGPUInfo())
	fmt.Printf(format, fmt.Sprintf("%sStorage%s", Yellow, Reset), getStorageInfo())
	fmt.Printf(format, fmt.Sprintf("%sRAM%s", Yellow, Reset), getRAMInfo())
	fmt.Printf(format, fmt.Sprintf("%sWM%s", Cyan, Reset), getWMInfo())
	fmt.Printf(format, fmt.Sprintf("%sDesktop Environment%s", Red, Reset), getDesktopEnvironment())
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
		return "Unknown"
	}

	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	// Count the number of lines in the output to get the package count
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	packageCount := len(lines)

	return fmt.Sprintf("%d", packageCount)
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
	return getUptime()
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

	return shell
}

func getCPUInfo() string {
	cmd := exec.Command("sh", "-c", "lscpu | grep 'Model name' | awk -F':' '{print $2}'")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return strings.TrimSpace(string(out))
}

func getGPUInfo() string {
	cmd := exec.Command("sh", "-c", "lspci | grep 'VGA' | awk -F':' '{print $3}'")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return strings.TrimSpace(string(out))
}

func getStorageInfo() string {
	cmd := exec.Command("df", "-h", "--output=size,used,avail,pcent,target", "/")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return strings.TrimSpace(string(out))
}

func getRAMInfo() string {
	cmd := exec.Command("sh", "-c", "free -h | grep 'Mem' | awk '{print $2, $3, $4, $7}'")
	out, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	return strings.TrimSpace(string(out))
}

func getWMInfo() string {
	// Try to get window manager information using XDG_SESSION_TYPE
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	if sessionType != "" {
		return sessionType
	}

	// Try to get window manager information using xprop
	cmd := exec.Command("xprop", "-root", "_NET_WM_NAME")
	out, err := cmd.Output()
	if err == nil {
		// Parse the output to extract the window manager name
		if parts := strings.SplitN(string(out), "=", 2); len(parts) == 2 {
			return strings.TrimSpace(parts[1])
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
