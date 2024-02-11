package main

import (
	"envi/cmd"
	"fmt"
	"os"
	// A simple example that shows how to render an animated progress bar. In this
	// example we bump the progress by 25% every two seconds, animating our
	// progress bar to its new target state.
	//
	// It's also possible to render a progress bar in a more static fashion without
	// transitions. For details on that approach see the progress-static example.
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
