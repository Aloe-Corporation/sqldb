package sqldb

import (
	"strconv"
	"testing"

	"github.com/Aloe-Corporation/sqldb/test"
	"github.com/stretchr/testify/assert"
)

type dockerManagerMySQLUpAndDownTestData struct {
	ContainerName string
	PathScript    string
	Config        Conf
	ShouldFail    bool
}

var dockerManagerMySQLUpAndDownTestCases = []dockerManagerMySQLUpAndDownTestData{
	{
		ContainerName: "test_mysql_docker_manager",
		PathScript:    test.PathScriptMySQL,
		Config: Conf{
			Driver: DRIVER_MYSQL,
			DSN:    "root:example@tcp(localhost:3306)/dbtest?loc=UTC&tls=false&parseTime=true",
		},
		ShouldFail: false,
	},
	{
		ContainerName: "test_mysql_docker_manager_fail",
		PathScript:    test.PathScriptMySQL,
		Config: Conf{
			Driver: DRIVER_MYSQL,
			DSN:    "root:example@tcp(localhost:-3306)/dbtest?loc=UTC&tls=false&parseTime=true",
		},
		ShouldFail: true,
	},
}

func TestDockerManagerMySQLUpAndDown(t *testing.T) {
	for i, testCase := range dockerManagerMySQLUpAndDownTestCases {
		t.Run("Case "+strconv.Itoa(i), func(t *testing.T) {
			dm := DockerManagerMySQL{
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
