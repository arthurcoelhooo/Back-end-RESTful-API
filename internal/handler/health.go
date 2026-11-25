package handler

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	escreverJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
