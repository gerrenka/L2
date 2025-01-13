package main

import "fmt"

// Product - конечный продукт
type Computer struct {
    CPU         string
    RAM         int
    Storage     int
    GraphicsCard string
}

// Builder - интерфейс строителя
type ComputerBuilder interface {
    SetCPU(cpu string) ComputerBuilder
    SetRAM(ram int) ComputerBuilder
    SetStorage(storage int) ComputerBuilder
    SetGraphicsCard(card string) ComputerBuilder
    Build() *Computer
}

// Конкретный строитель
type DesktopComputerBuilder struct {
    computer *Computer
}

func NewDesktopComputerBuilder() ComputerBuilder {
    return &DesktopComputerBuilder{
        computer: &Computer{},
    }
}

func (b *DesktopComputerBuilder) SetCPU(cpu string) ComputerBuilder {
    b.computer.CPU = cpu
    return b
}

func (b *DesktopComputerBuilder) SetRAM(ram int) ComputerBuilder {
    b.computer.RAM = ram
    return b
}

func (b *DesktopComputerBuilder) SetStorage(storage int) ComputerBuilder {
    b.computer.Storage = storage
    return b
}

func (b *DesktopComputerBuilder) SetGraphicsCard(card string) ComputerBuilder {
    b.computer.GraphicsCard = card
    return b
}

func (b *DesktopComputerBuilder) Build() *Computer {
    return b.computer
}

// Director - управляет процессом строительства
type ComputerAssembler struct {
    builder ComputerBuilder
}

func NewComputerAssembler(b ComputerBuilder) *ComputerAssembler {
    return &ComputerAssembler{builder: b}
}

// Пример предустановленной конфигурации
func (d *ComputerAssembler) ConstructGamingPC() *Computer {
    return d.builder.
        SetCPU("Intel Core i9").
        SetRAM(32).
        SetStorage(2000).
        SetGraphicsCard("NVIDIA RTX 4080").
        Build()
}

func main() {
    // Использование паттерна
    builder := NewDesktopComputerBuilder()
    assembler := NewComputerAssembler(builder)
    
    // Создание игрового компьютера
    gamingPC := assembler.ConstructGamingPC()
    fmt.Printf("Gaming PC: %+v\n", gamingPC)
    
    // Создание пользовательской конфигурации
    customPC := builder.
        SetCPU("AMD Ryzen 7").
        SetRAM(16).
        SetStorage(1000).
        SetGraphicsCard("NVIDIA RTX 3060").
        Build()
    fmt.Printf("Custom PC: %+v\n", customPC)
}


/*
# Паттерн Строитель (Builder Pattern) в Go

## Описание паттерна

Строитель - это порождающий паттерн проектирования, который позволяет создавать сложные объекты пошагово. Строитель даёт возможность использовать один и тот же код строительства для получения разных представлений объектов.

## Структура паттерна

1. **Builder** (Строитель) - интерфейс, определяющий методы для создания частей продукта
2. **ConcreteBuilder** (Конкретный строитель) - реализация интерфейса строителя
3. **Director** (Директор) - класс, определяющий порядок шагов строительства
4. **Product** (Продукт) - создаваемый объект

## Когда использовать паттерн Строитель

1. Когда процесс создания объекта должен быть независимым от составных частей объекта и способа их сборки
2. Когда необходимо обеспечить различные представления создаваемого объекта
3. Когда нужно контролировать процесс создания сложного объекта пошагово

## Преимущества

1. **Пошаговое конструирование**
   - Контроль над процессом создания объекта
   - Возможность создавать объекты с разными конфигурациями
   - Изоляция сложной логики конструирования

2. **Повторное использование кода**
   - Один и тот же код строительства для разных представлений
   - Возможность переиспользования компонентов
   - Уменьшение дублирования кода

3. **Улучшенная читаемость**
   - Чистый и понятный код благодаря цепочке методов
   - Явное разделение процесса конструирования
   - Улучшенная поддерживаемость кода

4. **Инкапсуляция**
   - Скрытие сложности создания объекта
   - Изоляция кода конструирования от бизнес-логики
   - Улучшенная модульность

## Недостатки

1. **Усложнение кодовой базы**
   - Требуется создание дополнительных классов
   - Увеличение количества кода
   - Повышение сложности архитектуры

2. **Привязка к конкретному типу**
   - Каждый строитель привязан к определенному типу продукта
   - Сложность при необходимости создания разных типов продуктов
   - Потенциальное дублирование кода при схожих продуктах

3. **Сложность расширения**
   - При добавлении новых свойств требуется изменение интерфейса
   - Возможное нарушение принципа открытости/закрытости
   - Усложнение поддержки при росте количества свойств

## Реальные примеры использования

### 1. Построение SQL-запросов
```go
query := NewQueryBuilder().
    Select("name, age").
    From("users").
    Where("age > 18").
    OrderBy("name DESC").
    Limit(10).
    Build()
```

### 2. Конфигурация HTTP-клиента
```go
client := NewHTTPClientBuilder().
    WithTimeout(time.Second * 30).
    WithRetries(3).
    WithProxy("http://proxy.example.com").
    WithTLS(true).
    Build()
```

### 3. Создание конфигурации сервера
```go
config := NewServerConfigBuilder().
    SetPort(8080).
    SetHost("localhost").
    SetMaxConnections(1000).
    SetSSLCertificate("cert.pem").
    SetSSLKey("key.pem").
    Build()
```

### 4. Формирование PDF-документов
```go
doc := NewPDFBuilder().
    AddHeader("Ежемесячный отчет").
    AddTitle("Продажи за январь 2024").
    AddTable(salesData).
    AddChart(monthlyStats).
    AddFooter("Страница 1 из 10").
    Build()
```

### 5. Конфигурация баз данных
```go
db := NewDatabaseBuilder().
    SetDriver("postgres").
    SetHost("localhost").
    SetPort(5432).
    SetDatabase("myapp").
    SetUser("admin").
    SetPassword("secret").
    SetMaxConnections(100).
    SetSSLMode("require").
    Build()
```

## Лучшие практики использования

1. **Именование методов**
   - Используйте префиксы Set*, With*, Add* для методов строителя
   - Делайте имена методов описательными
   - Соблюдайте единый стиль именования

2. **Валидация**
   - Проверяйте корректность входных данных в методах строителя
   - Реализуйте валидацию в методе Build()
   - Возвращайте ошибки при некорректных данных

3. **Цепочка методов**
   - Возвращайте указатель на строителя из каждого метода
   - Обеспечьте возможность вызова методов в любом порядке
   - Поддерживайте текучий интерфейс (fluent interface)

4. **Документация**
   - Документируйте каждый метод строителя
   - Указывайте обязательные и опциональные параметры
   - Предоставляйте примеры использования

## Антипаттерны при использовании Builder

1. **Избыточное использование**
   - Не используйте паттерн для простых объектов
   - Избегайте создания строителей для объектов с 2-3 полями
   - Не усложняйте код без необходимости

2. **Нарушение Single Responsibility**
   - Не добавляйте бизнес-логику в строитель
   - Разделяйте логику создания и использования объекта
   - Избегайте побочных эффектов в методах строителя

3. **Неправильная валидация**
   - Не откладывайте все проверки до метода Build()
   - Валидируйте данные при их установке
   - Обеспечивайте согласованность объекта

## Заключение

Паттерн Строитель является мощным инструментом для создания сложных объектов в Go. 
Он особенно полезен при работе с конфигурациями, построении запросов и создании документов. 
При правильном использовании паттерн помогает создавать чистый, поддерживаемый и расширяемый код.

*/