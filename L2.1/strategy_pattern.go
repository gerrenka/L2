// Package main демонстрирует реализацию паттерна Цепочка вызовов
package main

import (
	"fmt"
	"strings"
)

// LogLevel определяет уровни логирования
type LogLevel int

const (
	INFO LogLevel = iota
	WARNING
	ERROR
)

// LogEntry представляет запись в логе
type LogEntry struct {
	Level   LogLevel
	Message string
}

// Logger определяет интерфейс для обработчиков логов
type Logger interface {
	SetNext(logger Logger) Logger
	Log(entry LogEntry)
	LogMessage(level LogLevel, message string)
}

// BaseLogger предоставляет базовую функциональность для логгеров
type BaseLogger struct {
	nextLogger Logger
	level      LogLevel
}

// SetNext устанавливает следующий логгер в цепочке
func (l *BaseLogger) SetNext(next Logger) Logger {
	l.nextLogger = next
	return next
}

// Log обрабатывает запись лога
func (l *BaseLogger) Log(entry LogEntry) {
	if l.level <= entry.Level {
		l.logMessage(entry)
	}
	if l.nextLogger != nil {
		l.nextLogger.Log(entry)
	}
}

// LogMessage создает и обрабатывает новую запись лога
func (l *BaseLogger) LogMessage(level LogLevel, message string) {
	l.Log(LogEntry{
		Level:   level,
		Message: message,
	})
}

// logMessage абстрактный метод для логирования
func (l *BaseLogger) logMessage(entry LogEntry) {
	// Должен быть переопределен в конкретных логгерах
}

// ConsoleLogger выводит логи в консоль
type ConsoleLogger struct {
	BaseLogger
}

// NewConsoleLogger создает новый консольный логгер
func NewConsoleLogger(level LogLevel) *ConsoleLogger {
	return &ConsoleLogger{
		BaseLogger: BaseLogger{level: level},
	}
}

func (l *ConsoleLogger) logMessage(entry LogEntry) {
	levelStr := "INFO"
	switch entry.Level {
	case WARNING:
		levelStr = "WARNING"
	case ERROR:
		levelStr = "ERROR"
	}
	fmt.Printf("[Console] %s: %s\n", levelStr, entry.Message)
}

// FileLogger имитирует запись логов в файл
type FileLogger struct {
	BaseLogger
	filename string
}

// NewFileLogger создает новый файловый логгер
func NewFileLogger(level LogLevel, filename string) *FileLogger {
	return &FileLogger{
		BaseLogger: BaseLogger{level: level},
		filename:   filename,
	}
}

func (l *FileLogger) logMessage(entry LogEntry) {
	levelStr := "INFO"
	switch entry.Level {
	case WARNING:
		levelStr = "WARNING"
	case ERROR:
		levelStr = "ERROR"
	}
	fmt.Printf("[File: %s] %s: %s\n", l.filename, levelStr, entry.Message)
}

// AlertLogger отправляет уведомления для критических ошибок
type AlertLogger struct {
	BaseLogger
	adminEmail string
}

// NewAlertLogger создает новый логгер оповещений
func NewAlertLogger(level LogLevel, email string) *AlertLogger {
	return &AlertLogger{
		BaseLogger:  BaseLogger{level: level},
		adminEmail:  email,
	}
}

func (l *AlertLogger) logMessage(entry LogEntry) {
	if entry.Level >= ERROR {
		fmt.Printf("[Alert] Отправка уведомления на %s: %s\n", 
			l.adminEmail, entry.Message)
	}
}

// FilterLogger фильтрует логи по ключевым словам
type FilterLogger struct {
	BaseLogger
	keywords []string
}

// NewFilterLogger создает новый фильтрующий логгер
func NewFilterLogger(level LogLevel, keywords []string) *FilterLogger {
	return &FilterLogger{
		BaseLogger: BaseLogger{level: level},
		keywords:   keywords,
	}
}

func (l *FilterLogger) Log(entry LogEntry) {
	// Проверяем наличие ключевых слов
	for _, keyword := range l.keywords {
		if strings.Contains(strings.ToLower(entry.Message), 
			strings.ToLower(keyword)) {
			fmt.Printf("[Filter] Обнаружено ключевое слово '%s' в сообщении\n", 
				keyword)
			break
		}
	}
	
	// Продолжаем цепочку
	if l.nextLogger != nil {
		l.nextLogger.Log(entry)
	}
}

func main() {
	// Создаем логгеры
	consoleLogger := NewConsoleLogger(INFO)
	fileLogger := NewFileLogger(WARNING, "app.log")
	alertLogger := NewAlertLogger(ERROR, "admin@example.com")
	filterLogger := NewFilterLogger(INFO, []string{"password", "key"})

	// Строим цепочку обработчиков
	filterLogger.SetNext(consoleLogger).
		SetNext(fileLogger).
		SetNext(alertLogger)

	// Тестируем различные уровни логирования
	fmt.Println("=== Тест информационного сообщения ===")
	filterLogger.LogMessage(INFO, "Приложение запущено")

	fmt.Println("\n=== Тест предупреждения ===")
	filterLogger.LogMessage(WARNING, "Низкий заряд батареи")

	fmt.Println("\n=== Тест ошибки ===")
	filterLogger.LogMessage(ERROR, "Критическая ошибка: база данных недоступна")

	fmt.Println("\n=== Тест фильтрации ===")
	filterLogger.LogMessage(INFO, "User entered incorrect password")
}

