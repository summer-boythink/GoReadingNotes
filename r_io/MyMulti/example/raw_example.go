package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	r := strings.NewReader("aaa qqq")
	var buf1, buf2 bytes.Buffer
	w := io.MultiWriter(&buf1, &buf2)
	if _, err := io.Copy(w, r); err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf1.String(), buf2)

	r1 := strings.NewReader("11")
	r2 := strings.NewReader("22")
	rr := io.MultiReader(r1, r2)
	io.Copy(os.Stdout, rr)
}
