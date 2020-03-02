package handlers

import (
    "fmt"
    "github.com/gorilla/mux"
    "github.com/rs/zerolog/log"
    "lfs-postcode/api/services"
    "net/http"
)

type PostCodeMatchHandler struct {
    services.PostCodeMatchService
}

func NewPostCodeMatchHandler() *PostCodeMatchHandler {
    return &PostCodeMatchHandler{services.PostCodeMatchService{}}
}

func (p PostCodeMatchHandler) PostCodeMatchingHandler(w http.ResponseWriter, r *http.Request) {

    log.Debug().
        Str("client", r.RemoteAddr).
        Str("uri", r.RequestURI).
        Msg("Received match postcode request")

    vars := mux.Vars(r)
    year := vars["year"]
    month := vars["month"]

    // Convert year to int
    yr, err := IntConversion(year)
    if err != nil {
        ErrorResponse{
            Status:       Error,
            ErrorMessage: fmt.Sprintf("invalid year: %s, expected an integer", year)}.sendResponse(w, r)
        return
    }

    // Convert month to int
    mo, err := IntConversion(month)
    if err != nil || (mo < 1 || mo > 12) {
        ErrorResponse{
            Status:       Error,
            ErrorMessage: fmt.Sprintf("invalid month: [%d], expected one of 1-12", mo)}.sendResponse(w, r)
        return
    }

    res, err := p.PostCodeDifferences(yr)
    if err != nil {
        ErrorResponse{Status: Error, ErrorMessage: err.Error()}.sendResponse(w, r)
        return
    }

    if res == nil || len(res) == 0 {
        MatchingResponse{
            Status:  Success,
            Message: "Postcodes match",
        }.sendResponse(w, r)
        return
    }

    SendDataResponse{}.sendResponse(w, r, res)
}

