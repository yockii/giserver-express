package util

import "github.com/bwmarrin/snowflake"

var snowflakeNode *snowflake.Node

func InitNode(node int64) (err error) {
	snowflakeNode, err = snowflake.NewNode(1)
	if err != nil {
		return err
	}
	return nil
}

func SnowflakeId() int64 {
	return snowflakeNode.Generate().Int64()
}
func SnowflakeIdString() string {
	return snowflakeNode.Generate().String()
}
