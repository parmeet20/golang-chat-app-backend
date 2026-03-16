package healthcheck

import (
	"net/http"
	"time"

	"github.com/parmeet20/golang-chatapp/pkg/response"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	resp := HealthResponse{
		Status:    "ok",
		Service:   "golang-chatapp",
		Timestamp: time.Now().UTC(),
	}
	response.JSON(w, http.StatusOK, resp)
}
