package cmd

var (
	// The api file
	api string
	// The target dir
	dir string
	// The table name
	table string
	// The data source of database,like "root:password@tcp(127.0.0.1:3306)/database
	url string
	// The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority
	home string
	// describes the style of output files.
	style string
)
