package main

import (
	. "github.com/simon987/imhashdb"
	"github.com/simon987/imhashdb/hasher"
	api "github.com/simon987/imhashdb/web"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "web",
				Usage: "Start http API",
				Action: func(c *cli.Context) error {
					Init()
					return api.Main()
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "pg-user",
						Value:       "imhashdb",
						Usage:       "PostgreSQL user",
						EnvVars:     []string{"IMHASHDB_PG_USER"},
						Destination: &Conf.PgUser,
					},
					&cli.StringFlag{
						Name:        "pg-password",
						Value:       "imhashdb",
						Usage:       "PostgreSQL password",
						EnvVars:     []string{"IMHASHDB_PG_PASSWORD"},
						Destination: &Conf.PgPassword,
					},
					&cli.StringFlag{
						Name:        "pg-db",
						Value:       "imhashdb",
						Usage:       "PostgreSQL database",
						EnvVars:     []string{"IMHASHDB_PG_DATABASE"},
						Destination: &Conf.PgDb,
					},
					&cli.StringFlag{
						Name:        "pg-host",
						Value:       "localhost",
						Usage:       "PostgreSQL host",
						EnvVars:     []string{"IMHASHDB_PG_HOST"},
						Destination: &Conf.PgHost,
					},
					&cli.IntFlag{
						Name:        "pg-port",
						Value:       5432,
						Usage:       "PostgreSQL port",
						EnvVars:     []string{"IMHASHDB_PG_PORT"},
						Destination: &Conf.PgPort,
					},
					&cli.StringFlag{
						Name:        "redis-addr",
						Value:       "localhost:6379",
						Usage:       "redis address",
						EnvVars:     []string{"IMHASHDB_REDIS_ADDR"},
						Destination: &Conf.RedisAddr,
					},
					&cli.StringFlag{
						Name:        "redis-password",
						Value:       "",
						Usage:       "redis password",
						EnvVars:     []string{"IMHASHDB_REDIS_PASSWORD"},
						Destination: &Conf.RedisPassword,
					},
					&cli.IntFlag{
						Name:        "redis-db",
						Value:       0,
						Usage:       "redis db",
						EnvVars:     []string{"IMHASHDB_REDIS_DB"},
						Destination: &Conf.RedisDb,
					},
					&cli.StringFlag{
						Name:        "api-addr",
						Value:       "localhost:8080",
						Usage:       "HTTP api address",
						EnvVars:     []string{"IMHASHDB_API_ADDR"},
						Destination: &Conf.ApiAddr,
					},
					&cli.IntFlag{
						Name:        "query-concurrency",
						Value:       2,
						Usage:       "Number of background query workers",
						EnvVars:     []string{"IMHASHDB_QUERY_CONCURRENCY"},
						Destination: &Conf.QueryConcurrency,
					},
				},
			},
			{
				Name:  "hasher",
				Usage: "Start an hasher instance",
				Action: func(c *cli.Context) error {
					Init()
					return hasher.Main()
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "pg-user",
						Value:       "imhashdb",
						Usage:       "PostgreSQL user",
						EnvVars:     []string{"IMHASHDB_PG_USER"},
						Destination: &Conf.PgUser,
					},
					&cli.StringFlag{
						Name:        "pg-password",
						Value:       "imhashdb",
						Usage:       "PostgreSQL password",
						EnvVars:     []string{"IMHASHDB_PG_PASSWORD"},
						Destination: &Conf.PgPassword,
					},
					&cli.StringFlag{
						Name:        "pg-db",
						Value:       "imhashdb",
						Usage:       "PostgreSQL database",
						EnvVars:     []string{"IMHASHDB_PG_DATABASE"},
						Destination: &Conf.PgDb,
					},
					&cli.StringFlag{
						Name:        "pg-host",
						Value:       "localhost",
						Usage:       "PostgreSQL host",
						EnvVars:     []string{"IMHASHDB_PG_HOST"},
						Destination: &Conf.PgHost,
					},
					&cli.IntFlag{
						Name:        "pg-port",
						Value:       5432,
						Usage:       "PostgreSQL port",
						EnvVars:     []string{"IMHASHDB_PG_PORT"},
						Destination: &Conf.PgPort,
					},
					&cli.StringFlag{
						Name:        "redis-addr",
						Value:       "localhost:6379",
						Usage:       "redis address",
						EnvVars:     []string{"IMHASHDB_REDIS_ADDR"},
						Destination: &Conf.RedisAddr,
					},
					&cli.StringFlag{
						Name:        "redis-password",
						Value:       "",
						Usage:       "redis password",
						EnvVars:     []string{"IMHASHDB_REDIS_PASSWORD"},
						Destination: &Conf.RedisPassword,
					},
					&cli.IntFlag{
						Name:        "redis-db",
						Value:       0,
						Usage:       "redis db",
						EnvVars:     []string{"IMHASHDB_REDIS_DB"},
						Destination: &Conf.RedisDb,
					},
					&cli.StringFlag{
						Name:        "imgur-clientid",
						Value:       "546c25a59c58ad7",
						Usage:       "imgur API client id",
						EnvVars:     []string{"IMHASHDB_IMGUR_CLIENTID"},
						Destination: &Conf.ImgurClientId,
					},
					&cli.StringFlag{
						Name:        "hasher-pattern",
						Value:       "q.*",
						Usage:       "redis pattern for hasher input tasks",
						EnvVars:     []string{"IMHASHDB_HASHER_PATTERN"},
						Destination: &Conf.HasherPattern,
					},
					&cli.IntFlag{
						Name:        "hash-concurrency",
						Value:       4,
						Usage:       "Thread count for hasher",
						EnvVars:     []string{"IMHASHDB_HASH_CONCURRENCY"},
						Destination: &Conf.HasherConcurrency,
					},
					&cli.StringFlag{
						Name:        "store",
						Value:       "",
						Usage:       "If set, store downloaded images there",
						EnvVars:     []string{"IMHASHDB_STORE"},
						Destination: &Conf.Store,
					},
				},
			},
		},
		Authors: []*cli.Author{
			{
				Name:  "simon987",
				Email: "me@simon987.net",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
