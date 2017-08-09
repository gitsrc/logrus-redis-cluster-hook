package logredis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	SINGLE int = iota
	CLUSTER
)

// HookConfig stores configuration needed to setup the hook
type HookConfig struct {
	Key        string
	Format     string
	App        string
	Addrs      []string
	Password   string
	Hostname   string
	ConnOption int
	DB         int
}

// RedisHook to sends logs to Redis server
type RedisHook struct {
	RedisCluster   *redis.ClusterClient
	RedisClient    *redis.Client
	RedisAddrs     []string
	RedisKey       string
	LogstashFormat string
	AppName        string
	Hostname       string
	RedisPort      int
}

// NewHook creates a hook to be added to an instance of logger
func NewHook(config HookConfig) (*RedisHook, error) {
	var hook = &RedisHook{
		RedisAddrs: config.Addrs,
		RedisKey:   config.Key,
		AppName:    config.App,
		Hostname:   config.Hostname,
	}

	var err error

	if config.ConnOption == SINGLE {
		hook.RedisClient, err = newSingleConnHook(&config)
	} else {
		hook.RedisCluster, err = newClusterHook(&config)
	}

	if err != nil {
		return nil, err
	}

	if config.Format != "v0" && config.Format != "v1" {
		return nil, fmt.Errorf("unknown message format")
	} else {
		hook.LogstashFormat = config.Format
	}

	return hook, nil
}

// Fire is called when a log event is fired.
func (hook *RedisHook) Fire(entry *logrus.Entry) error {
	var msg interface{}

	switch hook.LogstashFormat {
	case "v0":
		msg = createV0Message(entry, hook.AppName, hook.Hostname)
	case "v1":
		msg = createV1Message(entry, hook.AppName, hook.Hostname)
	default:
		fmt.Println("Invalid LogstashFormat")
	}

	// Marshal into json message
	js, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error creating message for REDIS: %s", err)
	}

	// get connection and send msg
	if hook.RedisClient != nil {
		c := hook.RedisClient

		// send message
		_, err := c.RPush(hook.RedisKey, js).Result()
		if err != nil {
			return fmt.Errorf("error sending message to REDIS: %s", err)
		}
	} else {
		cc := hook.RedisCluster

		// send message
		_, err := cc.RPush(hook.RedisKey, js).Result()
		if err != nil {
			return fmt.Errorf("error sending message to REDIS: %s", err)
		}
	}

	return nil
}

// Levels returns the available logging levels.
func (hook *RedisHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (hook *RedisHook) CloseConn() error {
	var err error

	if hook.RedisClient != nil {
		err = hook.RedisClient.Close()
	}

	if hook.RedisCluster != nil {
		err = hook.RedisCluster.Close()
	}

	if err != nil {
		return fmt.Errorf("unable to disconnect to REDIS: %s", err)
	}

	return nil
}

func newSingleConnHook(config *HookConfig) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:        config.Addrs[0],
		Password:    config.Password,
		DB:          config.DB,
		PoolSize:    3,
		IdleTimeout: 240 * time.Second,
	})

	_, err := c.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to REDIS: %s", err)
	}

	return c, nil
}

func newClusterHook(config *HookConfig) (*redis.ClusterClient, error) {
	cc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       config.Addrs,
		Password:    config.Password,
		PoolSize:    3,
		IdleTimeout: 240 * time.Second,
	})

	_, err := cc.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to REDIS Cluster: %s", err)
	}

	return cc, nil
}

func createV0Message(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["@source_host"] = hostname
	m["@message"] = entry.Message

	// build map with additional fields
	fields := make(map[string]interface{})
	fields["level"] = entry.Level.String()
	fields["application"] = appName

	for k, v := range entry.Data {
		fields[k] = v
	}

	// add fields map to message
	m["@fields"] = fields

	// return full message
	return m
}

func createV1Message(entry *logrus.Entry, appName, hostname string) map[string]interface{} {
	m := make(map[string]interface{})
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	m["host"] = hostname
	m["message"] = entry.Message
	m["level"] = entry.Level.String()
	m["application"] = appName
	for k, v := range entry.Data {
		m[k] = v
	}

	// return full message
	return m
}
