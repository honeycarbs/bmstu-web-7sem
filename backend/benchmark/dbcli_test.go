package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/rocketlaunchr/dbq/v2"
	"log"
	"neatly/internal/model"
	"neatly/pkg/testutils"
	"testing"
)

var (
	db                  *sql.DB
	g                   *gorm.DB
	ctx                 = context.Background()
	notEnoughRecordsErr = errors.New("select failed: not enough records")
	limits              = []int{10, 50, 150, 500, 1000}
)

func dbqQuery(lim int) {
	q := fmt.Sprintf("SELECT * FROM users ORDER BY id LIMIT %d", lim)
	res, err := dbq.Qs(ctx, db, q, model.Account{}, nil)
	if err != nil {
		panic(err)
	}
	if len(res.([]*model.Account)) != lim {
		panic(notEnoughRecordsErr)
	}
}

func sqlxQuery(lim int) {
	dbx := sqlx.NewDb(db, "postgres")
	q := fmt.Sprintf("SELECT * FROM users ORDER BY id LIMIT %d", lim)

	var res []model.Account
	err := dbx.Select(&res, q)
	if err != nil {
		panic(err)
	}
	if len(res) != lim {
		panic(notEnoughRecordsErr)
	}
}

func gormQuery(lim int) {
	var res []model.Account

	err := g.Order("id").Limit(lim).Find(&res).Error
	if err != nil {
		panic(err)
	}
	if len(res) != lim {
		panic(notEnoughRecordsErr)
	}
}

var sqlBenchmarks = []struct {
	name     string
	function func(lim int)
}{
	{"qbq", dbqQuery},
	{"sqlx", sqlxQuery},
	{"gorm", gormQuery},
}

func benchmarkSQLCli() {
	fmt.Println("\nBenchmarks for SQL clients:")
	for _, lim := range limits {
		for _, benchmark := range sqlBenchmarks {
			//for i := 0; i < 10; i++ {

			pool, err := dockertest.NewPool("")
			if err != nil {
				log.Fatalf("Could not construct pool: %s", err)
			}

			err = pool.Client.Ping()
			if err != nil {
				log.Fatalf("Could not connect to Docker: %s", err)
			}

			resource, err := pool.RunWithOptions(&dockertest.RunOptions{
				Repository: "postgres",
				Env: []string{
					"POSTGRES_DB=test",
					"POSTGRES_USER=test",
					"POSTGRES_PASSWORD=pass",
				},
			}, func(config *docker.HostConfig) {
				config.AutoRemove = true
				config.RestartPolicy = docker.RestartPolicy{
					Name: "no",
				}
			})

			if err != nil {
				log.Fatalf("Could not start resource: %s", err)
			}
			resource.Expire(60)

			dsn := "postgres://test:pass@0.0.0.0:" + resource.GetPort("5432/tcp") + "/test?sslmode=disable"

			if err := pool.Retry(func() error {
				var err error

				db, err = sql.Open("postgres", dsn)
				if err != nil {
					return err
				}
				return db.Ping()
			}); err != nil {
				log.Fatalf("Could not connect to database: %s", err)
			}
			testutils.BenchDBSetup(db)

			if benchmark.name == "gorm" {
				g, err = gorm.Open("postgres", db)
				if err != nil {
					panic(err)
				}
			}
			benchmarkFunc := func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmark.function(lim)
				}
				b.StopTimer()
			}

			results := testing.Benchmark(benchmarkFunc)
			out := fmt.Sprintf(
				"%10s -- Records pulled: %5d Time: %5d ns/op %5d allocs/op %5d bytes/op",
				benchmark.name,
				lim,
				int(results.T)/results.N,
				results.AllocsPerOp(),
				results.AllocedBytesPerOp(),
			)
			fmt.Println(out)

			if err := pool.Purge(resource); err != nil {
				log.Fatalf("Could not purge resource: %s", err)
			}
		}
	}
}

func TestDBCil(t *testing.T) {
	benchmarkSQLCli()
}
