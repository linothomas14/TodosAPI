package helpers

import ( 
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func HttpResponse(w http.ResponseWriter, httpStatus int, response Response) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)
}