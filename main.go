package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"encoding/csv"
	"encoding/json"
	"html/template"
	"net/http"
	"io/ioutil"
	// "reflect"
	// "bytes"
	"github.com/mevinoth/csv2table/args"
)

type TableData struct {
	Header []string
	Rows [][]string
}

func main() {

	// Get the arguments value from the termial
	args := args.GetArgs()

	input_file, _ := args["csv"].(string)
	output_file, _ := args["save"].(string)
	serve, _ := args["serve"].(bool)
	// In case input file is present, we can serve the file in the browser
	if (input_file != "empty") {
		serve = true
	}
	if (output_file != "empty") {
		serve = false
	}
	// fmt.Println("reflect", reflect.TypeOf(serve).Kind())
	
	// Open the file
	csvfile, err := os.Open(input_file)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	var csv_headers []string
	var csv_rows [][]string

	// index value 
	i := 1

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// first row pushes to the csv headers, remaining pushes to the csv rows
		if i == 1 {
			csv_headers = record
		} else {
			csv_rows = append(csv_rows, record)
		}
		i = i + 1
	}
	// Read from csv and get header and body contents


	// Render csv to HTML
	if serve == true {
		tmpl := template.Must(template.ParseFiles("table.html"))
		
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data := TableData{
				Header: csv_headers,
				Rows: csv_rows,
			}
			tmpl.Execute(w, data)
		})
		fmt.Println("Listening on  http://localhost:8080")
		http.ListenAndServe(":8080", nil)
	} else if output_file != "empty" {
		var columnHeaders []map[string]string

		for _, v := range csv_headers {
			ch := make(map[string]string)
			ch["title"] = v
			columnHeaders = append(columnHeaders, ch)
		}
		
		// Marshal the map into a JSON string.
		jsonRowData, err := json.Marshal(csv_rows)
		rowstr := string(jsonRowData)
		jsonColData, err := json.Marshal(columnHeaders)
		columndata := string(jsonColData)

		t := 
`<!DOCTYPE html>
<html lang="en">
	<head>
		<script type="text/javascript" charset="utf8" src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.5.1/jquery.min.js"></script>
		<link href="https://cdn.datatables.net/1.10.21/css/jquery.dataTables.min.css" rel="stylesheet" type="text/css">
		<script src="https://cdn.datatables.net/1.10.21/js/jquery.dataTables.min.js" type="text/javascript"></script>
	</head>
	<body>
		<div>
			<table id="table"></table>
		</div>
		<script>
			$(document).ready( function () {
				$("#table").DataTable({
					data: `+rowstr+`,
					columns: `+columndata+`
				})
			});
		</script>
	</body>
</html>
`

		html := []byte(t)

		if err = ioutil.WriteFile(output_file, html, 0666); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
