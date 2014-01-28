package main
 
import (
    "net"
    "os"
    "time"
    "flag"
)

const (
    MiB = 1024 * 1024
)

var buffer int = 1

func init() {
    flag.IntVar(&buffer, "buffer", 1, "the buffer size for write/read (in MiBs)")
}

func main() {
    flag.Parse()
    args := flag.Args()
    if len(args) == 1 {
        client(args[0])
    } else if len(args) == 0 {
        server()
    } else {
        println("Bad args.")
        os.Exit(1)
    }
}

func client(addr string) {
    buf := make([]byte, buffer * MiB)
    conn, err := net.Dial("tcp", addr + ":54321")
    if err != nil {
        println("error connecting:", err.Error())
        os.Exit(1)
    }
    for {
        _, err := conn.Write(buf)
        if err != nil {
            println("Error reading:", err.Error())
            return
        }
    }
}

func server() {
    println("Starting the server")
 
    listener, err := net.Listen("tcp", "0.0.0.0:54321")
    if err != nil {
        println("error listening:", err.Error())
        os.Exit(1)
    }
 
    for {
        conn, err := listener.Accept()
        println("Connection from:", conn)
        if err != nil {
            println("Error accept:", err.Error())
            return
        }
        go recv_loop(conn)
    }
}
 
func recv_loop(conn net.Conn) {
    buf := make([]byte, buffer * MiB)

    bytes := 0
    start := time.Now()

    for {
        n, err := conn.Read(buf)
        if err != nil {
            println("Error reading:", err.Error())
            return
        }
        bytes += n
        duration := time.Since(start)
        if duration.Seconds() > 1 {
            mbs := int(float64(bytes) / duration.Seconds() / 1024. / 1024.)
            println("Speed:", mbs, " MB/s")
            bytes = 0
            start = time.Now()
        }
    }
}
