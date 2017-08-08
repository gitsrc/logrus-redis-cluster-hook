package main

import (
	"io/ioutil"

	"github.com/lazyjin/logrus-redis-cluster-hook"
	"github.com/sirupsen/logrus"
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
	// when hook is injected succesfully, logs will be send to redis server
	logrus.Debug("just some debug logging...")
	logrus.Info("just some info logging...")
	logrus.Warn("just some warn logging...")
	logrus.Error("just some error logging...")
	// logrus.Fatal("just some fatal logging...") // commented out, because it will kill the app
	// logrus.Panic("just some panic logging...") // commented out, because it will kill the app

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
