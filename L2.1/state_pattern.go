// Package main демонстрирует реализацию паттерна Состояние
package main

import (
	"fmt"
	"time"
)

// OrderState определяет интерфейс для различных состояний заказа
type OrderState interface {
	Name() string
	Next(order *Order) error
	Cancel(order *Order) error
	Deliver(order *Order) error
}

// Order представляет контекст, содержащий текущее состояние
type Order struct {
	state     OrderState
	isPaid    bool
	createdAt time.Time
}

// NewOrder создает новый заказ
func NewOrder() *Order {
	return &Order{
		state:     &NewOrderState{},
		createdAt: time.Now(),
	}
}

// SetState изменяет текущее состояние заказа
func (o *Order) SetState(state OrderState) {
	o.state = state
}

// GetStateName возвращает имя текущего состояния
func (o *Order) GetStateName() string {
	return o.state.Name()
}

// Next переводит заказ в следующее состояние
func (o *Order) Next() error {
	return o.state.Next(o)
}

// Cancel отменяет заказ
func (o *Order) Cancel() error {
	return o.state.Cancel(o)
}

// Deliver доставляет заказ
func (o *Order) Deliver() error {
	return o.state.Deliver(o)
}

// NewOrderState представляет начальное состояние заказа
type NewOrderState struct{}

func (s *NewOrderState) Name() string {
	return "Новый"
}

func (s *NewOrderState) Next(order *Order) error {
	if !order.isPaid {
		return fmt.Errorf("заказ должен быть оплачен перед переходом в следующее состояние")
	}
	order.SetState(&ProcessingState{})
	return nil
}

func (s *NewOrderState) Cancel(order *Order) error {
	order.SetState(&CancelledState{})
	return nil
}

func (s *NewOrderState) Deliver(order *Order) error {
	return fmt.Errorf("невозможно доставить новый заказ")
}

// ProcessingState представляет состояние обработки заказа
type ProcessingState struct{}

func (s *ProcessingState) Name() string {
	return "В обработке"
}

func (s *ProcessingState) Next(order *Order) error {
	order.SetState(&ShippedState{})
	return nil
}

func (s *ProcessingState) Cancel(order *Order) error {
	order.SetState(&CancelledState{})
	return nil
}

func (s *ProcessingState) Deliver(order *Order) error {
	return fmt.Errorf("заказ ещё обрабатывается")
}

// ShippedState представляет состояние отправленного заказа
type ShippedState struct{}

func (s *ShippedState) Name() string {
	return "Отправлен"
}

func (s *ShippedState) Next(order *Order) error {
	order.SetState(&DeliveredState{})
	return nil
}

func (s *ShippedState) Cancel(order *Order) error {
	return fmt.Errorf("невозможно отменить отправленный заказ")
}

func (s *ShippedState) Deliver(order *Order) error {
	order.SetState(&DeliveredState{})
	return nil
}

// DeliveredState представляет состояние доставленного заказа
type DeliveredState struct{}

func (s *DeliveredState) Name() string {
	return "Доставлен"
}

func (s *DeliveredState) Next(order *Order) error {
	return fmt.Errorf("заказ уже доставлен")
}

func (s *DeliveredState) Cancel(order *Order) error {
	return fmt.Errorf("невозможно отменить доставленный заказ")
}

func (s *DeliveredState) Deliver(order *Order) error {
	return fmt.Errorf("заказ уже доставлен")
}

// CancelledState представляет состояние отмененного заказа
type CancelledState struct{}

func (s *CancelledState) Name() string {
	return "Отменен"
}

func (s *CancelledState) Next(order *Order) error {
	return fmt.Errorf("невозможно обработать отмененный заказ")
}

func (s *CancelledState) Cancel(order *Order) error {
	return fmt.Errorf("заказ уже отменен")
}

func (s *CancelledState) Deliver(order *Order) error {
	return fmt.Errorf("невозможно доставить отмененный заказ")
}

// OrderProcessor обрабатывает заказы и демонстрирует их состояния
type OrderProcessor struct {
	order *Order
}

// NewOrderProcessor создает новый процессор заказов
func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{
		order: NewOrder(),
	}
}

