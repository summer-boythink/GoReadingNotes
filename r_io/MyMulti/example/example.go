package main

import (
	"MyMulti"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	r := strings.NewReader("qqq")
	var buf1, buf2 bytes.Buffer
	w := MyMulti.MultiWriter(&buf1, &buf2)
	io.Copy(w, r)
	fmt.Println(buf2.String(), buf1)

	r1 := strings.NewReader("11")
	r2 := strings.NewReader("22")
	rr := MyMulti.MultiReader(r1, r2)
	io.Copy(os.Stdout, rr)
}
