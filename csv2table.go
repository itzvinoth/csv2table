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
	// "reflect"
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

	var display_len int
	// Get the arguments value from the termial
	args := args.GetArgs()

	input_file, _ := args["csv"].(string)
	output_file, _ := args["save"].(string)
	serve, _ := args["serve"].(bool)
	dl, _ := args["dl"].(int)
	pagination, _ := args["pagination"].(bool)
	colvis, _ := args["colvis"].(bool)

	// In case input file is present, we can serve the file in the browser
	if (input_file != "empty") {
		serve = true
	}
	if (output_file != "empty") {
		serve = false
	}
	if (dl != 0 && dl != -1) {
		display_len = dl
	} else {
		display_len = -1
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
	
	// Render csv to HTML
	if serve == true {
		tmpl := template.Must(template.ParseFiles("table.html"))
		
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data := TableData{
				Header: csv_headers,
				Rows: csv_rows,
				PageLength: display_len,
				Pagination: pagination,
				Colvis: colvis,
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
		displaylen := string(display_len)
		spagination := strconv.FormatBool(pagination)

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
					data: `+rowstr+`,
					columns: `+columndata+`,
					pageLength: `+displaylen+`,
					paging: `+spagination+`
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
