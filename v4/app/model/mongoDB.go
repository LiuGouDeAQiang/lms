package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"os"

	"io/ioutil"
	"log"
	"path/filepath"
)

func StoreImagesInMongoDB(imageDir string, collectionName string) {
	// 获取集合

	collection := MB.Database("123").Collection(collectionName)
	// 读取目录下的所有图片文件
	files, err := os.ReadDir(imageDir)
	if err != nil {
		log.Fatal(err)
	}
	// 遍历文件列表并将图片存储到MongoDB中
	for _, file := range files {
		filePath := filepath.Join(imageDir, file.Name())

		// 读取图片文件内容
		imageData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file: %s, error: %s\n", filePath, err)
			continue
		}

		// 创建文档并插入到集合中
		document := bson.M{
			"filename": file.Name(),
			"image":    imageData,
		}
		_, err = collection.InsertOne(context.TODO(), document)
		if err != nil {
			log.Printf("Failed to insert document: %s, error: %s\n", file.Name(), err)
		} else {
			fmt.Printf("Inserted document: %s\n", file.Name())
		}
	}
}

func ReadDocuments(databaseName, collectionName string) {
	// 选择数据库和集合
	database := MB.Database(databaseName)
	collection := database.Collection(collectionName)

	// 查询所有文档
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// 遍历结果集
	for cursor.Next(context.TODO()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			log.Fatal(err)
		}
		// 处理每个文档
		fmt.Println(document)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
func CreateDocument() {
	collection := MB.Database("123").Collection("chat")
	doc := bson.M{"name": "John", "age": 29}
	insertResult, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted document ID:", insertResult.InsertedID)
}
func ReadDocument() {
	collection := MB.Database("123").Collection("chat")
	filter := bson.M{"name": "John"}
	var result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Found document:", result)
}
func UpdateDocument() {
	collection := MB.Database("123").Collection("chat")
	filter := bson.M{"name": "John"}
	update := bson.M{"$set": bson.M{"age": 35}}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified count:", updateResult.ModifiedCount)
}
func DeleteDocument() {
	collection := MB.Database("123").Collection("chat")
	filter := bson.M{"name": "John"}
	deleteResult, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted count:", deleteResult.DeletedCount)
}