/*
# Паттерн Цепочка вызовов (Chain of Responsibility) в Go

## Описание паттерна

Цепочка вызовов — это поведенческий паттерн проектирования, который позволяет передавать запросы последовательно по цепочке обработчиков. Каждый последующий обработчик решает, может ли он обработать запрос, и если нет, передаёт его дальше по цепочке.

## Структура паттерна

1. **Handler** (Обработчик) - интерфейс для обработки запросов
2. **BaseHandler** (Базовый обработчик) - реализует связывание обработчиков
3. **ConcreteHandler** (Конкретный обработчик) - выполняет обработку запросов
4. **Client** (Клиент) - отправляет запросы в цепочку обработчиков

## Когда использовать паттерн Цепочка вызовов

1. Когда программа должна обрабатывать разнообразные запросы несколькими способами
2. Когда порядок обработчиков важен
3. Когда набор обработчиков должен быть определён динамически
4. Когда нужно ослабить связанность между отправителем и получателями

## Преимущества

1. **Уменьшение связанности**
   - Отделение отправителя от получателей
   - Гибкость в выборе обработчиков
   - Динамическое изменение цепочки

2. **Принцип единственной ответственности**
   - Каждый обработчик выполняет одну задачу
   - Улучшенная модульность кода
   - Простота добавления новых обработчиков

3. **Гибкость конфигурации**
   - Динамическое построение цепочки
   - Изменение порядка обработки
   - Условное выполнение обработчиков

4. **Возможность пропуска обработки**
   - Обработчики могут пропускать запросы
   - Частичная обработка запросов
   - Выборочное применение обработчиков

## Недостатки

1. **Отсутствие гарантии обработки**
   - Запрос может не быть обработан
   - Сложность отладки
   - Необходимость мониторинга

2. **Возможные задержки**
   - Длинные цепочки замедляют обработку
   - Накладные расходы на передачу
   - Потенциальное снижение производительности

3. **Сложность конфигурации**
   - Правильный порядок обработчиков
   - Потенциальные циклы
   - Сложность отладки

## Реальные примеры использования

### 1. Middleware в веб-приложениях
```go
type Middleware interface {
    Handle(req *Request, next Handler) error
}

type AuthMiddleware struct{}
type LoggerMiddleware struct{}
type RateLimiterMiddleware struct{}
```

### 2. Валидация данных
```go
type Validator interface {
    Validate(data interface{}) error
    SetNext(validator Validator)
}

type SchemaValidator struct{}
type TypeValidator struct{}
type RangeValidator struct{}
```

### 3. Обработка событий
```go
type EventHandler interface {
    HandleEvent(event Event)
    SetNext(handler EventHandler)
}

type LogEventHandler struct{}
type MetricsEventHandler struct{}
type NotificationHandler struct{}
```

### 4. Фильтрация контента
```go
type ContentFilter interface {
    Filter(content string) string
    SetNext(filter ContentFilter)
}

type ProfanityFilter struct{}
type SpamFilter struct{}
type HTMLFilter struct{}
```

### 5. Система авторизации
```go
type AuthorizationChecker interface {
    Check(user User, resource Resource) bool
    SetNext(checker AuthorizationChecker)
}

type RoleChecker struct{}
type PermissionChecker struct{}
type OwnershipChecker struct{}
```

## Лучшие практики использования

1. **Проектирование обработчиков**
   - Соблюдайте единственную ответственность
   - Обеспечивайте независимость обработчиков
   - Используйте внедрение зависимостей

2. **Управление цепочкой**
   - Правильно определяйте порядок
   - Обрабатывайте граничные случаи
   - Реализуйте мониторинг цепочки

3. **Обработка ошибок**
   - Определите стратегию обработки ошибок
   - Обеспечьте логирование
   - Реализуйте восстановление

4. **Оптимизация производительности**
   - Ограничивайте длину цепочки
   - Используйте быстрые проверки
   - Кэшируйте результаты

## Антипаттерны при использовании Chain of Responsibility

1. **Длинные цепочки**
   - Избегайте слишком длинных цепочек
   - Разбивайте на подцепочки
   - Используйте композицию

2. **Нарушение Single Responsibility**
   - Не смешивайте разные обязанности
   - Разделяйте сложные обработчики
   - Сохраняйте фокус обработчиков

3. **Неправильная обработка ошибок**
   - Не игнорируйте ошибки
   - Обеспечивайте обратную связь
   - Реализуйте восстановление

## Заключение

Паттерн Цепочка вызовов является эффективным инструментом для:
- Построения гибких систем обработки запросов
- Уменьшения связанности компонентов
- Реализации сложной логики обработки
- Динамического конфигурирования обработчиков

Для эффективного использования паттерна:
- Тщательно проектируйте интерфейсы
- Правильно управляйте цепочкой
- Обеспечивайте мониторинг и отладку
- Оптимизируйте производительность
*/