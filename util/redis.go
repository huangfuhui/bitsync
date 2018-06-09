package util

import (
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego"
	"time"
)

var Redis RedisCli
var err error

type RedisCli struct {
	Cli redis.Conn
}

func init() {
	redisScheme := beego.AppConfig.String("redis_scheme")
	redisHost := beego.AppConfig.String("redis_host")
	redisPort := beego.AppConfig.String("redis_port")
	redisMaxActive, _ := beego.AppConfig.Int("redis_max_active")
	redisMaxIdle, _ := beego.AppConfig.Int("redis_max_idle")
	redisIdleTimeout, _ := beego.AppConfig.Int("redis_idle_timeout")
	redisWait, _ := beego.AppConfig.Bool("redis_wait")

	dbNum := redis.DialDatabase(1)
	pool := redis.Pool{
		MaxIdle:     redisMaxIdle,
		MaxActive:   redisMaxActive,
		IdleTimeout: time.Duration(redisIdleTimeout) * time.Second,
		Wait:        redisWait,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
		},
	}
	if err != nil {
		beego.Error(err)
	} else {
		beego.Info("初始化Redis连接池.")
	}

	Redis = RedisCli{Cli: pool.Get()}
}

// 切换数据库
func (r *RedisCli) DB(db int64) error {
	_, err := r.Cli.Do("select", db)
	return err
}

// 获取字符串
func (r *RedisCli) Get(key string) (string, error) {
	res, err := redis.String(r.Cli.Do("get", key))
	return res, err
}

// 设置字符串
func (r *RedisCli) Set(key, value string) error {
	_, err := r.Cli.Do("set", key, value)
	return err
}

// 设置字符串(含过期时间设置)
func (r *RedisCli) SetEx(key, value, expire string) error {
	_, err := r.Cli.Do("set", key, value, "ex", expire)
	return err
}
