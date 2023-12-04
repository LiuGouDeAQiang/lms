package tools

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

var snowNode *snowflake.Node

func GetUUID() string {
	id := uuid.New() //默认V4版本
	fmt.Printf("uuid:%s,version:%s", id.String(), id.Version().String())
	return id.String()
}
func GetUid() int64 {
	//node, _ := snowflake.NewNode(1)
	//return node.Generate().Int64()
	if snowNode == nil {
		snowNode, _ = snowflake.NewNode(1)
	}

	return snowNode.Generate().Int64()
}
