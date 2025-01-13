/*
# Паттерн Посетитель (Visitor Pattern) в Go

## Описание паттерна

Посетитель — это поведенческий паттерн проектирования, который позволяет добавлять новые операции к объектам без изменения самих объектов. Посетитель реализует принцип открытости/закрытости, позволяя добавлять новую функциональность без изменения существующего кода.

## Структура паттерна

1. **Visitor** (Посетитель) - интерфейс, объявляющий методы посещения для каждого типа элемента
2. **ConcreteVisitor** (Конкретный посетитель) - реализует интерфейс посетителя
3. **Element** (Элемент) - интерфейс, объявляющий метод accept для посетителя
4. **ConcreteElement** (Конкретный элемент) - реализует метод accept

## Когда использовать паттерн Посетитель

1. Когда нужно выполнить операцию над всеми элементами сложной структуры объектов
2. Когда классы структуры редко меняются, но часто добавляются новые операции
3. Когда связанная функциональность не должна быть распределена по всем классам
4. Когда структура данных и алгоритмы должны быть разделены

## Преимущества

1. **Разделение алгоритма от структуры**
   - Чистое разделение данных и операций
   - Упрощение добавления новых операций
   - Улучшенная модульность кода

2. **Централизация родственных операций**
   - Группировка связанной функциональности
   - Улучшенная организация кода
   - Упрощение поддержки

3. **Накопление состояния**
   - Возможность накапливать информацию при обходе
   - Простота передачи данных между операциями
   - Удобство сбора статистики

4. **Принцип открытости/закрытости**
   - Легкое добавление новых операций
   - Отсутствие изменений в существующих классах
   - Улучшенная расширяемость

## Недостатки

1. **Нарушение инкапсуляции**
   - Необходимость доступа к внутренним данным
   - Потенциальное раскрытие деталей реализации
   - Усложнение поддержки приватности

2. **Сложность добавления новых классов**
   - Необходимость обновления всех посетителей
   - Усложнение расширения иерархии классов
   - Возможное нарушение принципа открытости/закрытости

3. **Увеличение сложности**
   - Дополнительные классы и интерфейсы
   - Усложнение навигации по коду
   - Повышение порога вхождения

## Реальные примеры использования

### 1. Обработка документов
```go
type DocumentVisitor interface {
    VisitPDF(*PDFDocument) error
    VisitWord(*WordDocument) error
    VisitExcel(*ExcelDocument) error
}

type DocumentConverter struct{}

func (dc *DocumentConverter) VisitPDF(doc *PDFDocument) error {
    // Конвертация PDF в другой формат
}
```

### 2. Валидация форм
```go
type FormValidator interface {
    VisitTextInput(*TextInput) []error
    VisitCheckbox(*Checkbox) []error
    VisitSelect(*Select) []error
}

type ValidationVisitor struct {
    Rules map[string][]string
}
```

### 3. Экспорт данных
```go
type DataExporter interface {
    VisitUser(*User) []byte
    VisitOrder(*Order) []byte
    VisitProduct(*Product) []byte
}

type JSONExporter struct{}
type XMLExporter struct{}
```

### 4. Анализ кода
```go
type ASTVisitor interface {
    VisitFunction(*FunctionNode) error
    VisitVariable(*VariableNode) error
    VisitClass(*ClassNode) error
}

type CodeAnalyzer struct{
    Metrics map[string]int
}
```

### 5. Обработка платежей
```go
type PaymentVisitor interface {
    VisitCreditCard(*CreditCard) error
    VisitPayPal(*PayPal) error
    VisitBankTransfer(*BankTransfer) error
}

type PaymentProcessor struct{
    Fee float64
}
```

## Лучшие практики использования

1. **Проектирование интерфейса**
   - Используйте говорящие имена методов
   - Группируйте связанные операции
   - Продумывайте иерархию посетителей

2. **Управление состоянием**
   - Явно передавайте контекст
   - Избегайте глобального состояния
   - Используйте иммутабельные данные

3. **Обработка ошибок**
   - Определите стратегию обработки ошибок
   - Используйте понятные сообщения
   - Обеспечьте восстановление после ошибок

4. **Тестирование**
   - Тестируйте каждого посетителя отдельно
   - Проверяйте граничные случаи
   - Используйте моки для сложных структур

## Антипаттерны при использовании Visitor

1. **Избыточная функциональность**
   - Не добавляйте лишние методы в интерфейс
   - Разделяйте несвязанные операции
   - Избегайте "мусорных" посетителей

2. **Нарушение Single Responsibility**
   - Не смешивайте разные операции
   - Разделяйте посетителей по ответственности
   - Сохраняйте фокус каждого посетителя

3. **Неправильное управление состоянием**
   - Избегайте скрытых зависимостей
   - Не полагайтесь на порядок посещения
   - Не храните избыточное состояние

## Заключение

Паттерн Посетитель является мощным инструментом для работы со сложными структурами объектов в Go. Он особенно полезен когда:
- Структура объектов стабильна
- Часто добавляются новые операции
- Нужно разделить алгоритмы и структуры данных
- Требуется централизовать схожую функциональность

Для эффективного использования паттерна Посетитель:
- Тщательно проектируйте иерархию классов
- Следите за чистотой интерфейсов
- Разделяйте несвязанные операции
- Обеспечивайте хорошее покрытие тестами
*/