// ProcessOrder демонстрирует жизненный цикл заказа
func (p *OrderProcessor) ProcessOrder() {
	// Выводим начальное состояние
	fmt.Printf("\nНачальное состояние заказа: %s\n", p.order.GetStateName())

	// Пробуем доставить неоплаченный заказ
	fmt.Println("\nПопытка доставки неоплаченного заказа:")
	if err := p.order.Deliver(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	// Оплачиваем заказ
	fmt.Println("\nОплата заказа...")
	p.order.isPaid = true

	// Переводим заказ через различные состояния
	fmt.Println("\nПереход к обработке заказа:")
	if err := p.order.Next(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	fmt.Printf("Текущее состояние: %s\n", p.order.GetStateName())

	fmt.Println("\nОтправка заказа:")
	if err := p.order.Next(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	fmt.Printf("Текущее состояние: %s\n", p.order.GetStateName())

	fmt.Println("\nДоставка заказа:")
	if err := p.order.Deliver(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
	fmt.Printf("Текущее состояние: %s\n", p.order.GetStateName())

	// Пробуем выполнить недопустимые операции
	fmt.Println("\nПопытка отмены доставленного заказа:")
	if err := p.order.Cancel(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}
}

func main() {
	processor := NewOrderProcessor()
	processor.ProcessOrder()
}

/*
# Паттерн Состояние (State Pattern) в Go

## Описание паттерна

Состояние — это поведенческий паттерн проектирования, который позволяет объектам менять свое поведение в зависимости от внутреннего состояния. Создается впечатление, что объект меняет свой класс.

## Структура паттерна

1. **Context** (Контекст) - объект, содержащий текущее состояние
2. **State** (Состояние) - интерфейс, определяющий поведение состояния
3. **ConcreteState** (Конкретное состояние) - реализации различных состояний
4. **Client** (Клиент) - использует контекст

## Когда использовать паттерн Состояние

1. Когда поведение объекта зависит от его состояния и должно изменяться во время выполнения
2. Когда в коде встречается много условных операторов, зависящих от состояния объекта
3. Когда переходы между состояниями должны быть явными и контролируемыми
4. Когда состояния можно организовать в иерархию

## Преимущества

1. **Локализация поведения**
   - Поведение для каждого состояния изолировано
   - Упрощение поддержки кода
   - Улучшенная организация

2. **Управление переходами**
   - Явные и контролируемые переходы
   - Защита от недопустимых переходов
   - Централизованное управление

3. **Расширяемость**
   - Легкое добавление новых состояний
   - Независимость от существующего кода
   - Поддержка Open/Closed Principle

4. **Устранение условных операторов**
   - Замена условной логики полиморфизмом
   - Улучшение читаемости
   - Уменьшение сложности

## Недостатки

1. **Увеличение числа классов**
   - Отдельный класс для каждого состояния
   - Усложнение структуры проекта
   - Повышение накладных расходов

2. **Возможная избыточность**
   - Дублирование кода между состояниями
   - Сложность при малом количестве состояний
   - Дополнительные уровни абстракции

3. **Сложность начальной реализации**
   - Необходимость продумать все состояния
   - Определение всех переходов
   - Сложность рефакторинга

## Реальные примеры использования

### 1. Управление документами
```go
type DocumentState interface {
    Edit(doc *Document) error
    Review(doc *Document) error
    Publish(doc *Document) error
}

type DraftState struct{}
type ReviewState struct{}
type PublishedState struct{}
```

### 2. Управление заказами
```go
type OrderState interface {
    Process(order *Order) error
    Ship(order *Order) error
    Deliver(order *Order) error
    Cancel(order *Order) error
}

type NewOrderState struct{}
type ProcessingState struct{}
type ShippedState struct{}
```

### 3. Управление задачами
```go
type TaskState interface {
    Start(task *Task) error
    Pause(task *Task) error
    Complete(task *Task) error
    Fail(task *Task) error
}

type CreatedState struct{}
type RunningState struct{}
type PausedState struct{}
```

### 4. Аутентификация пользователей
```go
type UserState interface {
    Login(user *User) error
    Logout(user *User) error
    Block(user *User) error
}

type AnonymousState struct{}
type AuthenticatedState struct{}
type BlockedState struct{}
```

### 5. Игровые состояния
```go
type GameState interface {
    Start(game *Game) error
    Pause(game *Game) error
    Resume(game *Game) error
    End(game *Game) error
}

type MenuState struct{}
type PlayingState struct{}
type PausedState struct{}
```

## Лучшие практики использования

1. **Проектирование состояний**
   - Определите четкие границы состояний
   - Продумайте все переходы
   - Обеспечьте валидацию переходов

2. **Управление переходами**
   - Централизуйте логику переходов
   - Документируйте возможные переходы
   - Обрабатывайте ошибки переходов

3. **Обработка ошибок**
   - Определите стратегию обработки ошибок
   - Валидируйте переходы
   - Логируйте изменения состояний

4. **Оптимизация**
   - Используйте пул объектов состояний
   - Кэшируйте состояния
   - Минимизируйте создание объектов

## Антипаттерны при использовании State

1. **Нарушение Single Responsibility**
   - Слишком много ответственности в состояниях
   - Смешивание логики переходов и поведения
   - Сложные зависимости между состояниями

2. **Избыточное использование**
   - Создание состояний для простых случаев
   - Усложнение простых систем
   - Излишняя абстракция

3. **Неправильное управление состоянием**
   - Утечки состояния
   - Несогласованные переходы
   - Отсутствие валидации

## Связанные паттерны

1. **Strategy**
   - Похож на State, но без управления переходами
   - Фокус на алгоритмах, а не на состояниях
   - Более простая структура

2. **Command**
   - Может использоваться вместе с State
   - Инкапсуляция операций
   - Поддержка отмены операций

3. **Singleton**
   - Часто используется для состояний
   - Экономия памяти
   - Глобальный доступ к состояниям

## Заключение

Паттерн Состояние особенно полезен для:
- Управления сложными переходами между состояниями
- Инкапсуляции поведения, зависящего от состояния
- Упрощения условной логики
- Обеспечения типобезопасности при работе с состояниями

Для эффективного использования паттерна:
- Тщательно проектируйте состояния и переходы
- Следите за согласованностью состояний
- Обеспечивайте корректную обработку ошибок
- Документируйте возможные переходы между состояниями
- Оптимизируйте использование памяти

Паттерн Состояние - это мощный инструмент для управления поведением объектов в зависимости от их внутреннего состояния, 
который особенно полезен в сложных системах с множеством правил перехода между состояниями. 
При правильном применении паттерн значительно упрощает поддержку и расширение кода, делая его более понятным и надёжным.
*/