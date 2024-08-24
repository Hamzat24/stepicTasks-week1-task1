package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var otstup []bool

func main() {
	out := os.Stdout
	//out, _ := os.Create("file.txt")
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("Can't read directory %s: %s", path, err.Error())
	}

	if !printFiles {
		for i := 0; i < len(dirs); i++ {
			if !dirs[i].IsDir() {
				if i < len(dirs)-1 {
					dirs = append(dirs[:i], dirs[i+1:]...)
					i--
				} else {
					dirs = dirs[:i]
				}
			}
		}
	}
	for k, dir := range dirs {
		for _, v := range otstup {
			if v {
				fmt.Fprint(out, "│	")
			} else {
				fmt.Fprint(out, "	")
			}
		}

		if k < len(dirs)-1 {
			if !dir.IsDir() {
				f, _ := os.Stat(filepath.Join(path, dir.Name()))
				if f.Size() != 0 {
					fmt.Fprintf(out, "├───%s (%db)\n", dir.Name(), f.Size())
				} else {
					fmt.Fprintf(out, "├───%s (empty)\n", dir.Name())
				}
			} else {
				fmt.Fprintf(out, "├───%s\n", dir.Name())
				otstup = append(otstup, true)
			}

		} else {
			if !dir.IsDir() {
				f, _ := os.Stat(filepath.Join(path, dir.Name()))
				if f.Size() != 0 {
					fmt.Fprintf(out, "└───%s (%db)\n", dir.Name(), f.Size())
				} else {
					fmt.Fprintf(out, "└───%s (empty)\n", dir.Name())
				}
			} else {
				fmt.Fprintf(out, "└───%s\n", dir.Name())
				otstup = append(otstup, false)
			}
		}
		if dir.IsDir() {
			dirTree(out, path+"/"+dir.Name(), printFiles)
		}
	}
	if len(otstup) > 1 {
		otstup = otstup[:len(otstup)-1]
	} else {
		otstup = make([]bool, 0)
	}
	return nil
}
