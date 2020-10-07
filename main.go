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
	"strconv"
	"reflect"
	// "bytes"
	"github.com/mevinoth/csv2table/args"
)

type TableData struct {
	Header []string
	Rows [][]string
	PageLength int
	Pagination bool
	Colvis bool
}

func main() {

	var displayLength int
	// Get the arguments value from the termial
	args := args.GetArgs()

	inputCsvFile, _ := args["csv"].(string)
	outputHtmlFile, _ := args["save"].(string)
	serve, _ := args["serve"].(bool)
	dl, _ := args["dl"].(int)
	pagination, _ := args["pagination"].(bool)
	colvis, _ := args["colvis"].(bool)

	// In case input file is present, we can serve the file in the browser
	if (inputCsvFile != "empty") {
		serve = true
	}
	if (outputHtmlFile != "empty") {
		serve = false
	}
	if (dl != 0 && dl != -1) {
		displayLength = dl
	} else {
		displayLength = -1
	}
	
	// fmt.Println("reflect", reflect.TypeOf(serve).Kind())
	
	// Open the file
	csvfile, err := os.Open(inputCsvFile)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	var csvHeaders []string
	var csvRows [][]string

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
			csvHeaders = record
		} else {
			csvRows = append(csvRows, record)
		}
		i = i + 1
	}
	
	// Render csv to HTML
	if serve == true {
		tmpl := template.Must(template.ParseFiles("templates/table.html"))
		
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data := TableData{
				Header: csvHeaders,
				Rows: csvRows,
				PageLength: displayLength,
				Pagination: pagination,
				Colvis: colvis,
			}
			tmpl.Execute(w, data)
		})
		fmt.Println("Listening on  http://localhost:8080")
		http.ListenAndServe(":8080", nil)
	} else if outputHtmlFile != "empty" {
		var columnHeaders []map[string]string

		for _, v := range csvHeaders {
			colHeader := make(map[string]string)
			colHeader["title"] = v
			columnHeaders = append(columnHeaders, colHeader)
		}
		
		// Marshal the map into a JSON string.
		jsonRowDataBytes, err := json.Marshal(csvRows)
		jsonColDataBytes, err := json.Marshal(columnHeaders)

		t := 
`<!DOCTYPE html>
<html lang="en">
	<head>
		<link href="https://cdn.datatables.net/1.10.12/css/jquery.dataTables.css" rel="stylesheet" />
		<link href="https://cdn.datatables.net/buttons/1.2.2/css/buttons.dataTables.css" rel="stylesheet" />
		<script src="https://code.jquery.com/jquery-1.12.4.js"></script>
		<script src="https://cdn.datatables.net/1.10.16/js/jquery.dataTables.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/jszip/3.1.3/jszip.min.js"></script>
		<script src="https://cdn.datatables.net/buttons/1.4.2/js/dataTables.buttons.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/pdfmake/0.1.32/pdfmake.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/pdfmake/0.1.32/vfs_fonts.js"></script>
		<script src="https://cdn.datatables.net/buttons/1.4.2/js/buttons.html5.min.js"></script>
		<script src="https://cdn.datatables.net/buttons/1.6.2/js/buttons.print.min.js"></script>
	</head>
	<body>
		<div class="container">
			<table id="table" class="display nowrap" width="100%"></table>
		</div>
		<script>
			$(document).ready( function () {
				$("#table").DataTable({
					dom: 'Bfrtip',
					buttons: [{
						extend: 'print',
						title: 'Customized Print Title',
						filename: 'customized_print_file_name'
					}, {
						extend: 'copy',
						title: 'Customized Copy Title',
						filename: 'customized_copy_file_name'
					}, {
						extend: 'pdf',
						title: 'Customized PDF Title',
						filename: 'customized_pdf_file_name'
					}, {
						extend: 'excel',
						title: 'Customized EXCEL Title',
						filename: 'customized_excel_file_name'
					}, {
						extend: 'csv',
						filename: 'customized_csv_file_name'
					}],
					data: `+string(jsonRowDataBytes)+`,
					columns: `+string(jsonColDataBytes)+`,
					pageLength: `+strconv.Itoa(displayLength)+`,
					paging: `+strconv.FormatBool(pagination)+`
				})
			});
		</script>
	</body>
</html>
`

		html := []byte(t)

		if err = ioutil.WriteFile(outputHtmlFile, html, 0666); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
