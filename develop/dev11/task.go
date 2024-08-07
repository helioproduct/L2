package main

import (
	"dev11/endpoints"
	"dev11/logic"
	"dev11/middleware"
	"dev11/models"
	"flag"
	"fmt"
	"log"
	"net/http"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API:
	POST /create_event
	POST /update_event
	POST /delete_event
	GET /events_for_day
	GET /events_for_week
	GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3.  В случае ошибки бизнес-логики сервер должен возвращать HTTP 503.
		В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400.
		В случае остальных ошибок сервер должен возвращать HTTP 500.
		Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

type muxBuilder struct {
	mux *http.ServeMux
}

func NewMuxBuilder() *muxBuilder {
	return &muxBuilder{
		mux: http.NewServeMux(),
	}
}

func (m *muxBuilder) AddEventHTTP(e *endpoints.EventHTTP) {
	m.mux.Handle("/create_event", middleware.WithMethod("POST", http.HandlerFunc(e.CreateHandle)))
	m.mux.Handle("/update_event", middleware.WithMethod("POST", http.HandlerFunc(e.UpdateHandle)))
	m.mux.Handle("/delete_event", middleware.WithMethod("POST", http.HandlerFunc(e.DeleteHandle)))

	m.mux.Handle("/events_for_day", middleware.WithMethod("GET", http.HandlerFunc(e.ForDayHandle)))
	m.mux.Handle("/events_for_week", middleware.WithMethod("GET", http.HandlerFunc(e.ForWeekHandle)))
	m.mux.Handle("/events_for_month", middleware.WithMethod("GET", http.HandlerFunc(e.ForMonthHandle)))
}

func (m *muxBuilder) Build() http.Handler {
	// Добавляем логирование запросов
	withLogger := middleware.Logging(m.mux)
	return withLogger
}

type config struct {
	addr string
}

func parseConfig() *config {
	host := flag.String("h", "127.0.0.1", "адрес, который будет прослушиваться")
	port := flag.Int("p", 8080, "порт, на котором будет запущен http сервер")
	flag.Parse()
	return &config{
		addr: fmt.Sprintf("%s:%d", *host, *port),
	}
}

func main() {
	// Получаем конфиги
	cfg := parseConfig()

	// Модель, позволяющая взаимодействовать с БД событий
	eventModel := models.NewEventModelMemory()
	// Слой бизнес-логики
	eventApi := logic.NewEventAPI(eventModel)
	// http ручки
	eventHTTP := endpoints.NewEventHTTP(eventApi)
	log.Println("eventHTTP ready")

	// Строим мультиплексер
	mb := NewMuxBuilder()
	mb.AddEventHTTP(eventHTTP)
	// Получаем его
	mux := mb.Build()
	log.Println("mux build")

	// Настраиваем сервер
	log.Printf("HTTP ListenAndServe: %s\n", cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
