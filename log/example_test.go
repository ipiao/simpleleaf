package log_test

import (
	l "log"

	"github.com/ipiao/simpleleaf/log"
)

func Example() {
	name := "Leaf"

	log.Debugf("My name is %v", name)
	log.Releasef("My name is %v", name)
	log.Errorf("My name is %v", name)
	// log.Fatal("My name is %v", name)

	logger, err := log.New("release", "", l.LstdFlags)
	if err != nil {
		return
	}
	defer logger.Close()

	logger.Debugf("will not print")
	logger.Releasef("My name is %v", name)

	log.Export(logger)

	log.Debugf("will not print")
	log.Releasef("My name is %v", name)
}
