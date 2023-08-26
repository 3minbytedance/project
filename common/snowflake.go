package common

import (
	"github.com/bwmarrin/snowflake"
)

var (
	Node *snowflake.Node
)

// InitSnowflake 初始化生成器
func InitSnowflake(nodeId int64) (err error) {
	Node, err = snowflake.NewNode(nodeId)
	return
}

func GetUid() (id uint) {
	return uint(Node.Generate().Int64())
}
