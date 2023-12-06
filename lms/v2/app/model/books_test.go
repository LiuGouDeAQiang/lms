package model

import (
	"fmt"
	"testing"
)

func TestGetImg(t *testing.T) {
	NewMysql()
	img, err := GetImg("红与黑")
	if err != nil {
		return
	}
	fmt.Println("你大爷" + img)
	Close()
}
