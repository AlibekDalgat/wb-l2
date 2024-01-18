package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

// Event представляет событие в календаре.
type Event struct {
	UserId  int       `json:"user_id"`
	EventID int       `json:"event_id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
}

func newErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	errorResponse := map[string]string{"error": err.Error()}
	jsonResponse, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

// Промежуточное ПО для логирования запросов.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func getUserId(r *http.Request) (int, error) {
	// Извлечение значения параметра "user_id" из queryString
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		return 0, errors.New(`Отсутствует параметр "user_id" в queryString`)
	}

	// Преобразование параметра "user_id" в целое число
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, errors.New(`Ошибка при преобразовании "user_id" в число: ` + err.Error())
	}
	return userID, nil
}

func getDate(r *http.Request) (time.Time, error) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		return time.Time{}, errors.New(`Отсутствует параметр "user_id" в queryString`)
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, errors.New("ошибка парсинга значения времени")
	}
	return date, nil
}

// Обработчик HTTP для /create_event
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		newErrorResponse(w, errors.New("Метод не поддерживается"), http.StatusMethodNotAllowed)
		return
	}
	userID, err := getUserId(r)
	if err != nil {
		newErrorResponse(w, err, http.StatusBadRequest)
	}

	// Чтение тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при чтении тела запроса: "+err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Декодирование JSON-данных в структуру Event
	var event Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при декодировании JSON: "+err.Error()), http.StatusInternalServerError)
		return
	}
	event.UserId = userID

	// Создание JSON-ответа
	eventRes, err := createEventService(event)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при работе сервиса: "+err.Error()), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(eventRes)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при создании JSON-ответа: "+err.Error()), http.StatusInternalServerError)
		return
	}

	// Установка заголовков и отправка JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Обработчик HTTP для /update_event
func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		newErrorResponse(w, errors.New("Метод не поддерживается"), http.StatusMethodNotAllowed)
		return
	}
	userID, err := getUserId(r)
	if err != nil {
		newErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Чтение тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при чтении тела запроса: "+err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Декодирование JSON-данных в структуру Event
	var updateIvent Event
	err = json.Unmarshal(body, &updateIvent)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при декодировании JSON: "+err.Error()), http.StatusInternalServerError)
		return
	}
	updateIvent.UserId = userID

	// Создание JSON-ответа
	eventRes, err := updateEventService(updateIvent)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при работе сервиса: "+err.Error()), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(eventRes)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при создании JSON-ответа: "+err.Error()), http.StatusInternalServerError)
		return
	}

	// Установка заголовков и отправка JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Обработчик HTTP для /delete_event
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		newErrorResponse(w, errors.New("Метод не поддерживается"), http.StatusMethodNotAllowed)
		return
	}
	userID, err := getUserId(r)
	if err != nil {
		newErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Чтение тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при чтении тела запроса: "+err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Декодирование JSON-данных в структуру Event
	var deleteIvent Event
	err = json.Unmarshal(body, &deleteIvent)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при декодировании JSON: "+err.Error()), http.StatusInternalServerError)
		return
	}
	deleteIvent.UserId = userID

	// Создание JSON-ответа
	err = deleteEventService(deleteIvent)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при работе сервиса: "+err.Error()), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(map[string]string{"result": "ok"})
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при создании JSON-ответа: "+err.Error()), http.StatusInternalServerError)
		return
	}

	// Установка заголовков и отправка JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Обработчик HTTP для /events_for_day
func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		newErrorResponse(w, errors.New("Метод не поддерживается"), http.StatusMethodNotAllowed)
		return
	}
	userID, err := getUserId(r)
	if err != nil {
		newErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	date, err := getDate(r)
	if err != nil {
		newErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Создание JSON-ответа
	eventsRes, err := getEventsForRangeService(userID, date)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при работе сервиса: "+err.Error()), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(eventsRes)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при создании JSON-ответа: "+err.Error()), http.StatusInternalServerError)
		return
	}

	// Установка заголовков и отправка JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Обработчик HTTP для /events_for_week
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		newErrorResponse(w, errors.New("Метод не поддерживается"), http.StatusMethodNotAllowed)
		return
	}
	userID, err := getUserId(r)
	if err != nil {
		newErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	date, err := getDate(r)
	if err != nil {
		newErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Создание JSON-ответа
	eventsRes, err := getEventsForRangeService(userID, date)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при работе сервиса: "+err.Error()), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(eventsRes)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при создании JSON-ответа: "+err.Error()), http.StatusInternalServerError)
		return
	}

	// Установка заголовков и отправка JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Обработчик HTTP для /events_for_month
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		newErrorResponse(w, errors.New("Метод не поддерживается"), http.StatusMethodNotAllowed)
		return
	}
	userID, err := getUserId(r)
	if err != nil {
		newErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	date, err := getDate(r)
	if err != nil {
		newErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Создание JSON-ответа
	eventsRes, err := getEventsForRangeService(userID, date)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при работе сервиса: "+err.Error()), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(eventsRes)
	if err != nil {
		newErrorResponse(w, errors.New("Ошибка при создании JSON-ответа: "+err.Error()), http.StatusInternalServerError)
		return
	}

	// Установка заголовков и отправка JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	// Настройте обработчики HTTP
	http.HandleFunc("/create_event", createEventHandler)
	http.HandleFunc("/update_event", updateEventHandler)
	http.HandleFunc("/delete_event", deleteEventHandler)
	http.HandleFunc("/events_for_day", eventsForDayHandler)
	http.HandleFunc("/events_for_week", eventsForWeekHandler)
	http.HandleFunc("/events_for_month", eventsForMonthHandler)

	// Настройте промежуточное ПО для логирования
	http.Handle("/", loggingMiddleware(http.DefaultServeMux))

	// Укажите порт из конфигурации
	port := ":8080"

	// Запустите сервер
	log.Printf("Сервер слушает порт %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// работа бизнес логики
func createEventService(event Event) (Event, error) {
	if event.Title == "" && event.Date.IsZero() {
		return Event{}, errors.New("Отсутствуют обязательные поля в событии")
	}
	return Event{EventID: 1, UserId: event.UserId, Title: event.Title, Date: event.Date}, nil
}

func updateEventService(event Event) (Event, error) {
	if event.Title == "" && event.Date.IsZero() {
		return Event{}, errors.New("Отсутствуют редактируемые поля события")
	}
	return Event{EventID: 1, UserId: event.UserId, Title: event.Title, Date: event.Date}, nil
}

func deleteEventService(event Event) error {
	if event.EventID == 0 {
		return errors.New("Отсутствует выбранная заметка")
	}
	return nil
}

func getEventsForRangeService(userID int, date time.Time) ([]Event, error) {
	return []Event{{EventID: 1, UserId: userID, Title: "title", Date: date}, {EventID: 2, UserId: userID, Title: "title", Date: date}}, nil
}
