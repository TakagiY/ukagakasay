package main

import (
    "net"
    "fmt"
    "os"
    "github.com/gazitt/flago"
    "bufio"
    "strings"
)

func main() {
    port := flago.Int("port", 'P', 9821, "SSTP Port (Default 9821)", nil)
    address := flago.String("address", 'A', "localhost", "SSTP Address (Default localhost)", nil)
    unyuu := flago.Bool("unyuu", 'u', false, "Speaking by partner", nil)
    flago.Usage = func() {
        fmt.Println("Usage: ukagakasay [options..] [message]\n")
        flago.PrintDefaults()
        fmt.Println("If message is not specified, pipe stdin.\n    (In terminal, send eof by `Ctrl D`.)")
    }
    flago.Parse()

    conn, derr := net.Dial("tcp", fmt.Sprint(*address, ":", *port))
    if derr != nil {
        fmt.Println("Dial error: ", derr)
        return
    }
    defer conn.Close()

    header := "SEND STTP/1.1\r\n" +
    "Sender: ukagakasay\r\n" +
    "Charset: UTF-8\r\n" +
    "Script: "
    var message string

    if flago.NArg() != 0 {
        message = strings.Join(flago.Args(), " ")
    } else {
        stdin := bufio.NewScanner(os.Stdin)
        message = ""
        for stdin.Scan() {
            message += stdin.Text()
            message += "\n"
        }
    }
    message = strings.NewReplacer(
        `\`, `\\`,
        "%", `\%`,
        "\r\n", `\n`,
        "\r", `\n`,
        "\n", `\n`,
    ).Replace(message)
    if *unyuu {
        message = `\u` + message
    } else {
        message = `\h` + message
    }
    message += `\e` + "\r\n\r\n"

    _, werr := conn.Write([]byte(header + message))
    if werr != nil {
        fmt.Println("Write error: ", werr)
        return
    }
}
