package tools

import (
	"fmt"
	"testing"
)

func TestGetUid(t *testing.T) {
	id := GetUid()
	fmt.Println(id)
}
