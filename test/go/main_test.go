package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.RemoveAll("../../tmp")
	os.RemoveAll("../examples/output")
	
	code := m.Run()
	
	os.RemoveAll("../../tmp")
	os.RemoveAll("../examples/output")
	
	os.Exit(code)
}