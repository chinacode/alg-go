package demo

import (
	"context"
	"github.com/garyburd/redigo/redis"
	//"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

var RedisClient *redis.Pool

var RedisConf = map[string]string{
	"name":    "redis",
	"type":    "tcp",
	"address": "192.168.1.200",
	"auth":    "123456",
}

func ExampleClient() {
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle: 16, //最初的连接数量
		//MaxActive:1000000,    //最大连接数量
		MaxActive:   0,                 //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300 * time.Second, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			c, err := redis.Dial(RedisConf["type"], RedisConf["address"])
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", RedisConf["auth"]); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	// 从池里获取连接
	rc := RedisClient.Get()
	// 用完后将连接放回连接池
	//defer rc.Close()

	//key := "table_1"
	start := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		//rdb.SetBit(ctx, key, int64(i), 1)
		//rc.GetBit(ctx, key, int64(i))
		//rc.Do("SETBIT", key, i)
		//v, _ := redis.Int64(rc.Do("GETBIT", key, i))
		v, _ := redis.String(rc.Do("GET", "aaaa"))
		println(v)
	}

	println((time.Now().UnixNano() - start) / 1000 / 1000)
}

/*func ExampleClientRedisGo() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.200:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})

	key := "table_1"
	start := time.Now().UnixNano()
	for i := 0; i < 100000; i++ {
		//rdb.SetBit(ctx, key, int64(i), 1)
		rdb.GetBit(ctx, key, int64(i))
	}

	println((time.Now().UnixNano() - start) / 1000 / 1000)

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	//Output: key value
	//key2 does not exist
}*/
