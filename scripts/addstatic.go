package main

import (
	"io"
	"io/ioutil"
	"os"
)

func main() {
	fs, _ := ioutil.ReadDir("../static")
	out, _ := os.Create("static.go")
	out.Write([]byte("package main \n\nconst (\n"))
	for _, f := range fs {
		out.Write([]byte(f.Name() + " = `"))
		f, _ := os.Open(f.Name())
		io.Copy(out, f)
		out.Write([]byte("`\n"))
	}
	out.Write([]byte(")\n"))
}
