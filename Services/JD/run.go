package main

import (
	"fmt"
	"os"
)

func castDebug(msg string) {
	if DEBUG {
		fmt.Println("DEBUG MESSAGE:", msg)
	}
}

func main() {
	if os.Getenv("DEBUG") != "" {
		fmt.Println("debug mode is on")
		DEBUG = true
	}

	var jd *JD = createJD()
	jd.run()
}
