package model

import (
	"fmt"
	"testing"
)

func TestGoModFilePath(t *testing.T) {

	got := GoModFilePath()
	fmt.Println(got)

}
