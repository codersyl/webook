//go:build k8s

package config

var Config = WebookConfig{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:11309)/webook",
	},
	Redis: RedisConfig{
		Addr:     "localhost:6379",
		DB:       1,
		Password: "",
	},
}
