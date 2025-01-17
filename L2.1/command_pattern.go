package main

import (
	"fmt"
	"time"
)

// Command определяет интерфейс для выполнения операций
type Command interface {
	Execute() error
	Undo() error
}

// TextEditor представляет собой получателя команд - текстовый редактор
type TextEditor struct {
	content string
}

// NewTextEditor создает новый экземпляр текстового редактора
func NewTextEditor() *TextEditor {
	return &TextEditor{}
}

// WriteText добавляет текст в редактор
func (te *TextEditor) WriteText(text string) {
	te.content += text
}

// DeleteText удаляет последние n символов из текста
func (te *TextEditor) DeleteText(n int) {
	if len(te.content) < n {
		te.content = ""
		return
	}
	te.content = te.content[:len(te.content)-n]
}

// GetContent возвращает текущее содержимое редактора
func (te *TextEditor) GetContent() string {
	return te.content
}

// WriteCommand представляет команду для записи текста
type WriteCommand struct {
	editor    *TextEditor
	text      string
	timestamp time.Time
}

// NewWriteCommand создает новую команду записи
func NewWriteCommand(editor *TextEditor, text string) *WriteCommand {
	return &WriteCommand{
		editor:    editor,
		text:      text,
		timestamp: time.Now(),
	}
}

// Execute выполняет команду записи текста
func (c *WriteCommand) Execute() error {
	c.editor.WriteText(c.text)
	return nil
}

// Undo отменяет команду записи текста
func (c *WriteCommand) Undo() error {
	c.editor.DeleteText(len(c.text))
	return nil
}

// DeleteCommand представляет команду для удаления текста
type DeleteCommand struct {
	editor     *TextEditor
	numChars   int
	deletedText string
}

// NewDeleteCommand создает новую команду удаления
func NewDeleteCommand(editor *TextEditor, numChars int) *DeleteCommand {
	return &DeleteCommand{
		editor:   editor,
		numChars: numChars,
	}
}

// Execute выполняет команду удаления текста
func (c *DeleteCommand) Execute() error {
	if len(c.editor.GetContent()) < c.numChars {
		return fmt.Errorf("недостаточно символов для удаления")
	}
	
	// Сохраняем удаляемый текст для возможности отмены
	text := c.editor.GetContent()
	c.deletedText = text[len(text)-c.numChars:]
	c.editor.DeleteText(c.numChars)
	return nil
}

// Undo отменяет команду удаления текста
func (c *DeleteCommand) Undo() error {
	c.editor.WriteText(c.deletedText)
	return nil
}

// CommandInvoker управляет выполнением команд и их отменой
type CommandInvoker struct {
	commands    []Command
	undoStack   []Command
	maxUndoSize int
}

// NewCommandInvoker создает новый инвокер команд
func NewCommandInvoker(maxUndoSize int) *CommandInvoker {
	return &CommandInvoker{
		commands:    make([]Command, 0),
		undoStack:   make([]Command, 0),
		maxUndoSize: maxUndoSize,
	}
}

// ExecuteCommand выполняет команду и добавляет её в историю
func (ci *CommandInvoker) ExecuteCommand(cmd Command) error {
	if err := cmd.Execute(); err != nil {
		return err
	}

	ci.commands = append(ci.commands, cmd)
	ci.undoStack = append(ci.undoStack, cmd)

	// Ограничиваем размер стека отмены
	if len(ci.undoStack) > ci.maxUndoSize {
		ci.undoStack = ci.undoStack[1:]
	}

	return nil
}

// Undo отменяет последнюю выполненную команду
func (ci *CommandInvoker) Undo() error {
	if len(ci.undoStack) == 0 {
		return fmt.Errorf("нет команд для отмены")
	}

	lastIdx := len(ci.undoStack) - 1
	lastCmd := ci.undoStack[lastIdx]
	ci.undoStack = ci.undoStack[:lastIdx]

	return lastCmd.Undo()
}

func main() {
	// Создаем редактор и инвокер команд
	editor := NewTextEditor()
	invoker := NewCommandInvoker(10)

	// Выполняем команды
	writeCmd1 := NewWriteCommand(editor, "Привет, ")
	writeCmd2 := NewWriteCommand(editor, "мир!")
	deleteCmd := NewDeleteCommand(editor, 4)

	fmt.Println("=== Выполнение команд ===")
	
	invoker.ExecuteCommand(writeCmd1)
	fmt.Printf("Текст после первой записи: %q\n", editor.GetContent())

	invoker.ExecuteCommand(writeCmd2)
	fmt.Printf("Текст после второй записи: %q\n", editor.GetContent())

	invoker.ExecuteCommand(deleteCmd)
	fmt.Printf("Текст после удаления: %q\n", editor.GetContent())

	fmt.Println("\n=== Отмена команд ===")
	
	invoker.Undo() // Отмена удаления
	fmt.Printf("Текст после первой отмены: %q\n", editor.GetContent())

	invoker.Undo() // Отмена второй записи
	fmt.Printf("Текст после второй отмены: %q\n", editor.GetContent())

	invoker.Undo() // Отмена первой записи
	fmt.Printf("Текст после третьей отмены: %q\n", editor.GetContent())
}

