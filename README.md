# csv2table (Not finished)

Command-line utility written in golang ([python](https://github.com/vividvilla/csvtotable/)) to convert CSV files to searchable and
sortable HTML table. Supports large datasets and horizontal scrolling for large number of columns.

## Run command 
`csv2table -csv=aapl.csv --serve` 


## Usage 
```
-h or --help              Used to print the list of arguments
	--serve                   Open html output in a web browser.
	-csv=table.csv            Mention csv file name.
	-save=output.html         Mention html filename to save the output, .
```