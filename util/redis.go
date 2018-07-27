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
	_, err := pool.Get().Do("ping")
	if err != nil {
		beego.Error(err)
	}
	beego.Info("初始化Redis价格连接池成功.")

	Redis = RedisCli{Pool: pool}
}

// 新建一个连接池
func NewPool(db int) RedisCli {
	dbNum := redis.DialDatabase(db)
	pool := &redis.Pool{
		MaxIdle:     redisMaxIdle,
		MaxActive:   redisMaxActive,
		IdleTimeout: time.Duration(redisIdleTimeout) * time.Second,
		Wait:        redisWait,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
		},
	}
	_, err := pool.Get().Do("ping")
	if err != nil {
		beego.Error("【Redis】", err)
		return RedisCli{}
	}

	return RedisCli{Pool: pool}
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

func (r *RedisCli) LPush(con redis.Conn, key, value string) error {
	_, err := con.Do("lpush", key, value)
	defer r.Close(con)
	return err
}

func (r *RedisCli) Lpop(con redis.Conn, key string) (string, error) {
	res, err := redis.String(con.Do("lpop", key))
	defer r.Close(con)
	return res, err
}

func (r *RedisCli) SAdd(con redis.Conn, key, member string) error {
	_, err := con.Do("sadd", key, member)
	defer r.Close(con)
	return err
}

func (r *RedisCli) SRem(con redis.Conn, key, member string) error {
	_, err := con.Do("srem", key, member)
	defer r.Close(con)
	return err
}

func (r *RedisCli) SMembers(con redis.Conn, key string) ([]string, error) {
	members, err := redis.Strings(con.Do("smembers", key))
	defer r.Close(con)
	return members, err
}

func (r *RedisCli) SIsMember(con redis.Conn, key, member string) (bool, error) {
	exist, err := redis.Bool(con.Do("sismember", key, member))
	defer r.Close(con)
	return exist, err
}

func (r *RedisCli) SCard(con redis.Conn, key string) (int, error) {
	nums, err := redis.Int(con.Do("scard", key))
	defer r.Close(con)
	return nums, err
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

func (cli *Cli) Del(key string) error {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	_, err := con.Do("del", key)
	defer con.Close()
	return err
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

func (cli *Cli) Lpush(list, value string) error {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	_, err := con.Do("lpush", list, value)
	defer con.Close()
	return err
}

func (cli *Cli) Ltrim(list, start, end string) error {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	_, err := con.Do("ltrim", start, end)
	defer con.Close()
	return err
}

func (cli *Cli) Llen(list string) (int, error) {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	llen, err := redis.Int(con.Do("llen", list))
	defer con.Close()
	return llen, err
}

func (cli *Cli) Lindex(list, index string) (string, error) {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	str, err := redis.String(con.Do("lindex", list, index))
	defer con.Close()
	return str, err
}

func (cli *Cli) SetEx(key string, seconds int) error {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	_, err := con.Do("expire", key, seconds)
	defer con.Close()
	return err
}

func (cli *Cli) Exists(key string) (bool, error) {
	dbNum := redis.DialDatabase(cli.DB)
	con, _ := redis.Dial(redisScheme, redisHost+":"+redisPort, dbNum)
	exists, err := redis.Bool(con.Do("exists", key))
	defer con.Close()
	return exists, err
}

func (cli *Cli) Select(db int) {
	cli.DB = db
}
