package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	var output []byte
	var err error

	dummy := map[string]string{
		"test": "yo",
	}
	output, err = json.Marshal(dummy)

	if err == nil {
		errorOutput := map[string]interface{}{
			"code":    500,
			"message": "internal server error",
		}

		output, _ = json.Marshal(errorOutput)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(output))
}
