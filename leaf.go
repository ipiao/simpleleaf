package leaf

import (
	"os"
	"os/signal"

	"github.com/ipiao/simpleleaf/conf"
	"github.com/ipiao/simpleleaf/console"
	"github.com/ipiao/simpleleaf/log"
	"github.com/ipiao/simpleleaf/module"
)

func Run(mods ...module.Module) {
	// logger
	if conf.LogLevel != "" {
		logger, err := log.New(conf.LogLevel, conf.LogPath, conf.LogFlag)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	log.Infof("Leaf starting up")

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Infof("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	module.Destroy()
}
