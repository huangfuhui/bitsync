package util

import (
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego"
	"time"
)

var Redis RedisCli

type RedisCli struct {
	Pool *redis.Pool
}

var (
	redisScheme         = beego.AppConfig.String("redis_scheme")
	redisHost           = beego.AppConfig.String("redis_host")
	redisPort           = beego.AppConfig.String("redis_port")
	redisMaxActive, _   = beego.AppConfig.Int("redis_max_active")
	redisMaxIdle, _     = beego.AppConfig.Int("redis_max_idle")
	redisIdleTimeout, _ = beego.AppConfig.Int("redis_idle_timeout")
	redisWait, _        = beego.AppConfig.Bool("redis_wait")
)

func init() {
	dbNum := redis.DialDatabase(1)
	pool := &redis.Pool{
		MaxIdle:     redisMaxIdle,
		MaxActive:   redisMaxActive,
		IdleTimeout: time.Duration(redisIdleTimeout) * time.Second,
		Wait:        redisWait,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
		},
	}
	beego.Info("初始化Redis连接池.")
	_, err := pool.Get().Do("ping")
	if err != nil {
		beego.Error(err)
	}

	Redis = RedisCli{Pool: pool}
}

// 切换数据库
func (r *RedisCli) DB(con redis.Conn, db int) error {
	_, err := con.Do("select", db)
	return err
}

// 获取字符串
func (r *RedisCli) Get(con redis.Conn, key string) (string, error) {
	res, err := redis.String(con.Do("get", key))
	defer r.Close(con)
	return res, err
}

// 设置字符串
func (r *RedisCli) Set(con redis.Conn, key, value string) error {
	_, err := con.Do("set", key, value)
	defer r.Close(con)
	return err
}

// 设置字符串(含过期时间设置)
func (r *RedisCli) SetEx(con redis.Conn, key, value, expire string) error {
	_, err := con.Do("set", key, value, "ex", expire)
	defer r.Close(con)
	return err
}

// 获取连接
func (r *RedisCli) Con() redis.Conn {
	return r.Pool.Get()
}

// 释放连接
func (r *RedisCli) Close(con redis.Conn) error {
	return con.Close()
}

// 关闭连接池
func (r *RedisCli) ClosePool(con redis.Conn) error {
	return r.Pool.Close()
}

type Cli struct {
	DB int
	Ex int
}

func (cli *Cli) Get(key string) (string, error) {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	res, err := redis.String(con.Do("get", key))
	defer con.Close()
	return res, err
}

func (cli *Cli) Set(key, value string) error {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	_, err := con.Do("set", key, value)
	defer con.Close()
	return err
}

func (cli *Cli) SetEx(key, value string) error {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	_, err := con.Do("set", key, value, "ex", cli.Ex)
	defer con.Close()
	return err
}

func (cli *Cli) Select(db int) {
	cli.DB = db
}
