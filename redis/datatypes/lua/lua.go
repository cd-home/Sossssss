package lua

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
	"redis/client"
)

var luaScript = `
	local value = redis.call("Get", KEYS[1])
	if( value - KEYS[2] >= 0 ) then
		local leftStock = redis.call("DecrBy" , KEYS[1],KEYS[2])
	return leftStock
	else
		return value - KEYS[2]
	end
	return -1
`

var StockSha string

func init() {
	ctx := context.Background()
	cli := client.NewSimpleClient()
	StockSha, _ = cli.ScriptLoad(ctx, luaScript).Result()
}

func RunLuaScript(luaScript string) {
	ctx := context.Background()
	cli := client.NewSimpleClient()
	//cli.Set(ctx, "stock", "10", redis.KeepTTL)
	luaCmd := redis.NewScript(luaScript)
	r, err := luaCmd.Run(ctx, cli, []string{"stock", "6"}).Result()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(r.(int64))
}

func EvalLuaSha() {
	ctx := context.Background()
	cli := client.NewSimpleClient()
	r, err := cli.EvalSha(ctx, StockSha, []string{"stock", "6"}).Result()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(r.(int64))
}
