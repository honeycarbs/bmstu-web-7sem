package testdocker

import (
	"fmt"
	"github.com/ory/dockertest"
	"neatly/pkg/dbclient"
)

func GetTestResource(migrationPath string) (*dbclient.Client, *dockertest.Resource, *dockertest.Pool, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker pool: %s", err)
	}

	resource, err := pool.Run("postgres", "alpine",
		[]string{"POSTGRES_DB=test",
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=pass"})

	//resource.Expire(100)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not start resource: %s", err)
	}

	var cli *dbclient.Client

	if err := pool.Retry(func() error {
		cli, err = dbclient.NewIntegrationClinent(resource, migrationPath)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker resource: %s", err)
	}

	return cli, resource, pool, nil
}