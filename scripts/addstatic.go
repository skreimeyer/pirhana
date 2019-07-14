package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	fs, _ := ioutil.ReadDir("../static")
	out, _ := os.Create("../pirhana/static.go")
	out.Write([]byte("package main \n\nconst (\n"))
	for _, f := range fs {
		out.Write([]byte(f.Name() + " = `"))
		f, err := os.Open("../static/" + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		w, err := io.Copy(out, f)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(w)
		}

		out.Write([]byte("`\n\n"))
	}
	out.Write([]byte(")\n"))
}
