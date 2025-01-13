package main

import (
    "flag"
    "fmt"
    "io"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Разбор аргументов командной строки
    timeout := flag.Duration("timeout", 10*time.Second, "таймаут для подключения")
    flag.Parse()

    args := flag.Args()
    if len(args) != 2 {
        fmt.Fprintf(os.Stderr, "Использование: %s [--timeout=10s] хост порт\n", os.Args[0])
        os.Exit(1)
    }

    host, port := args[0], args[1]
    address := net.JoinHostPort(host, port)

    // Создание подключения с таймаутом
    conn, err := net.DialTimeout("tcp", address, *timeout)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Ошибка подключения: %v\n", err)
        os.Exit(1)
    }
    defer conn.Close()

    // Обработка Ctrl+D и других сигналов
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
    done := make(chan bool, 1)

    // Обработка чтения из соединения
    go func() {
        _, err := io.Copy(os.Stdout, conn)
        if err != nil {
            if err != io.EOF {
                fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
            }
        }
        done <- true
    }()

    // Обработка записи в соединение
    go func() {
        _, err := io.Copy(conn, os.Stdin)
        if err != nil {
            if err != io.EOF {
                fmt.Fprintf(os.Stderr, "Ошибка записи: %v\n", err)
            }
        }
        done <- true
    }()

    // Ожидание сигнала завершения или прерывания
    select {
    case <-signalChan:
        return
    case <-done:
        return
    }
}