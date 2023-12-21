package model

import (
	"fmt"
	"testing"
)

func TestNewTelephone(t *testing.T) {
	code, err := NewTelephone("15993855695")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err)
	fmt.Println(code)
}
func TestNewMongoDB(t *testing.T) {
	NewMongoDB()
}
