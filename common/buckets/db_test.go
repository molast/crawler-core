package buckets

import (
	"os"
	"testing"
)

// Ensure we can open/close a buckets db.
func TestOpen(t *testing.T) {
	bx, err := Open(tempfile())
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(bx.Path())
	defer bx.Close()
}
