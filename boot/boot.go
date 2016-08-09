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
    file := flag.String("file", "", "file to save ip address variable to (optional)")
    varName := flag.String("var", "PUBLIC_IPV4", "ip address variable name")

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
        if err := os.MkdirAll(filepath.Dir(*file), os.FileMode(0777)); err != nil {
            application.Exit(fmt.Sprintf("Failed to create parent directory: %v", err))
        }
        if err := ioutil.WriteFile(*file, []byte(fmt.Sprintf("%s=%s", *varName, ip)), os.FileMode(0666)); err != nil {
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
