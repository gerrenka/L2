package main

import "fmt"

// Подсистема: Система безопасности
type SecuritySystem struct{}

func (s *SecuritySystem) Arm() {
    fmt.Println("Система безопасности активирована")
}

func (s *SecuritySystem) Disarm() {
    fmt.Println("Система безопасности деактивирована")
}

// Подсистема: Система освещения
type LightingSystem struct{}

func (l *LightingSystem) TurnOn() {
    fmt.Println("Освещение включено")
}

func (l *LightingSystem) TurnOff() {
    fmt.Println("Освещение выключено")
}

// Подсистема: Система климат-контроля
type ClimateControl struct{}

func (c *ClimateControl) SetTemperature(temp int) {
    fmt.Printf("Установлена температура: %d градусов\n", temp)
}

// Фасад для управления умным домом
type SmartHomeFacade struct {
    security *SecuritySystem
    lighting *LightingSystem
    climate  *ClimateControl
}

// Конструктор фасада
func NewSmartHomeFacade() *SmartHomeFacade {
    return &SmartHomeFacade{
        security: &SecuritySystem{},
        lighting: &LightingSystem{},
        climate:  &ClimateControl{},
    }
}

// Методы фасада для типовых сценариев
func (f *SmartHomeFacade) LeaveHome() {
    f.security.Arm()
    f.lighting.TurnOff()
    f.climate.SetTemperature(18) // экономный режим
}

func (f *SmartHomeFacade) ReturnHome() {
    f.security.Disarm()
    f.lighting.TurnOn()
    f.climate.SetTemperature(22) // комфортная температура
}

func main() {
    // Использование фасада
    smartHome := NewSmartHomeFacade()
    
    fmt.Println("Уходим из дома:")
    smartHome.LeaveHome()
    
    fmt.Println("\nВозвращаемся домой:")
    smartHome.ReturnHome()
}


/*
# Паттерн Фасад (Facade Pattern) в Go

## Описание паттерна

Фасад - это структурный паттерн проектирования, который предоставляет простой интерфейс к сложной системе классов, библиотеке или фреймворку. Фасад предлагает интерфейс более высокого уровня, который упрощает использование системы.

## Структура паттерна

1. **Facade** (Фасад) - предоставляет унифицированный интерфейс к подсистемам
2. **Subsystems** (Подсистемы) - реализуют функциональность системы
3. **Client** (Клиент) - работает с системой через фасад

## Когда использовать паттерн Фасад

1. Когда нужно представить простой или урезанный интерфейс к сложной системе
2. Когда необходимо разложить систему на подсистемы
3. Когда требуется уменьшить связанность между клиентами и компонентами системы
4. Когда нужно создать точку входа к каждому уровню подсистемы

## Преимущества

1. **Изоляция клиентов от компонентов сложной подсистемы**
   - Уменьшение сложности использования системы
   - Сокращение количества зависимостей
   - Упрощение миграции к новым версиям
   
2. **Уменьшение связанности**
   - Слабая связь между клиентами и подсистемами
   - Возможность изменения подсистем без влияния на клиентов
   - Улучшенная модульность

3. **Упрощение использования системы**
   - Единая точка входа в систему
   - Понятный высокоуровневый интерфейс
   - Сокращение количества кода у клиентов

4. **Улучшение поддерживаемости**
   - Централизованное управление подсистемами
   - Изолированные изменения
   - Упрощенное тестирование

## Недостатки

1. **Риск превращения в божественный объект**
   - Фасад может стать слишком большим
   - Возможное нарушение принципа единственной ответственности
   - Сложность поддержки при росте функциональности

2. **Увеличение накладных расходов**
   - Дополнительный слой абстракции
   - Потенциальное снижение производительности
   - Увеличение количества кода

3. **Ограничение гибкости**
   - Не все возможности подсистем доступны через фасад
   - Возможная избыточная простота интерфейса
   - Потеря специфичной функциональности

## Реальные примеры использования

### 1. Работа с базой данных
```go
type DatabaseFacade struct {
    connection *Connection
    cache     *Cache
    logger    *Logger
}

func (f *DatabaseFacade) GetUser(id int) (*User, error) {
    // Проверка кеша
    // Подключение к БД
    // Логирование операции
    // Возврат результата
}
```

### 2. Платежная система
```go
type PaymentFacade struct {
    validator   *PaymentValidator
    processor   *PaymentProcessor
    notifier    *NotificationService
}

func (f *PaymentFacade) ProcessPayment(payment *Payment) error {
    // Валидация платежа
    // Обработка транзакции
    // Отправка уведомления
    // Возврат результата
}
```

### 3. Система аутентификации
```go
type AuthFacade struct {
    userService  *UserService
    tokenService *TokenService
    permissions  *PermissionService
}

func (f *AuthFacade) Login(username, password string) (*Session, error) {
    // Проверка учетных данных
    // Создание токена
    // Проверка прав доступа
    // Создание сессии
}
```

### 4. Работа с файловой системой
```go
type FileSystemFacade struct {
    reader     *FileReader
    writer     *FileWriter
    encryptor  *Encryptor
}

func (f *FileSystemFacade) SaveEncryptedFile(data []byte, path string) error {
    // Шифрование данных
    // Создание файла
    // Запись данных
    // Проверка целостности
}
```

### 5. API клиент
```go
type APIClientFacade struct {
    httpClient  *HTTPClient
    rateLimit   *RateLimiter
    cache       *Cache
}

func (f *APIClientFacade) FetchData(url string) (*Response, error) {
    // Проверка кеша
    // Проверка лимитов
    // Выполнение запроса
    // Сохранение в кеш
}
```

## Лучшие практики использования

1. **Проектирование интерфейса**
   - Делайте интерфейс простым и понятным
   - Скрывайте сложность реализации
   - Группируйте связанные операции

2. **Управление зависимостями**
   - Внедряйте зависимости через конструктор
   - Используйте интерфейсы для подсистем
   - Применяйте принцип инверсии зависимостей

3. **Обработка ошибок**
   - Предоставляйте понятные сообщения об ошибках
   - Не пропускайте ошибки подсистем
   - Логируйте важные операции

4. **Тестирование**
   - Тестируйте фасад как единое целое
   - Используйте моки для подсистем
   - Покрывайте тестами граничные случаи

## Антипаттерны при использовании Facade

1. **Нарушение единственной ответственности**
   - Не добавляйте бизнес-логику в фасад
   - Избегайте прямой работы с данными
   - Разделяйте большие фасады на меньшие

2. **Избыточная простота**
   - Не упрощайте интерфейс в ущерб функциональности
   - Сохраняйте возможность прямого доступа к подсистемам
   - Обеспечивайте необходимую гибкость

3. **Нарушение инкапсуляции**
   - Не предоставляйте доступ к внутреннему состоянию
   - Скрывайте детали реализации
   - Защищайте подсистемы от неправильного использования

## Заключение

Паттерн Фасад является отличным инструментом для упрощения работы со сложными системами в Go. Он особенно полезен при:
- Интеграции с внешними сервисами
- Работе с legacy кодом
- Создании API для библиотек
- Упрощении сложных подсистем

Для эффективного использования паттерна Фасад:
- Тщательно проектируйте интерфейс
- Следите за размером и ответственностью фасада
- Обеспечивайте правильную обработку ошибок
- Поддерживайте хорошее покрытие тестами
*/
