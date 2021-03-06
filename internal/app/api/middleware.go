package api

import (
	"compress/gzip"
	"context"
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

func gzipMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			r.Body = gz
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func authMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		uid := ""
		if err != nil || !auth.CheckUserIDCookieValue(cookie.Value) {
			uid, err = auth.GenerateUserID()
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			cookie := http.Cookie{
				Name:    "user_id",
				Value:   auth.GetUserIDCookieValue(uid),
				Expires: time.Now().Add(1 * time.Hour),
			}
			http.SetCookie(w, &cookie)
		} else {
			uid = cookie.Value[:32]
		}
		ctx := context.WithValue(r.Context(), auth.ContextUserIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func jsonEncMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const ContentType = "application/json"

		if strings.ToLower(r.Header.Get("Content-Type")) != ContentType {
			http.Error(w, "Content-Type should be "+ContentType, http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}
