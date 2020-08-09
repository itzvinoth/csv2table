package args

import (
	"os"
	"flag"
	"fmt"
	// "reflect"
)

func GetArgs () map[string]interface{} {
	var h bool
	flag.BoolVar(&h, "h", false, "Used to print the list of arguments")

	var help bool
	flag.BoolVar(&help, "help", false, "Used to print the list of arguments")

	var csv string
	flag.StringVar(&csv, "csv", "empty", "Mention csv filename -- e.g., -csv=table.csv")

	var save string
	flag.StringVar(&save, "save", "empty", "Mention html filename to save the output -- e.g., -save=output.html")
	
	var serve bool
	flag.BoolVar(&serve, "serve", false, "Open html output in a web browser")
	
	var dl int
	flag.IntVar(&dl, "dl", 20, "Number of rows to show by default. Defaults to 20 (show all rows when -1 and 0)")

	var pagination bool
	flag.BoolVar(&pagination, "pagination", true, "Enable/disable pagination. Enabled by default.")

	var colvis bool
	flag.BoolVar(&colvis, "colvis", false, "Enable/disable pagination. Enabled by default.")

	// to execute the command-line parsing
	flag.Parse()
	args := make(map[string]interface{})
	args["help"] = help
	args["h"] = h
	args["csv"] = csv
	args["save"] = save
	args["serve"] = serve
	args["dl"] = dl
	args["pagination"] = pagination
	args["colvis"] = colvis

	if help || h {
		fmt.Fprintf(os.Stderr, 
`
	-h or --help              Used to print the list of arguments
	--serve                   Open html output in a web browser.
	-csv=table.csv            Mention csv file name.
	-save=output.html         Mention html filename to save the output, .
	-dl=25                    Number of rows to show by default. Defaults to -1 (show all rows).
	-pagination=false         Enable/disable pagination. Enabled by default.
	-colvis=true              Enable/disable pagination. Disabled by default.
`)
		os.Exit(2)
	}

	if csv == "empty" {
		fmt.Fprintf(os.Stderr, "missing required -%s flag, check more with -h or --help\n", "csv")
        os.Exit(2) // the same exit code flag.Parse uses
	}
	// fmt.Println("svar: ", save, serve, reflect.TypeOf(save).Kind())

	// argsWithoutProg := os.Args[1:]
	return args
}
