package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ParseFile(w http.ResponseWriter, r *http.Request) {
	// Parse out the multipart form data
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println("Error getting file")
		fmt.Println(err)
		return
	}

	// Get the file data and create a copy on the backend
	for _, h := range r.MultipartForm.File["file"] {
		file, _ := h.Open()
		tmpfile, _ := os.Create("./" + h.Filename)
		io.Copy(tmpfile, file)
		tmpfile.Close()
		file.Close()
	}

	// Need to do file parsing here

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/parsefile", ParseFile)

	fmt.Println("Running the API at http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST"}))(r)))

}
