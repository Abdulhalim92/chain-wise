package health

import (
	"net/http"
)

// Handler возвращает HTTP-обработчик для проверки здоровья (200 OK). Ошибка Write игнорируется (best-effort).
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
