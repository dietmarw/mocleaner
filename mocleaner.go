package main

import "os"
import "strings"
import "bufio"

/* This function takes a string and returns
   a (potentially nil) error object */
func TTWS(filename string) error {
	/* Open the input file */
	inf, err := os.Open(filename);
	/* In case this function generates a "panic", be sure to close this file */
	defer inf.Close();
	/* Did we open it successfully?  If not, close and return. */
	if (err!=nil) { inf.Close(); return err; }

	/* Open the output file */
	outf, err := os.Create("stripped.txt");
	/* In case this function generates a "panic", be sure to close this file */
	defer outf.Close();
	/* Did we open it succesfully?  If not, close all and return. */
	if (err!=nil) { inf.Close(); outf.Close(); return err; }

	/* Create a scanner object to break this in to lines */
	scanner := bufio.NewScanner(inf)
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

	/* No errors, so we return nil */
	return nil;
}

func main() { TTWS("test.mo"); }
