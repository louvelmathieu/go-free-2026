package api

import (
	"errors"
	"net/http"

	"gitlab.example.com/free/provisio/internal/domain"
)

// ErrorResponse est le format d'erreur de l'API.
type ErrorResponse struct {
	Error     string `json:"error"`
	Field     string `json:"field,omitempty"`
	Step      string `json:"step,omitempty"`
	Retryable bool   `json:"retryable,omitempty"`
}

// writeError traduit une erreur métier en réponse HTTP. C'est la frontière :
// à l'intérieur on manipule des erreurs Go, à l'extérieur des codes HTTP.
func writeError(w http.ResponseWriter, err error) {
	resp := ErrorResponse{Error: "erreur interne"}
	status := http.StatusInternalServerError

	var verr *domain.ValidationError
	var serr *domain.StepError

	switch {
	case errors.As(err, &verr):
		status = http.StatusBadRequest
		resp = ErrorResponse{Error: verr.Reason, Field: verr.Field}

	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		resp = ErrorResponse{Error: "commande introuvable"}

	case errors.Is(err, domain.ErrNotEligible):
		status = http.StatusUnprocessableEntity
		resp = ErrorResponse{Error: "adresse non éligible"}

	case errors.Is(err, domain.ErrNoStock):
		status = http.StatusConflict
		resp = ErrorResponse{Error: "aucun équipement disponible"}

	case errors.As(err, &serr):
		// Une erreur d'étape non reconnue : on expose l'étape et si
		// la commande peut être rejouée, pas le message interne.
		status = http.StatusBadGateway
		resp = ErrorResponse{Error: "échec d'activation", Step: serr.Step, Retryable: domain.IsRetryable(err)}
	}

	writeJSON(w, status, resp)
}
