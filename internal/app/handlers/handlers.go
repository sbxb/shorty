package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/storage"
	u "github.com/sbxb/shorty/internal/app/url"
)

// URLHandler defines a container for handlers and their dependencies
type URLHandler struct {
	store  storage.Storage
	config config.Config
}

func NewURLHandler(st storage.Storage, cfg config.Config) URLHandler {
	return URLHandler{
		store:  st,
		config: cfg,
	}
}

// GetHandler process GET /{id} request
// ... Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор
// сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL
// в HTTP-заголовке Location ...
func (uh URLHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, err := uh.store.GetURL(id)
	if err != nil {
		http.Error(w, "Server failed to process URL", http.StatusInternalServerError)
		return
	}
	if url == "" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// PostHandler process POST / request
// ... Эндпоинт POST / принимает в теле запроса строку URL для сокращения
// и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой
// строки в теле ...
func (uh URLHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	var reader io.Reader
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}
	b, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, "Server failed to read the request's body", http.StatusInternalServerError)
		return
	}

	url := string(b)
	// TODO There should be some kind of URL validation
	if url == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	status := http.StatusCreated

	userID := GetUserID(r.Context())

	ue := u.URLEntry{
		ShortURL:    u.ShortID(url),
		OriginalURL: url,
	}

	err = uh.store.AddURL(r.Context(), ue, userID)

	if IsConflictError(err) {
		status = http.StatusConflict
	} else if err != nil {
		http.Error(w, "Server failed to store URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, "%s/%s", uh.config.BaseURL, ue.ShortURL)
}

// JSONBatchPostHandler process POST /api/shorten/batch request with JSON array payload
// ... хендлер POST /api/shorten/batch, принимающий в теле запроса множество
// URL для сокращения в формате:
// [
//    {
//        "correlation_id": "<строковый идентификатор>",
//        "original_url": "<URL для сокращения>"
//    }, ...
// ]
//
// В качестве ответа хендлер должен возвращать данные в формате:
//
// [
//    {
//        "correlation_id": "<строковый идентификатор из объекта запроса>",
//        "short_url": "<результирующий сокращённый URL>"
//    }, ...
// ]
func (uh URLHandler) JSONBatchPostHandler(w http.ResponseWriter, r *http.Request) {
	const ContentType = "application/json"
	batch := []u.BatchURLRequestEntry{}

	if r.Header.Get("Content-Type") != ContentType {
		http.Error(w, "Bad request: Content-Type should be "+ContentType, http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&batch); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// check if request array contains empty or incomplete records
	if !u.IsBatchURLRequestValid(batch) {
		http.Error(w, "Bad request: empty or incomplete records received", http.StatusBadRequest)
		return
	}

	userID := GetUserID(r.Context())

	// we're ready to start processing
	respBatch := make([]u.BatchURLEntry, 0, len(batch))
	for _, entry := range batch {
		ne := u.BatchURLEntry{
			CorrelationID: entry.CorrelationID,
			OriginalURL:   entry.OriginalURL,
			ShortURL:      u.ShortID(entry.OriginalURL),
		}
		respBatch = append(respBatch, ne)
	}

	// Storage staff starts here
	if err := uh.store.AddBatchURL(r.Context(), respBatch, userID); err != nil {
		http.Error(w, "Server failed to store URL(s)", http.StatusInternalServerError)
		return
	}
	// Storage stuff stops here

	for i := range respBatch {
		respBatch[i].ShortURL = uh.config.BaseURL + "/" + respBatch[i].ShortURL
	}

	//w.Write([]byte(fmt.Sprintf("%+v", respBatch)))
	jr, err := json.Marshal(respBatch)
	if err != nil {
		http.Error(w, "Server failed to process response result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write(jr)
}

// JSONPostHandler process POST /api/shorten request with JSON payload
// ... эндпоинт POST /api/shorten, принимающий в теле запроса JSON-объект
// {"url": "<some_url>"} и возвращающий в ответ объект {"result": "<shorten_url>"}
func (uh URLHandler) JSONPostHandler(w http.ResponseWriter, r *http.Request) {
	const ContentType = "application/json"
	var req u.URLRequest

	if r.Header.Get("Content-Type") != ContentType {
		http.Error(w, "Bad request: Content-Type should be "+ContentType, http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// is request an empty struct
	if req == (u.URLRequest{}) {
		http.Error(w, "Bad request: empty object received", http.StatusBadRequest)
		return
	}

	status := http.StatusCreated

	userID := GetUserID(r.Context())

	ue := u.URLEntry{
		ShortURL:    u.ShortID(req.URL),
		OriginalURL: req.URL,
	}

	err := uh.store.AddURL(r.Context(), ue, userID)

	if IsConflictError(err) {
		status = http.StatusConflict
	} else if err != nil {
		http.Error(w, "Server failed to store URL", http.StatusInternalServerError)
		return
	}

	jr, err := json.Marshal(
		u.URLResponse{
			Result: fmt.Sprintf("%s/%s", uh.config.BaseURL, ue.ShortURL),
		},
	)

	if err != nil {
		http.Error(w, "Server failed to process response result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(status)
	w.Write(jr)
}

// UserGetHandler process GET /user/urls request
// ... хендлер GET /user/urls, который сможет вернуть пользователю все
// когда-либо сокращённые им URL в формате:
// [
//     {
//         "short_url": "http://...",
//         "original_url": "http://..."
//     }, ...
// ]
// При отсутствии сокращённых пользователем URL хендлер должен отдавать
// HTTP-статус 204 No Content ...
func (uh URLHandler) UserGetHandler(w http.ResponseWriter, r *http.Request) {
	const ContentType = "application/json"

	uid := GetUserID(r.Context())

	urls, _ := uh.store.GetUserURLs(uid)

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	for i := range urls {
		urls[i].ShortURL = uh.config.BaseURL + "/" + urls[i].ShortURL
	}
	jr, err := json.Marshal(urls)

	if err != nil {
		http.Error(w, "Server failed to process response result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(http.StatusOK)
	w.Write(jr)
}

// PingGetHandler process GET /ping request
// ... хендлер GET /ping, который при запросе проверяет соединение с базой
// данных. При успешной проверке хендлер должен вернуть HTTP-статус 200 OK,
// при неуспешной — 500 Internal Server Error ...
func (uh URLHandler) PingGetHandler(w http.ResponseWriter, r *http.Request) {
	dbStore, ok := uh.store.(*storage.DBStorage)
	if !ok {
		http.Error(w, "Server failed to open DB", http.StatusInternalServerError)
		return
	}

	if err := dbStore.Ping(); err != nil {
		http.Error(w, "Server failed to ping DB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}
