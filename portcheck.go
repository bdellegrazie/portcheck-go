package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime/debug"
	"strings"

	"github.com/f1bonacc1/netstat/netstat"
)

var version = "v0.0.1"

type PortcheckArgs struct {
	Addr     string
	Port     uint
	Protocol string
	Silent   bool
}

func printVersion(verbose bool) {
	fmt.Fprintln(os.Stderr, "version", version)
	if verbose {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			log.Fatal("Cannot get build information from binary")
		}
		fmt.Fprintln(os.Stderr, buildInfo.String())
	}
}

func doTcpPortCheck(args PortcheckArgs, hostIP net.IP) (bool, error) {
	// get only listening TCP sockets that match our IP and port
	tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return (s.LocalAddr.IP.Equal(hostIP) &&
			s.LocalAddr.Port == uint16(args.Port) &&
			s.State == netstat.Listen)
	})

	if err != nil {
		return false, err
	} else {
		return len(tabs) > 0, nil
	}
}

func doTcp6PortCheck(args PortcheckArgs, hostIP net.IP) (bool, error) {
	// get only listening TCP sockets that match our IP and port
	tabs, err := netstat.TCP6Socks(func(s *netstat.SockTabEntry) bool {
		return (s.LocalAddr.IP.Equal(hostIP) &&
			s.LocalAddr.Port == uint16(args.Port) &&
			s.State == netstat.Listen)
	})

	if err != nil {
		return false, err
	} else {
		return len(tabs) > 0, nil
	}
}

func doUdpPortCheck(args PortcheckArgs, hostIP net.IP) (bool, error) {
	// For UDP, only "Close" and "Established" states are supported
	tabs, err := netstat.UDPSocks(func(s *netstat.SockTabEntry) bool {
		return (s.LocalAddr.IP.Equal(hostIP) &&
			s.LocalAddr.Port == uint16(args.Port) &&
			s.State == netstat.Close)
	})

	if err != nil {
		return false, err
	} else {
		return len(tabs) > 0, nil
	}
}

func doUdp6PortCheck(args PortcheckArgs, hostIP net.IP) (bool, error) {
	// For UDP, only "Close" and "Established" states are supported
	tabs, err := netstat.UDP6Socks(func(s *netstat.SockTabEntry) bool {
		return (s.LocalAddr.IP.Equal(hostIP) &&
			s.LocalAddr.Port == uint16(args.Port) &&
			s.State == netstat.Close)
	})

	if err != nil {
		return false, err
	} else {
		return len(tabs) > 0, nil
	}
}

func main() {
	args := PortcheckArgs{}
	versionFlagPtr := flag.Bool("version", false, "Get application version")
	flag.StringVar(&args.Addr, "addr", "0.0.0.0", "IP address to check against")
	flag.StringVar(&args.Protocol, "protocol", "tcp", "Protocol, one of: tcp, tcp6, udp, udp6")
	flag.UintVar(&args.Port, "port", 8000, "Port to check (0 > n < 65536)")
	flag.BoolVar(&args.Silent, "silent", false, "Be silent")

	flag.Parse()

	if *versionFlagPtr {
		printVersion(true)
		os.Exit(0)
	}

	hostIP := net.ParseIP(args.Addr)
	if hostIP == nil {
		log.Fatal("Could not parse Addr: '", args.Addr, "'")
	}

	if args.Port == 0 || args.Port >= 65536 {
		log.Fatal("Port must be greater than 0 and less than 65536")
	}

	var result bool
	var err error

	switch strings.ToLower(args.Protocol) {
	case "tcp":
		if !args.Silent {
			log.Printf("tcp://%s:%d", hostIP.String(), args.Port)
		}
		result, err = doTcpPortCheck(args, hostIP)
	case "tcp6":
		if !args.Silent {
			log.Printf("tcp6://%s:%d", hostIP.String(), args.Port)
		}
		result, err = doTcp6PortCheck(args, hostIP)
	case "udp":
		if !args.Silent {
			log.Printf("udp://%s:%d", hostIP.String(), args.Port)
		}
		result, err = doUdpPortCheck(args, hostIP)
	case "udp6":
		if !args.Silent {
			log.Printf("udp6://%s:%d", hostIP.String(), args.Port)
		}
		result, err = doUdp6PortCheck(args, hostIP)
	default:
		log.Fatal("unknown Protocol '", args.Protocol, "'")
	}

	if err != nil {
		log.Fatal("Error during socket check ", err)
	}

	if !result {
		os.Exit(10)
	}
}
