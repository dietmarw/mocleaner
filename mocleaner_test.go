package main

import "testing"

import "os"
import "fmt"
import "path"
import "io/ioutil"

var sample = "ABC\n  CDE\n  FGH  \nIJK  \n L M N\n";
var trimmed = "ABC\n  CDE\n  FGH\nIJK\n L M N\n";

func read(filename string, t *testing.T) string {
	r, err := os.Open(filename);
	if (err!=nil) { t.Fatal("Couldn't open file "+filename+" for reading: "+err.Error()); }
	defer r.Close();
	data, err := ioutil.ReadAll(r);
	if (err!=nil) { t.Fatal("Couldn't read file "+filename+": "+err.Error()); }
	r.Close();

	return string(data);
}

func TestTTWS(t *testing.T) {
	/* Create a temporary directory */
	tdir, err := ioutil.TempDir("", "");
	defer os.RemoveAll(tdir);
	if (err!=nil) { t.Fatal("Couldn't create temporary directory: "+err.Error()); }

	/* Create a temporary file in this directory */
	f1 := path.Join(tdir, "temp1.mo");

	/* Write some stuff to it */
	w, err := os.Create(f1);
	if (err!=nil) { t.Fatal("Couldn't open temp file for writing: "+err.Error()); }
	defer w.Close();
	w.Write([]byte(sample));
	w.Close();

	data := read(f1, t);

	if (string(data)!=sample) { t.Fatal("Initial string didn't match"); }

	err = TTWS(f1, true);
	if (err!=nil) { t.Fatal("Error from TTWS: "+err.Error()); }

	data = read(f1, t);

	if (string(data)!=trimmed) {
		fmt.Println("'"+trimmed+"'");
		fmt.Println("vs.");
		fmt.Println("'"+string(data)+"'");
		t.Fatal("Trimmed string didn't match");
	}

	os.RemoveAll(tdir);
}

func TestWalking(t *testing.T) {
}
