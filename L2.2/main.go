package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

// Список NTP серверов для попытки подключения
var ntpServers = []string{
	"pool.ntp.org",
	"time.google.com",
	"time.windows.com",
	"time.apple.com",
}

func main() {
	// Получаем локальное время
	localTime := time.Now()
	fmt.Printf("Локальное время: %v\n", localTime.Format(time.RFC3339))

	// Пробуем получить время с разных серверов
	var ntpTime time.Time
	var lastErr error
	
	for _, server := range ntpServers {
		fmt.Printf("Попытка подключения к %s...\n", server)
		
		// Устанавливаем увеличенный таймаут в 5 секунд
		ntpTime, lastErr = ntp.Time(server)
		if lastErr == nil {
			fmt.Printf("Успешное подключение к %s\n", server)
			break
		}
		fmt.Printf("Ошибка при подключении к %s: %v\n", server, lastErr)
	}

	if lastErr != nil {
		fmt.Fprintf(os.Stderr, "Не удалось получить NTP время ни с одного сервера. Последняя ошибка: %v\n", lastErr)
		os.Exit(1)
	}

	fmt.Printf("Точное время (NTP): %v\n", ntpTime.Format(time.RFC3339))
	fmt.Printf("Разница между NTP и локальным временем: %v\n", ntpTime.Sub(localTime))
}