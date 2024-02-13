package main

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/scorpiotzh/toolib"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
	"time"
	"twitter-to-lark/cache"
	"twitter-to-lark/config"
	"twitter-to-lark/task"
)

var (
	log               = logger.NewLogger("main", logger.LevelDebug)
	exit              = make(chan struct{})
	ctxServer, cancel = context.WithCancel(context.Background())
	wgServer          = sync.WaitGroup{}
)

func main() {
	log.Debugf("start：")
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: runServer,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx *cli.Context) error {
	// config file
	configFilePath := ctx.String("config")
	if err := config.InitCfg(configFilePath); err != nil {
		return err
	}

	// config file watcher
	watcher, err := config.AddCfgFileWatcher(configFilePath)
	if err != nil {
		return fmt.Errorf("AddCfgFileWatcher err: %s", err.Error())
	}
	// ============= service start =============

	// redis
	red, err := toolib.NewRedisClient(config.Cfg.Cache.Redis.Addr, config.Cfg.Cache.Redis.Password, config.Cfg.Cache.Redis.DbNum)
	if err != nil {
		log.Info("NewRedisClient err: %s", err.Error())
		//return fmt.Errorf("NewRedisClient err:%s", err.Error())
	} else {
		log.Info("redis ok")
	}
	rc := cache.Initialize(red)
	if err := rc.SetTweets2lark("test11"); err != nil {
		fmt.Println("set redis err11111")
	} else {
		fmt.Println("set redis ok22222")
	}
	// task
	twitterTask := task.Task{
		Ctx: ctxServer,
		Wg:  &wgServer,
		Rc:  rc,
	}
	twitterTask.Run()
	log.Infof("monitor_task ok")

	// ============= service end =============
	toolib.ExitMonitoring(func(sig os.Signal) {
		log.Warn("ExitMonitoring:", sig.String())
		if watcher != nil {
			log.Warn("close watcher ... ")
			_ = watcher.Close()
		}
		cancel()
		wgServer.Wait()
		log.Warn("success exit server. bye bye!")
		time.Sleep(time.Second)
		exit <- struct{}{}
	})

	<-exit

	return nil
}