// Package main демонстрирует реализацию паттерна Посетитель (Visitor Pattern)
package main

import (
	"fmt"
)

// Интерфейс Visitor определяет методы для посещения каждого типа элемента
type Visitor interface {
	VisitCircle(*Circle) string
	VisitRectangle(*Rectangle) string
	VisitTriangle(*Triangle) string
}

// Shape определяет интерфейс для геометрических фигур
type Shape interface {
	Accept(Visitor) string
}

// Circle представляет круг
type Circle struct {
	Radius float64
}

func (c *Circle) Accept(v Visitor) string {
	return v.VisitCircle(c)
}

// Rectangle представляет прямоугольник
type Rectangle struct {
	Width  float64
	Height float64
}

func (r *Rectangle) Accept(v Visitor) string {
	return v.VisitRectangle(r)
}

// Triangle представляет треугольник
type Triangle struct {
	Base   float64
	Height float64
}

func (t *Triangle) Accept(v Visitor) string {
	return v.VisitTriangle(t)
}

// AreaCalculator вычисляет площадь фигур
type AreaCalculator struct{}

func (ac *AreaCalculator) VisitCircle(c *Circle) string {
	area := 3.14 * c.Radius * c.Radius
	return fmt.Sprintf("Площадь круга: %.2f", area)
}

func (ac *AreaCalculator) VisitRectangle(r *Rectangle) string {
	area := r.Width * r.Height
	return fmt.Sprintf("Площадь прямоугольника: %.2f", area)
}

func (ac *AreaCalculator) VisitTriangle(t *Triangle) string {
	area := 0.5 * t.Base * t.Height
	return fmt.Sprintf("Площадь треугольника: %.2f", area)
}

// PerimeterCalculator вычисляет периметр фигур
type PerimeterCalculator struct{}

func (pc *PerimeterCalculator) VisitCircle(c *Circle) string {
	perimeter := 2 * 3.14 * c.Radius
	return fmt.Sprintf("Периметр круга: %.2f", perimeter)
}

func (pc *PerimeterCalculator) VisitRectangle(r *Rectangle) string {
	perimeter := 2 * (r.Width + r.Height)
	return fmt.Sprintf("Периметр прямоугольника: %.2f", perimeter)
}

func (pc *PerimeterCalculator) VisitTriangle(t *Triangle) string {
	// Упрощенный расчет для равнобедренного треугольника
	perimeter := t.Base + 2*t.Height
	return fmt.Sprintf("Периметр треугольника: %.2f", perimeter)
}

// DrawingVisitor отрисовывает фигуры (имитация)
type DrawingVisitor struct{}

func (dv *DrawingVisitor) VisitCircle(c *Circle) string {
	return fmt.Sprintf("Рисуем круг с радиусом %.2f", c.Radius)
}

func (dv *DrawingVisitor) VisitRectangle(r *Rectangle) string {
	return fmt.Sprintf("Рисуем прямоугольник %.2fx%.2f", r.Width, r.Height)
}

func (dv *DrawingVisitor) VisitTriangle(t *Triangle) string {
	return fmt.Sprintf("Рисуем треугольник с основанием %.2f и высотой %.2f", t.Base, t.Height)
}

func main() {
	// Создаем фигуры
	shapes := []Shape{
		&Circle{Radius: 5},
		&Rectangle{Width: 4, Height: 6},
		&Triangle{Base: 3, Height: 4},
	}

	// Создаем посетителей
	areaCalc := &AreaCalculator{}
	perimeterCalc := &PerimeterCalculator{}
	drawer := &DrawingVisitor{}

	// Применяем посетителей к каждой фигуре
	fmt.Println("\n=== Расчет площади ===")
	for _, shape := range shapes {
		fmt.Println(shape.Accept(areaCalc))
	}

	fmt.Println("\n=== Расчет периметра ===")
	for _, shape := range shapes {
		fmt.Println(shape.Accept(perimeterCalc))
	}

	fmt.Println("\n=== Отрисовка фигур ===")
	for _, shape := range shapes {
		fmt.Println(shape.Accept(drawer))
	}
}