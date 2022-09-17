package main

import (
	"MyPipe"
	"fmt"
	"io"
	"os"
)

func main() {
	r, w := MyPipe.Pipe()
	go func() {
		fmt.Fprintln(w, "aaa")
		fmt.Fprintln(w, "qwe")
		w.Close()
	}()
	io.Copy(os.Stdout, r)
}
