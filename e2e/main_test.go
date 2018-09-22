// +build integration

package e2e

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	baseTmpDir   = "/tmp/contactbook"
	testUser     = "user"
	testPassword = "password"
)

var (
	binDir string
)

func TestMain(m *testing.M) {

	defBinDir, _ := filepath.Abs("../build")

	flag.StringVar(&binDir, "bindir", defBinDir, "The directory containing contactnook binary")
	flag.Parse()

	if err := os.MkdirAll(baseTmpDir, 0700); err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(baseTmpDir)

	os.Exit(m.Run())
}
