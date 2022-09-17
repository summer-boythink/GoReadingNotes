package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	r, w := io.Pipe()
	go func() {
		fmt.Fprintln(w, "ww")
		w.Close()
	}()

	io.Copy(os.Stdout, r)
}
