package main

import (
	"envi/cmd"
	"envi/internal/ui"
	"fmt"
	"os"
	"time"
	// A simple example that shows how to render an animated progress bar. In this
	// example we bump the progress by 25% every two seconds, animating our
	// progress bar to its new target state.
	//
	// It's also possible to render a progress bar in a more static fashion without
	// transitions. For details on that approach see the progress-static example.
)

func main() {
	prg := ui.ProgressBarProgram()
	time.AfterFunc(3*time.Second, func() {
		prg.Send(ui.MsgComplete{})
	})

	_, err := prg.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
