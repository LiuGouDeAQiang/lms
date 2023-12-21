package model

import "testing"

func TestCreateDocument(t *testing.T) {
	NewMongoDB()
	CreateDocument()

}
func TestStoreImagesInMongoDB(t *testing.T) {
	NewMongoDB()
	imageDir := "F:\\go.code\\src\\go_code\\lms\\app\\img"
	collectionName := "chat"
	StoreImagesInMongoDB(imageDir, collectionName)
}
