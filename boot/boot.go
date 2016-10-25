package ip_resolver

import (
    "flag"
    "fmt"
    "github.com/maxmanuylov/utils/application"
    "io/ioutil"
    "net"
    "os"
    "path/filepath"
    "time"
)

func Run() {
    network := flag.String("network", "tcp4", "network to use")
    address := flag.String("address", "", "remote address to connect to")
    timeout := flag.Int("timeout", 5, "timeout in seconds")
    file := flag.String("file", "", "file to save the variables to (optional)")
    ipv4VarName := flag.String("ipv4-var", "PUBLIC_IPV4", "IP v4 address variable name")
    ipv4NetworkVarName := flag.String("ipv4-network-var", "PUBLIC_IPV4_NETWORK", "IP v4 network variable name")

    flag.Parse()

    if *address == "" {
        fmt.Fprintln(os.Stderr, "\"address\" is not specified")
        fmt.Fprintln(os.Stderr, "Available options:")
        flag.PrintDefaults()
        os.Exit(255)
    }

    ip, err := ResolveLocalIP(*network, *address, time.Duration(*timeout) * time.Second)
    if err != nil {
        application.Exit(err.Error())
    }

    if *file == "" {
        fmt.Fprintf(os.Stdout, "%s", ip)
    } else {
        ipNetwork, err := FindIPNetwork(ip)
        if err != nil {
            application.Exit(err.Error())
        }

        if err := os.MkdirAll(filepath.Dir(*file), os.FileMode(0777)); err != nil {
            application.Exit(fmt.Sprintf("Failed to create parent directory: %v", err))
        }

        fileText := fmt.Sprintf("%s=%s\n%s=%s", *ipv4VarName, ip, *ipv4NetworkVarName, ipNetwork)

        if err := ioutil.WriteFile(*file, []byte(fileText), os.FileMode(0666)); err != nil {
            application.Exit(fmt.Sprintf("Failed to save file: %v", err))
        }
    }
}

func ResolveLocalIP(network, address string, timeout time.Duration) (string, error) {
    conn, err := net.DialTimeout(network, address, timeout)
    if err != nil {
        return "", fmt.Errorf("Connection failed: %v", err)
    }

    defer conn.Close()

    localAddress := conn.LocalAddr().String()
    ip, _, err := net.SplitHostPort(localAddress)
    if err != nil {
        return "", fmt.Errorf("Failed to parse local address (%s): %v", localAddress, err)
    }

    return ip, nil
}

func FindIPNetwork(ipStr string) (string, error) {
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return "", fmt.Errorf("Failed to parse IP address: %s", ipStr)
    }

    ifaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }

    for _, iface := range ifaces {
        addrs, err := iface.Addrs()
        if err != nil {
            return "", err
        }

        for _, addr := range addrs {
            if ipNet, ok := addr.(*net.IPNet); ok && ipNet.Contains(ip) {
                return ipNet.String(), nil
            }
        }
    }

    return "", fmt.Errorf("Failed to find network interface for IP address: %s", ipStr)
}
