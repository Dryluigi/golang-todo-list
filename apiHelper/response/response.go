package response

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BuildErrorResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Status:  code,
		Error:   message,
	})
}

func BuildServerErrorResponse(w http.ResponseWriter, err error) {
	godotenv.Load(".env")

	message := "internal server error"

	if debugMode, errParse := strconv.ParseBool(os.Getenv("DEBUG_MODE")); errParse == nil && debugMode {
		message = err.Error()
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Status:  http.StatusInternalServerError,
		Error:   message,
	})
}

func BuildSuccessResponse(w http.ResponseWriter, code int, message string, isSuccess bool, data interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: isSuccess,
		Status:  code,
		Message: message,
		Data:    data,
	})
}
