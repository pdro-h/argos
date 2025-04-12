package main

import (
    "fmt"
    "net"
    "sort"
    "sync"
    "time"
)

const maxPort = 1024
const timeout = 500 * time.Millisecond // 0.5 seconds

func scanPort(wg *sync.WaitGroup, host string, port int, openPorts *[]int, mu *sync.Mutex) {
    defer wg.Done()
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err == nil && conn != nil {
        mu.Lock()
        *openPorts = append(*openPorts, port)
        mu.Unlock()
        conn.Close()
    }
}

func main() {
    var host string
    fmt.Print("Type the host: ")
    fmt.Scanln(&host)

    var wg sync.WaitGroup
    var mu sync.Mutex
    openPorts := []int{}

    for port := 1; port <= maxPort; port++ {
        wg.Add(1)
        go scanPort(&wg, host, port, &openPorts, &mu)
    }

    wg.Wait()

    sort.Ints(openPorts)
    fmt.Println("Open ports:")
    for _, port := range openPorts {
        fmt.Printf(" - %d open\n", port)
    }
}