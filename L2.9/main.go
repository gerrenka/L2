package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "os/exec"
    "strings"
    "strconv"
    "syscall"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("$ ") // shell prompt
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            continue
        }

        input = strings.TrimSpace(input)
        if input == "\\quit" {
            break
        }

        if err := executeCommand(input); err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}

func executeCommand(input string) error {
    // Разбиваем команду на пайпы
    pipedCmds := strings.Split(input, "|")
    var commands [][]string

    // Парсим каждую команду в пайпе
    for _, cmd := range pipedCmds {
        args := strings.Fields(strings.TrimSpace(cmd))
        if len(args) == 0 {
            continue
        }
        commands = append(commands, args)
    }

    if len(commands) == 0 {
        return nil
    }

    // Если команда одна - выполняем её напрямую
    if len(commands) == 1 {
        return executeSingleCommand(commands[0])
    }

    return executePipeline(commands)
}

func executeSingleCommand(args []string) error {
    switch args[0] {
    case "cd":
        if len(args) < 2 {
            return os.Chdir(os.Getenv("HOME"))
        }
        return os.Chdir(args[1])
    case "pwd":
        dir, err := os.Getwd()
        if err != nil {
            return err
        }
        fmt.Println(dir)
        return nil
    case "echo":
        fmt.Println(strings.Join(args[1:], " "))
        return nil
    case "kill":
        if len(args) != 2 {
            return fmt.Errorf("kill: неверное количество аргументов")
        }
        pid, err := strconv.Atoi(args[1])
        if err != nil {
            return fmt.Errorf("kill: неверный PID: %v", err)
        }
        return syscall.Kill(pid, syscall.SIGTERM)
    case "ps":
        cmd := exec.Command("ps", "aux")
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        return cmd.Run()
    default:
        cmd := exec.Command(args[0], args[1:]...)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Stdin = os.Stdin
        return cmd.Run()
    }
}

func executePipeline(commands [][]string) error {
    var pipes [][2]io.ReadWriteCloser
    var processes []*exec.Cmd

    // Создаем пайпы между командами
    for i := 0; i < len(commands)-1; i++ {
        readPipe, writePipe := io.Pipe()
        pipes = append(pipes, [2]io.ReadWriteCloser{readPipe, writePipe})
    }

    // Создаем команды
    for i, command := range commands {
        cmd := exec.Command(command[0], command[1:]...)

        // Настраиваем stdin для первой команды
        if i == 0 {
            cmd.Stdin = os.Stdin
        } else {
            cmd.Stdin = pipes[i-1][0]
        }

        // Настраиваем stdout для последней команды
        if i == len(commands)-1 {
            cmd.Stdout = os.Stdout
        } else {
            cmd.Stdout = pipes[i][1]
        }

        cmd.Stderr = os.Stderr
        processes = append(processes, cmd)
    }

    // Запускаем все процессы
    for _, cmd := range processes {
        if err := cmd.Start(); err != nil {
            return fmt.Errorf("ошибка запуска команды: %v", err)
        }
    }

    // Ждем завершения всех процессов
    for _, cmd := range processes {
        if err := cmd.Wait(); err != nil {
            return fmt.Errorf("ошибка выполнения команды: %v", err)
        }
    }

    // Закрываем все пайпы
    for _, pipe := range pipes {
        pipe[0].Close()
        pipe[1].Close()
    }

    return nil
}