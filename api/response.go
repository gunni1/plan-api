package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int
	Data   interface{}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewResponse(status int, data interface{}) Response {
	return Response{
		Status: status,
		Data:   data,
	}
}

func (r Response) SendJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/JSON; charset=UTF-8")
	w.WriteHeader(r.Status)

	if err := json.NewEncoder(w).Encode(r.Data); err != nil {
		status := http.StatusInternalServerError
		w.WriteHeader(status)
		r.Status = status
		r.Data = ErrorResponse{Error: err.Error()}

		_ = json.NewEncoder(w).Encode(r.Data)
	}
}

func SendErrorJSON(status int, errorKey string, w http.ResponseWriter) {
	NewResponse(status, ErrorResponse{errorKey}).SendJSON(w)
}
