package sqldb

import (
	"strconv"
	"testing"

	"github.com/Aloe-Corporation/sqldb/test"
	"github.com/stretchr/testify/assert"
)

type dockerManagerPostgresUpAndDownTestData struct {
	ContainerName string
	PathScript    string
	Config        Conf
	ShouldFail    bool
}

var dockerManagerPostgresUpAndDownTestCases = []dockerManagerPostgresUpAndDownTestData{
	{
		ContainerName: "test_postgres_docker_manager",
		PathScript:    test.PathScriptPostgres,
		Config: Conf{
			Driver: DRIVER_POSTGRES,
			DSN:    "user=postgres password=example dbname=postgres host=localhost port=5432 sslmode=disable TimeZone=UTC",
		},
		ShouldFail: false,
	},
	{
		ContainerName: "test_postgres_docker_manager_fail",
		PathScript:    test.PathScriptPostgres,
		Config: Conf{
			Driver: DRIVER_POSTGRES,
			DSN:    "user=postgres password=example dbname=postgres host=localhost port=-5432 sslmode=disable TimeZone=UTC",
		},
		ShouldFail: true,
	},
}

func TestDockerManagerPostgresUpAndDown(t *testing.T) {
	for i, testCase := range dockerManagerPostgresUpAndDownTestCases {
		t.Run("Case "+strconv.Itoa(i), func(t *testing.T) {
			dm := DockerManagerPostgres{
				ContainerName: testCase.ContainerName,
				PathScript:    testCase.PathScript,
				Config:        testCase.Config,
			}

			if testCase.ShouldFail {
				err := dm.Up()
				assert.Error(t, err)
				err = dm.Down()
				assert.Error(t, err)
			} else {
				err := dm.Up()
				assert.NoError(t, err)
				err = dm.Down()
				assert.NoError(t, err)
			}
		})
	}
}
