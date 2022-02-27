package helpers

import (
	"testing"
)

func AssertStatusCode(t *testing.T, expected, got int) {
	if expected != got {
		t.Errorf("Expect response code to be %v, but got %v instead", expected, got)
	}
}