/*
# Паттерн Команда (Command Pattern) в Go

## Описание паттерна

Команда — это поведенческий паттерн проектирования, который превращает запросы в объекты, позволяя передавать их как аргументы при вызове методов, ставить запросы в очередь, логировать их, а также поддерживать отмену операций.

## Структура паттерна

1. **Command** (Команда) - интерфейс, объявляющий методы выполнения и отмены операций
2. **ConcreteCommand** (Конкретная команда) - реализует интерфейс команды
3. **Receiver** (Получатель) - объект, который выполняет действия
4. **Invoker** (Инициатор) - управляет выполнением команд
5. **Client** (Клиент) - создает и настраивает команды

## Когда использовать паттерн Команда

1. Когда нужно параметризовать объекты выполняемым действием
2. Когда требуется ставить операции в очередь, выполнять их по расписанию или передавать по сети
3. Когда необходима поддержка отмены операций
4. Когда нужно поддерживать логирование изменений и аудит

## Преимущества

1. **Отделение исполнителя от вызывающего кода**
   - Уменьшение связанности
   - Гибкость в выборе исполнителя
   - Улучшенная модульность

2. **Расширяемость**
   - Легкое добавление новых команд
   - Независимость от конкретных исполнителей
   - Простое комбинирование команд

3. **Поддержка отмены операций**
   - Встроенный механизм undo/redo
   - История изменений
   - Восстановление состояния

4. **Отложенное выполнение**
   - Возможность планирования выполнения
   - Поддержка очередей
   - Асинхронное выполнение

## Недостатки

1. **Усложнение кода**
   - Увеличение количества классов
   - Дополнительный уровень абстракции
   - Повышение сложности навигации

2. **Накладные расходы**
   - Затраты памяти на хранение истории
   - Дополнительные объекты команд
   - Возможное снижение производительности

3. **Сложность отмены**
   - Необходимость хранения состояния
   - Сложность реализации для некоторых операций
   - Потенциальные утечки памяти

## Реальные примеры использования

### 1. Система редактирования документов
```go
type DocumentCommand interface {
    Execute() error
    Undo() error
}

type PasteCommand struct {
    document *Document
    position Position
    content  string
}
```

### 2. Управление транзакциями
```go
type TransactionCommand interface {
    Execute() error
    Rollback() error
}

type TransferCommand struct {
    from    *Account
    to      *Account
    amount  decimal.Decimal
}
```

### 3. Управление задачами
```go
type TaskCommand interface {
    Execute() error
    Cancel() error
}

type ScheduledTask struct {
    task      Task
    startTime time.Time
    priority  int
}
```

### 4. Обработка событий UI
```go
type UICommand interface {
    Execute() error
    Revert() error
}

type ButtonClickCommand struct {
    button  *Button
    handler func()
}
```

### 5. Управление конфигурацией
```go
type ConfigCommand interface {
    Apply() error
    Rollback() error
}

type UpdateConfigCommand struct {
    config  *Config
    key     string
    newValue interface{}
    oldValue interface{}
}
```

## Лучшие практики использования

1. **Проектирование команд**
   - Делайте команды неизменяемыми
   - Храните всю необходимую информацию
   - Обеспечивайте атомарность операций

2. **Управление состоянием**
   - Сохраняйте состояние для отмены
   - Используйте снимки состояния
   - Обрабатывайте побочные эффекты

3. **Обработка ошибок**
   - Реализуйте корректную обработку ошибок
   - Обеспечивайте откат при сбоях
   - Поддерживайте согласованность данных

4. **Оптимизация производительности**
   - Ограничивайте размер истории
   - Используйте пул объектов
   - Реализуйте ленивое выполнение

## Антипаттерны при использовании Command

1. **Избыточное состояние**
   - Хранение лишних данных
   - Дублирование информации
   - Неэффективное использование памяти

2. **Нарушение Single Responsibility**
   - Смешивание логики команд
   - Сложные команды с множеством ответственностей
   - Нарушение принципа единственной ответственности

3. **Неправильная обработка ошибок**
   - Игнорирование ошибок
   - Незавершенные откаты
   - Несогласованное состояние

## Заключение

Паттерн Команда является мощным инструментом для работы с операциями в Go. Он особенно полезен когда:
- Требуется поддержка отмены операций
- Необходимо отложенное или асинхронное выполнение
- Нужно логирование и аудит действий
- Требуется параметризация объектов действиями

Для эффективного использования паттерна Команда:
- Тщательно проектируйте интерфейсы команд
- Правильно управляйте состоянием
- Обеспечивайте корректную обработку ошибок
- Оптимизируйте использование памяти
*/