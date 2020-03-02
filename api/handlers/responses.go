package handlers

import (
    "encoding/json"
    "github.com/rs/zerolog/log"
    "net/http"
    "strconv"
)

const (
    Error   = "ERROR"
    OK      = "OK"
    Success = "SUCCESS"
)

type Response interface {
    sendResponse(w http.ResponseWriter, r *http.Request)
}

type SendDataResponse struct{}

func (re SendDataResponse) sendResponse(w http.ResponseWriter, r *http.Request, d interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(d); err != nil {
        log.Error().
            Str("client", r.RemoteAddr).
            Str("uri", r.RequestURI).
            Msg("json.NewEncoder() failed in sendDataResponse")
    }
}

type ErrorResponse struct {
    Status       string `json:"status"`
    ErrorMessage string `json:"errorMessage"`
}

type MatchingResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

func (response MatchingResponse) sendResponse(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    sendResponse(w, r, response)
}

type OkayResponse struct {
    Status string `json:"status"`
}

func (response OkayResponse) sendResponse(w http.ResponseWriter, r *http.Request) {
    response.Status = OK
    w.WriteHeader(http.StatusOK)
    sendResponse(w, r, response)
}

func (response ErrorResponse) sendResponse(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    response.Status = Error
    w.WriteHeader(http.StatusBadRequest)
    sendResponse(w, r, response)
}

func sendResponse(w http.ResponseWriter, r *http.Request, response Response) {
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Error().
            Str("client", r.RemoteAddr).
            Str("uri", r.RequestURI).
            Msg("json.NewEncoder() failed in PostCodeMatchHandler")
    }
}

func IntConversion(str string) (int, error) {
    i, err := strconv.Atoi(str)
    if err != nil {
        return 0, err
    }
    return i, nil
}

