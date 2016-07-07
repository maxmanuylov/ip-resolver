package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "net"
    "os"
    "path/filepath"
    "time"
)

func main() {
    network := flag.String("network", "tcp4", "network to use")
    address := flag.String("address", "", "remote address to connect to")
    timeout := flag.Int("timeout", 5, "timeout in seconds")
    file := flag.String("file", "", "file to save ip address variable to (optional)")
    varName := flag.String("var", "PUBLIC_IPV4", "ip address variable name")

    flag.Parse()

    if *address == "" {
        fmt.Fprintf(os.Stderr, "\"address\" is not specified\n")
        fmt.Fprintf(os.Stderr, "Available options:\n")
        flag.PrintDefaults()
        os.Exit(255)
    }

    conn, err := net.DialTimeout(*network, *address, time.Duration(*timeout) * time.Second)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Connection failed: %s\n", err.Error())
        os.Exit(255)
    }

    defer conn.Close()

    localAddress := conn.LocalAddr().String()
    ip, _, err := net.SplitHostPort(localAddress)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to parse local address (%s): %s\n", localAddress, err.Error())
        os.Exit(255)
    }

    if *file == "" {
        fmt.Fprintf(os.Stdout, "%s", ip)
    } else {
        if err := os.MkdirAll(filepath.Dir(*file), os.FileMode(0777)); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to create parent directory: %s\n", err.Error())
            os.Exit(255)
        }
        if err := ioutil.WriteFile(*file, []byte(fmt.Sprintf("%s=%s", *varName, ip)), os.FileMode(0666)); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to save file: %s\n", err.Error())
            os.Exit(255)
        }
    }
}
