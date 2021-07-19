package main

import (
	"encoding/json"
	"fmt"
	"io"
	"jacobsmi/server/src/dbutils"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
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

	// Get values necessary for parsing the Excel file
	// Such as where the questions answers and comments are in the file
	questions, _ := strconv.Atoi(r.MultipartForm.Value["questions"][0])
	answers, _ := strconv.Atoi(r.MultipartForm.Value["answers"][0])
	comments, _ := strconv.Atoi(r.MultipartForm.Value["comments"][0])

	h := r.MultipartForm.File["file"][0] // Create a copy of the file that was sent on the server
	file, _ := h.Open()
	tmpfile, _ := os.Create("./" + h.Filename)
	io.Copy(tmpfile, file)
	tmpfile.Close()
	file.Close()

	defer os.Remove("./" + h.Filename)

	// Parse it with excelize
	f, err := excelize.OpenFile(h.Filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	for sheetNum := range f.GetSheetList() {
		sheetName := f.GetSheetName(sheetNum)
		rows, _ := f.GetRows(sheetName)
		for row := range rows {
			question := rows[row][questions-1]
			answer := rows[row][answers-1]
			// Error when comment does not exist
			var comment string
			if len(rows[row]) > comments-1 {
				comment = rows[row][comments-1]
			} else {
				comment = ""
			}
			sqlStatement := `INSERT INTO known_qa (question, answer, comment) VALUES ($1, $2, $3)`
			_, err := dbutils.DB.Exec(sqlStatement, question, answer, comment)
			if err != nil {
				fmt.Println("Error inserting")
				fmt.Println(err)
				return
			}
		}
	}

	// Need to do file parsing here

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

func main() {
	defer dbutils.DB.Close()
	r := mux.NewRouter()

	r.HandleFunc("/parsefile", ParseFile)

	fmt.Println("Running the API at http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST"}))(r)))

}
