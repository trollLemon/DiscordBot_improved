package main

import (
	"context"
	"flag"

	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"goManip/jobdispatch"
	"goManip/jobs"
	"goManip/server"
	"goManip/worker"
)

func main() {
	prettyPrint := flag.Bool("pretty_print", false, "Enable pretty print output instead of structured json")
	numWorkers := flag.Int("num_workers", runtime.NumCPU(), "Number of workers")
	port := flag.String("port", "8080", "port to listen on")
	address := flag.String("address", "localhost", "address to bind to")
	flag.Parse()

	if *prettyPrint {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	jobReqs := make(chan *jobs.JobRequest, *numWorkers)
	maxTime := time.Second * 10
	jobDispatcher := jobdispatch.NewJobDispatcher(jobReqs, maxTime)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		server.GraceFullShutdown(jobDispatcher, wg, cancel)
		os.Exit(0)
	}()

	for workerId := range *numWorkers {
		log.Info().Msgf("Starting worker #%d", workerId+1)
		wg.Add(1)
		go worker.Worker(ctx, workerId+1, jobReqs, wg)
	}
	e := echo.New()
	server.InitRouting(e, jobDispatcher)
	server.Start(e, *address, *port)
}
