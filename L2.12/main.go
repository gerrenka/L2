// main.go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"
)

// Event представляет событие в календаре
type Event struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Date        time.Time `json:"date"`
}

// CalendarService представляет бизнес-логику календаря
type CalendarService struct {
    events map[int]Event
    nextID int
}

// NewCalendarService создает новый экземпляр сервиса календаря
func NewCalendarService() *CalendarService {
    return &CalendarService{
        events: make(map[int]Event),
        nextID: 1,
    }
}

// CreateEvent создает новое событие
func (s *CalendarService) CreateEvent(userID int, title, description string, date time.Time) (Event, error) {
    if title == "" {
        return Event{}, fmt.Errorf("title cannot be empty")
    }

    event := Event{
        ID:          s.nextID,
        UserID:      userID,
        Title:       title,
        Description: description,
        Date:        date,
    }

    s.events[s.nextID] = event
    s.nextID++

    return event, nil
}

// UpdateEvent обновляет существующее событие
func (s *CalendarService) UpdateEvent(id, userID int, title, description string, date time.Time) error {
    event, exists := s.events[id]
    if !exists {
        return fmt.Errorf("event not found")
    }

    if event.UserID != userID {
        return fmt.Errorf("unauthorized")
    }

    event.Title = title
    event.Description = description
    event.Date = date

    s.events[id] = event
    return nil
}

// DeleteEvent удаляет событие
func (s *CalendarService) DeleteEvent(id, userID int) error {
    event, exists := s.events[id]
    if !exists {
        return fmt.Errorf("event not found")
    }

    if event.UserID != userID {
        return fmt.Errorf("unauthorized")
    }

    delete(s.events, id)
    return nil
}

// GetEventsForDay возвращает события на указанный день
func (s *CalendarService) GetEventsForDay(userID int, date time.Time) []Event {
    var result []Event
    for _, event := range s.events {
        if event.UserID == userID && isSameDay(event.Date, date) {
            result = append(result, event)
        }
    }
    return result
}

// GetEventsForWeek возвращает события на указанную неделю
func (s *CalendarService) GetEventsForWeek(userID int, date time.Time) []Event {
    var result []Event
    weekStart := date.AddDate(0, 0, -int(date.Weekday()))
    weekEnd := weekStart.AddDate(0, 0, 7)

    for _, event := range s.events {
        if event.UserID == userID && event.Date.After(weekStart) && event.Date.Before(weekEnd) {
            result = append(result, event)
        }
    }
    return result
}

// GetEventsForMonth возвращает события на указанный месяц
func (s *CalendarService) GetEventsForMonth(userID int, date time.Time) []Event {
    var result []Event
    for _, event := range s.events {
        if event.UserID == userID && 
           event.Date.Year() == date.Year() && 
           event.Date.Month() == date.Month() {
            result = append(result, event)
        }
    }
    return result
}

// Handler представляет HTTP обработчик
type Handler struct {
    service *CalendarService
    logger  *log.Logger
}

// LoggingMiddleware реализует middleware для логирования запросов
func (h *Handler) LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        h.logger.Printf("Started %s %s", r.Method, r.URL.Path)
        next(w, r)
        h.logger.Printf("Completed in %v", time.Since(start))
    }
}

// parseDate парсит дату из строки
func parseDate(dateStr string) (time.Time, error) {
    return time.Parse("2006-01-02", dateStr)
}

// isSameDay проверяет, относятся ли две даты к одному дню
func isSameDay(date1, date2 time.Time) bool {
    y1, m1, d1 := date1.Date()
    y2, m2, d2 := date2.Date()
    return y1 == y2 && m1 == m2 && d1 == d2
}

// writeJSON отправляет JSON-ответ
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// handleCreateEvent обрабатывает создание события
func (h *Handler) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
        return
    }

    userID, err := strconv.Atoi(r.Form.Get("user_id"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
        return
    }

    date, err := parseDate(r.Form.Get("date"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
        return
    }

    event, err := h.service.CreateEvent(
        userID,
        r.Form.Get("title"),
        r.Form.Get("description"),
        date,
    )

    if err != nil {
        writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{"result": event})
}

// handleUpdateEvent обрабатывает обновление события
func (h *Handler) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
        return
    }

    eventID, err := strconv.Atoi(r.Form.Get("event_id"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid event_id"})
        return
    }

    userID, err := strconv.Atoi(r.Form.Get("user_id"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
        return
    }

    date, err := parseDate(r.Form.Get("date"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
        return
    }

    err = h.service.UpdateEvent(
        eventID,
        userID,
        r.Form.Get("title"),
        r.Form.Get("description"),
        date,
    )

    if err != nil {
        writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
        return
    }

    writeJSON(w, http.StatusOK, map[string]string{"result": "event updated"})
}

// handleDeleteEvent обрабатывает удаление события
func (h *Handler) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
        return
    }

    eventID, err := strconv.Atoi(r.Form.Get("event_id"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid event_id"})
        return
    }

    userID, err := strconv.Atoi(r.Form.Get("user_id"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
        return
    }

    err = h.service.DeleteEvent(eventID, userID)
    if err != nil {
        writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
        return
    }

    writeJSON(w, http.StatusOK, map[string]string{"result": "event deleted"})
}

// handleEventsForDay обрабатывает получение событий за день
func (h *Handler) handleEventsForDay(w http.ResponseWriter, r *http.Request) {
    userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
        return
    }

    date, err := parseDate(r.URL.Query().Get("date"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
        return
    }

    events := h.service.GetEventsForDay(userID, date)
    writeJSON(w, http.StatusOK, map[string]interface{}{"result": events})
}

// handleEventsForWeek обрабатывает получение событий за неделю
func (h *Handler) handleEventsForWeek(w http.ResponseWriter, r *http.Request) {
    userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
        return
    }

    date, err := parseDate(r.URL.Query().Get("date"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
        return
    }

    events := h.service.GetEventsForWeek(userID, date)
    writeJSON(w, http.StatusOK, map[string]interface{}{"result": events})
}

// handleEventsForMonth обрабатывает получение событий за месяц
func (h *Handler) handleEventsForMonth(w http.ResponseWriter, r *http.Request) {
    userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
        return
    }

    date, err := parseDate(r.URL.Query().Get("date"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
        return
    }

    events := h.service.GetEventsForMonth(userID, date)
    writeJSON(w, http.StatusOK, map[string]interface{}{"result": events})
}

func main() {
    service := NewCalendarService()
    logger := log.New(log.Writer(), "HTTP: ", log.LstdFlags)
    handler := &Handler{service: service, logger: logger}

    // Регистрация обработчиков
    http.HandleFunc("/create_event", handler.LoggingMiddleware(handler.handleCreateEvent))
    http.HandleFunc("/update_event", handler.LoggingMiddleware(handler.handleUpdateEvent))
    http.HandleFunc("/delete_event", handler.LoggingMiddleware(handler.handleDeleteEvent))
    http.HandleFunc("/events_for_day", handler.LoggingMiddleware(handler.handleEventsForDay))
    http.HandleFunc("/events_for_week", handler.LoggingMiddleware(handler.handleEventsForWeek))
    http.HandleFunc("/events_for_month", handler.LoggingMiddleware(handler.handleEventsForMonth))

    port := ":8080" // Порт можно вынести в конфиг
    logger.Printf("Starting server on port %s", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        logger.Fatal(err)
    }
}