package main

import (
	"flag"
	"github.com/millken/logger"
	"runtime"
	"sync"
)

var (
	VERSION    string = "0.1"
	config     tomlConfig
	configPath string
	gitVersion string
	workerCh   = make(chan string)
)

func init() {
	if len(gitVersion) > 0 {
		VERSION = VERSION + "/" + gitVersion
	}
}

func main() {
	var err error
	flag.StringVar(&configPath, "config", "ngx2es.toml", "config path")
	flag.Parse()
	logger.Info("Loading config : %s, version: %s", configPath, VERSION)
	err = LoadConfig(configPath)
	if err != nil {
		logger.Exitf("Read config failed.Err = %s", err.Error())
	}

	numCpus := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpus)

	wg := new(sync.WaitGroup)
	for i := 0; i < numCpus; i++ {
		wg.Add(1)
		go startWorker(i, workerCh, wg)
	}

	startKafkaService()

	close(workerCh)
	wg.Wait()

}
