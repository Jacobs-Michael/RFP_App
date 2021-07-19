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
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func processError(w http.ResponseWriter, err error, apiResp string, status int) {
	fmt.Println(err)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "msg": apiResp})
}

func ParseFile(w http.ResponseWriter, r *http.Request) {
	// Parse out the multipart form data
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println("Error getting file")
		fmt.Println(err)
		return
	}

	// Get value of ignoredRows from the Multipart Form
	ignoredRows := r.MultipartForm.Value["ignoredRows"][0]
	ignoredRowsList := strings.Split(ignoredRows, ",")
	// Get other values from the form
	questions, _ := strconv.Atoi(r.MultipartForm.Value["questions"][0])
	answers, _ := strconv.Atoi(r.MultipartForm.Value["answers"][0])
	comments, _ := strconv.Atoi(r.MultipartForm.Value["comments"][0])

	ignoredRowsMap := map[int]bool{}
	for i := range ignoredRowsList {
		ignoredRowsList[i] = strings.ReplaceAll(ignoredRowsList[i], " ", "")
		val, err := strconv.Atoi(ignoredRowsList[i])
		if err != nil {
			processError(w, err, "ignored_rows_err", http.StatusBadGateway)
			return
		}
		ignoredRowsMap[val] = true
	}

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
		// Gets all the sheets in the excel file
		sheetName := f.GetSheetName(sheetNum)
		rows, _ := f.GetRows(sheetName)

		for i := 0; i < len(rows); i++ {

			if _, ok := ignoredRowsMap[i+1]; !ok {
				question := rows[i][questions-1]
				answer := rows[i][answers-1]
				// Error when comment does not exist
				var comment string
				if len(rows[i]) > comments-1 {
					comment = rows[i][comments-1]
				} else {
					comment = ""
				}
				sqlStatement := `INSERT INTO known_qa (question, answer, comment) VALUES ($1, $2, $3)`
				_, err := dbutils.DB.Exec(sqlStatement, question, answer, comment)
				if err != nil {
					if err, ok := err.(*pq.Error); ok {
						if err.Code != "23505" {
							processError(w, err, "db_insert_err", http.StatusInternalServerError)
							return
						}
					}
				}
			} else {
				fmt.Printf("Ignoring row %d \n", i)
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
