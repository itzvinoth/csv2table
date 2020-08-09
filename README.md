# csv2table (Not finished)

Command-line utility to convert CSV files to searchable and
sortable HTML table. Supports large datasets and horizontal scrolling for large number of columns. This is just a copycat of ([csvtotable](https://github.com/vividvilla/csvtotable/)) utility written in python. Only for learning purpose &mdash; doing this utility library in golang.

**Can't use for now.**

## Run command 
`csv2table -csv=example.csv --serve` 


## Usage 
```
-h or --help              Used to print the list of arguments
--serve                   Open html output in a web browser.
-csv=table.csv            Mention csv file name.
-save=output.html         Mention html filename to save the output, .
-dl=20                    Number of rows to show by default. Defaults to 20 (show all rows when -1 or 0).
-pagination=false         Enable/disable pagination. Enabled by default.
-colvis=true              Enable/disable Column visualisation. Disabled by default.
```
