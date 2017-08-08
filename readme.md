# Redis Cluster Hook for [Logrus](https://github.com/Sirupsen/logrus)
[![Build Status](https://travis-ci.org/lazyjin/logrus-redis-cluster-hook.svg?branch=master)](https://travis-ci.org/lazyjin/logrus-redis-cluster-hook)

logrus-redis-cluster-hook is a redis hook for [Logrus](https://github.com/Sirupsen/logrus), based on [logrus-redis-cluster]( https://github.com/rogierlommers/logrus-redis-hook). 

[go-redis](https://github.com/go-redis/redis) is used for redis connection, And both single connection and cluster connection are supported.



## Install

```shell
$ go get github.com/lazyjin/logrus-redis-cluster-hook
```



## Usage

```go
package main

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/lazyjin/logrus-redis-cluster-hook"
)

func init() {
	hookConfig := logredis.HookConfig{
		Addrs:      []string{"localhost:6379"},
		ConnOption: logredis.SINGLE,
		// Addrs:      []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"},
		// ConnOption: logredis.CLUSTER,
		Key:      "my_redis_key",
		Format:   "v0",
		App:      "my_app_name",
		Hostname: "my_app_hostmame",
		DB:       0,
	}

	hook, err := logredis.NewHook(hookConfig)
	if err == nil {
		logrus.AddHook(hook)
	} else {
		logrus.Errorf("logredis error: %q", err)
	}
}

func main() {
	// when hook is injected succesfully, logs will be sent to redis server
	logrus.Info("just some info logging...")

	// we also support log.WithFields()
	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"foo":    "bar",
		"this":   "that"}).
		Info("additional fields are being logged as well")

	// If you want to disable writing to stdout, use setOutput
	logrus.SetOutput(ioutil.Discard)
	logrus.Info("This will only be sent to Redis")
}
```
