package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)


/* This function takes a string and returns
   a (potentially nil) error object */
func TTWS(filename string) error {
	/* Open the input file */
	inf, err := os.Open(filename);
	/* In case this function generates a "panic", be sure to close this file */
	defer inf.Close();
	/* Did we open it successfully?  If not, close and return. */
	if (err!=nil) { return err; }

	/* Open the output file in system temp dir*/
	outf, err := ioutil.TempFile("","");
	/* In case this function generates a "panic", be sure to close this file */
	defer outf.Close();
	/* Did we open it succesfully?  If not, close all and return. */
	if (err!=nil) { inf.Close(); return err; }

	/* Create a scanner object to break this in to lines */
	scanner := bufio.NewScanner(inf);
	/* Declare a variable for the line */
	var line string;
	/* Loop over lines */
	for scanner.Scan() {
		/* Trim right space and then add the \n back on the end before writing */
		line = strings.TrimRight(scanner.Text(), " \t")+"\n"
		outf.Write([]byte(line));
	}
	/* Close all open files */
	inf.Close();
	outf.Close();

	/* Replace the source file by the trimmed file */
	os.Rename(outf.Name(), filename);

	/* No errors, so we return nil */
	return nil;
}

func WalkFunc(path string, fi os.FileInfo, err error) error {
	/* list of directories to ignore */
	blacklist := []string{".bzr", ".cvs", ".git", ".hg", ".svn"}
	if contains(path, blacklist){
		fmt.Printf("Skipping version control dir: %s\n", path)
		return filepath.SkipDir
	} else {
		inf, err := os.Open(path)
		defer inf.Close();
		if (err!=nil) { return err; }
		readStart := io.LimitReader(inf, 512);
		data, err := ioutil.ReadAll(readStart);
		/* Close all open files */
		inf.Close();

		if (err!=nil) { return err; }

		/* Determine file type */
		fileType := http.DetectContentType(data);

		if (fi.IsDir()) return nil; // Now you don't need ot check this

		/* only act on text files */
		if (strings.Contains(fileType, "text/plain")) {
			fmt.Printf("Trimming: %v\n", path);
			TTWS(path);
		} else {
			fmt.Printf("Skipping file of type '%v': %v\n", fileType, path)
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
