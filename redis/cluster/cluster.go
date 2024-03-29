package cluster

import (
	"context"
	"github.com/go-redis/redis/v9"
	"log"
)

func NewClusterClient() *redis.ClusterClient {
	// if readOnly == true   slave
	// if readOnly == true and RouteByLatency 选择对应slot 延迟最低的master 或者 slave
	// if readOnly == true and RouteRandomly  随机选择对应slot master 或者 slave
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"10.211.55.61:6379", // master
			"10.211.55.61:6380", // slave
			"10.211.55.62:6379", // master
			"10.211.55.62:6380", // slave
			"10.211.55.63:6379", // master
			"10.211.55.63:6380", // slave
		},
		ReadOnly:       true,
		RouteByLatency: false,
		RouteRandomly:  false,
		Username:       "webapp",        // 此用户有key权限限制 acl list 查看
		Password:       "webapp2023@xp", // 本地测试集群的fake密码
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		log.Println(err)
		return nil
	}
	return cli
}
