# Fetch

Fetch is a simple command-line utility written in Go (Golang) that provides detailed information about your system.
It gathers information about the operating system, packages, kernel, uptime, shell, CPU, GPU, storage, RAM, and window manager.

![Fetch](assets/fetch1.png)

![Fetch](assets/fetch2.png)

![Fetch](assets/fetch3.png)

## Usage

### Prerequisites

Make sure that you have Go installed on your system.
If not, you can download and install it from the [official Go website](https://golang.org/dl/).

### Void Linux

```bash
sudo xbps-install -S go
```

### Arch Linux

```bash
sudo pacman -S go
```

### Installation

Clone the Fetch repository and build the executable:

```bash
git clone https://github.com/elbachir-one/fetch
cd fetch
go build fetch.go
```

### Run Fetch

After building the executable, run Fetch to display system information:

```bash
./fetch
```

As for installation, just put it in your PATH.
