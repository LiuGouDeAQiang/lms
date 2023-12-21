package model

import (
	"fmt"
	"testing"
)

func TestGetAdminList(t *testing.T) {
	NewMysql()
	admin := make([]Admin, 0)
	admin = GetAdminList(5, 0)
	fmt.Println(admin)
}
