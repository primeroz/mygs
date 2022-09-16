package main

import (
	"flag"
)

func main() {
	// =====================
	// Get OS parameter
	// =====================
	var bind string
	flag.StringVar(&bind, "bind", "0.0.0.0:9104", "bind")
	flag.Parse()
}
