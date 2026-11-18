package httpx

import (
	"net"
	"net/http"
	"time"
)

// NewClient renvoie un client HTTP avec des timeouts explicites.
// http.DefaultClient n'en a aucun : un serveur distant muet bloque
// l'appelant indéfiniment. Chaque service de l'équipe passe par ici.
func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: timeout,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
		},
	}
}

// WithRequestID copie l'identifiant de requête du contexte dans la requête
// sortante, pour suivre un traitement de service en service.
func WithRequestID(req *http.Request, requestID string) *http.Request {
	if requestID != "" {
		req.Header.Set(HeaderRequestID, requestID)
	}
	return req
}
