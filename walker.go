package main

import (
  "path/filepath"
  "os"
  "flag"
  "fmt"
  "net/http"
  "io"
  "io/ioutil"
  "strings"
)

func WalkFunc(path string, info os.FileInfo, err error) error {
	blacklist := []string{".bzr", ".cvs", ".git", ".hg", ".svn"}
	if contains(path, blacklist){
		fmt.Printf("Skipping version control dir: %s\n", path)
		return filepath.SkipDir
	} else {
		inf, err := os.Open(path)
		defer inf.Close();
		if (err!=nil) { inf.Close(); return err; }
		readStart := io.LimitReader(inf, 512);
		data, err := ioutil.ReadAll(readStart);
		fileType := http.DetectContentType(data);
		if strings.Contains(fileType, "text/plain"){
			fmt.Printf("Trimming: %v\n", path)
		} else {
			fmt.Printf("Skipping file of type: %v: %v\n", fileType, path)
		}
	}
	return nil
}

func contains(x string, a []string) bool {
	for _, e := range(a) {
		if (x==e) { return true; }
	}
	return false;
}


func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, WalkFunc)
	fmt.Printf("filepath.Walk() returned %v\n", err)
}
