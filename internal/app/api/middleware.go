package api

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sbxb/shorty/internal/app/auth"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipWrapper(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func cookieAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil || !auth.CheckUserIdCookieValue(cookie.Value) {
			uid, err := auth.GenerateUserId()
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			cookie := http.Cookie{
				Name:    "user_id",
				Value:   auth.GetUserIdCookieValue(uid),
				Expires: time.Now().Add(1 * time.Hour),
			}
			http.SetCookie(w, &cookie)
		}
		next.ServeHTTP(w, r)
	})
}
