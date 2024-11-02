package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/natecw/minily/api"
	"github.com/natecw/minily/cache"
	"github.com/natecw/minily/storage"
	"github.com/urfave/cli"
)

const (
	apiAddrFlagname       string = "addr"
	apiStorageDatabaseUrl string = "database-url"
)

func main() {
	if err := app().Run(os.Args); err != nil {
		log.Fatalf("could not run %v\n", err)
	}
}

func app() *cli.App {
	return &cli.App{
		Name:  "api-server",
		Usage: "Api",
		Commands: []cli.Command{
			serverCmd(),
		},
	}
}

func serverCmd() cli.Command {
	return cli.Command{
		Name:  "start",
		Usage: "starts the API server",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: apiAddrFlagname, EnvVar: "API_SERVER_ADDR"},
			&cli.StringFlag{Name: apiStorageDatabaseUrl, EnvVar: "DATABSE_URL"},
		},
		Action: func(c *cli.Context) error {
			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			stopper := make(chan struct{})
			go func() {
				<-done
				close(stopper)
			}()

			port, err := strconv.ParseInt(os.Getenv("REDIS_PORT"), 10, 64)
			if err != nil {
				return err
			}

			cache := cache.NewCache(os.Getenv("REDIS_HOST"), port)

			databaseUrl := c.String(apiStorageDatabaseUrl)
			s, err := storage.NewStorage(databaseUrl, cache)
			if err != nil {
				return fmt.Errorf("could not initialize storage: %w", err)
			}

			addr := c.String(apiAddrFlagname)
			server, err := api.NewApi(addr, s, slog.Default())
			if err != nil {
				return err
			}

			return server.Start(stopper)
		},
	}
}
