package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbxb/shorty/internal/app/config"
	"github.com/sbxb/shorty/internal/app/logger"
	"github.com/sbxb/shorty/internal/app/storage"
	"github.com/sbxb/shorty/internal/app/storage/psql"
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

	url, err := uh.store.GetURL(r.Context(), id)
	if err != nil {
		if IsDeletedError(err) {
			http.Error(w, "Record deleted", http.StatusGone)
		} else {
			http.Error(w, "Server failed to process URL", http.StatusInternalServerError)
		}
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
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Server failed to read the request's body", http.StatusInternalServerError)
		return
	}

	url := string(b)

	if !u.IsValidInputURL(url) {
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

	if err := uh.store.AddBatchURL(r.Context(), respBatch, userID); err != nil {
		http.Error(w, "Server failed to store URL(s)", http.StatusInternalServerError)
		return
	}

	for i := range respBatch {
		respBatch[i].ShortURL = uh.config.BaseURL + "/" + respBatch[i].ShortURL
	}

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

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// is request an empty struct
	if req == (u.URLRequest{}) || !u.IsValidInputURL(req.URL) {
		http.Error(w, "Bad request: non-valid object received", http.StatusBadRequest)
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

// UserDeleteHandler process DELETE /api/user/urls request
// ... асинхронный хендлер DELETE /api/user/urls, который принимает список
// идентификаторов сокращённых URL для удаления в формате:
// [ "a", "b", "c", "d", ...]
// В случае успешного приёма запроса хендлер должен возвращать HTTP-статус
// 202 Accepted. Фактический результат удаления может происходить позже -
// каким-либо образом оповещать пользователя об успешности или
// неуспешности не нужно.
// Успешно удалить URL может пользователь, его создавший.
// При запросе удалённого URL с помощью хендлера GET /{id} нужно вернуть
// статус 410 Gone.
func (uh URLHandler) UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var deleteIDs []string

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&deleteIDs); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(deleteIDs) == 0 {
		http.Error(w, "Bad request: no ids provided", http.StatusBadRequest)
		return
	}

	userID := GetUserID(r.Context())

	go func() {
		if uh.config.DatabaseDSN == "" {
			// Either MapStorage or FileMapStorage is used (for testing
			// purposes only), delete batch straightforward since map-based
			// storage can not really benefit from concurrency due to heavy
			// locking and fullscan
			err := uh.store.DeleteBatch(context.Background(), deleteIDs, userID)
			if err != nil {
				logger.Warningf("UserDeleteHandler : DeleteBatch failed: %v", err)
			}
		} else {
			// Real Database is used, increment 14 requires some concurrency here
			ConcurrentDeleteBatch(uh.store, deleteIDs, userID)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}

// UserGetHandler process GET /user/urls request
// ... хендлер GET /api/user/urls, который сможет вернуть пользователю все
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

	userID := GetUserID(r.Context())

	urls, _ := uh.store.GetUserURLs(r.Context(), userID)

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
	dbStore, ok := uh.store.(*psql.DBStorage)
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
