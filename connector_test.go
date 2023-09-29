package sqldb

import (
	"database/sql"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Aloe-Corporation/sqldb/test"
	"github.com/stretchr/testify/assert"
)

type TableStructure struct {
	UserID    int
	Username  string
	Password  string
	Email     string
	CreatedOn time.Time
	LastLogin time.Time
}

var (
	configPostgres = Conf{
		Driver: DRIVER_POSTGRES,
		DSN:    "user=postgres password=example dbname=postgres host=localhost port=15432 sslmode=disable TimeZone=UTC",
	}
	configMysql = Conf{
		Driver: DRIVER_MYSQL,
		DSN:    "root:example@tcp(localhost:13306)/dbtest?loc=UTC&tls=false&parseTime=true",
	}

	contentTables = []TableStructure{
		{
			UserID:    1,
			Username:  "clange",
			Password:  "passwd",
			Email:     "clange@wgf.fr",
			CreatedOn: time.Date(2020, 12, 4, 10, 42, 42, 0, time.UTC),
			LastLogin: time.Date(2020, 12, 5, 10, 42, 42, 0, time.UTC),
		}, {
			UserID:    2,
			Username:  "tmazzotti",
			Password:  "passwd",
			Email:     "tmazzotti@wgf.fr",
			CreatedOn: time.Date(2020, 12, 4, 10, 42, 43, 0, time.UTC),
		}, {
			UserID:    3,
			Username:  "frichard",
			Password:  "passwd",
			Email:     "frichard@wgf.fr",
			CreatedOn: time.Date(2020, 12, 4, 10, 42, 44, 0, time.UTC),
			LastLogin: time.Date(2020, 12, 5, 10, 42, 44, 0, time.UTC),
		}, {
			UserID:    4,
			Username:  "acolin",
			Password:  "passwd",
			Email:     "acolin@wgf.fr",
			CreatedOn: time.Date(2020, 12, 4, 10, 42, 45, 0, time.UTC),
			LastLogin: time.Date(2020, 12, 5, 10, 42, 45, 0, time.UTC),
		},
	}

	insertValue = TableStructure{
		UserID:    5,
		Username:  "blemaitre",
		Password:  "passwd",
		Email:     "blemaitre@wgf.fr",
		CreatedOn: time.Date(2020, 12, 4, 10, 42, 46, 0, time.UTC),
	}

	dockerNamePostgres = "test_postgres"
	dockerNameMySQL    = "test_mysql"
	timeout            = 20
)

func TestMain(m *testing.M) {
	dmPostgres := DockerManagerPostgres{
		ContainerName: dockerNamePostgres,
		PathScript:    test.PathScriptPostgres,
		Config:        configPostgres,
	}

	dmMysql := DockerManagerMySQL{
		ContainerName: dockerNameMySQL,
		PathScript:    test.PathScriptMySQL,
		Config:        configMysql,
	}

	if err := dmPostgres.Up(); err != nil {
		panic(err)
	}

	if err := dmMysql.Up(); err != nil {
		_ = dmPostgres.Down()
		panic(err)
	}
	time.Sleep(2 * time.Second)

	r := m.Run()

	if err := dmPostgres.Down(); err != nil {
		_ = dmMysql.Down()
		panic(err)
	}

	if err := dmMysql.Down(); err != nil {
		panic(err)
	}

	os.Exit(r)
}

type factoryConnectorTestData struct {
	Config              Conf
	ShouldFail          bool
	LastInsertIdAvaible bool
}

var factoryConnectorTestCases = []factoryConnectorTestData{
	{
		Config:              configPostgres,
		ShouldFail:          false,
		LastInsertIdAvaible: false,
	},
	{
		Config:              configMysql,
		ShouldFail:          false,
		LastInsertIdAvaible: true,
	},
	{
		Config: Conf{
			Driver: "wrong",
			DSN:    "wrong",
		},
		ShouldFail: true,
	},
}

func TestFactoryConnector(t *testing.T) {
	for i, testCase := range factoryConnectorTestCases {
		t.Run("Case "+strconv.Itoa(i), func(t *testing.T) {
			connector, err := FactoryConnector(testCase.Config)
			if testCase.ShouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, connector)
			}
		})
	}
}

type connectorTryConnectionTestData struct {
	Config     Conf
	ShouldFail bool
}

var connectorTryConnectionTestCases = []connectorTryConnectionTestData{
	{
		Config:     configPostgres,
		ShouldFail: false,
	},
	{
		Config:     configMysql,
		ShouldFail: false,
	},
	{
		Config: Conf{
			Driver: configPostgres.Driver,
			DSN:    "user=bad_user password=example dbname=postgres host=localhost port=15432 sslmode=disable TimeZone=UTC",
		},
		ShouldFail: true,
	},
}

func TestConnectorTryConnection(t *testing.T) {
	for i, testCase := range connectorTryConnectionTestCases {
		t.Run("Case "+strconv.Itoa(i), func(t *testing.T) {
			connector, err := FactoryConnector(testCase.Config)
			assert.NoError(t, err)
			assert.NotNil(t, connector)

			err = connector.TryConnection(timeout)
			if testCase.ShouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type connectorCommitTestData struct {
	Config Conf
}

var connectorCommitTestCases = [...]connectorCommitTestData{
	{Config: configPostgres},
	{Config: configMysql},
}

func TestTransactionControl(t *testing.T) {
	for i, testCase := range connectorCommitTestCases {
		t.Run("Case "+strconv.Itoa(i), func(t *testing.T) {
			db, err := FactoryConnector(testCase.Config)
			assert.NoError(t, err)
			assert.NotNil(t, db)

			err = db.TryConnection(timeout)
			assert.NoError(t, err)

			tx, err := db.Begin()
			assert.NoError(t, err)
			assert.NotNil(t, tx)

			err = db.Commit(tx)
			assert.NoError(t, err)
		})
	}
}

func TestExecQueryRow(t *testing.T) {
	db, err := FactoryConnector(configPostgres)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.TryConnection(timeout)
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)

	row, err := db.ExecQueryRow(
		tx,
		"INSERT INTO table_insert_test(user_id, username, password, email, created_on, last_login) VALUES ($1, $2, $3, $4, $5, NULL) RETURNING user_id",
		insertValue.UserID, insertValue.Username, insertValue.Password, insertValue.Email, insertValue.CreatedOn,
	)
	if !assert.NoError(t, err) || !assert.NotNil(t, row) {
		t.FailNow()
	}

	var id int
	err = row.Scan(&id)
	assert.NoError(t, err)
	assert.Equal(t, 5, id)

	err = db.Commit(tx)
	assert.NoError(t, err)

	row = db.QueryRow("SELECT * FROM table_insert_test WHERE user_id=$1", id)
	rowResult := TableStructure{}
	var lastLoginNullTime sql.NullTime

	err = row.Scan(&rowResult.UserID, &rowResult.Username, &rowResult.Password, &rowResult.Email, &rowResult.CreatedOn, &lastLoginNullTime)
	assert.NoError(t, err)
	if lastLoginNullTime.Valid {
		rowResult.LastLogin = lastLoginNullTime.Time
	}

	assert.Equal(t, rowResult, insertValue)
}

func TestExecQueryRowFail(t *testing.T) {
	db, err := FactoryConnector(configPostgres)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.TryConnection(timeout)
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)

	_, err = db.ExecQueryRow(tx, "NOT SQL QUERY, test for r******d proof")
	assert.Error(t, err)
}

type connectorExecTestData struct {
	Config Conf
}

var connectorExecTestCases = [...]connectorExecTestData{
	{Config: configPostgres},
	{Config: configMysql},
}

func TestConnectorExec(t *testing.T) {
	for i, testCase := range connectorExecTestCases {
		t.Run("Case "+strconv.Itoa(i), func(t *testing.T) {
			db, err := FactoryConnector(testCase.Config)
			assert.NoError(t, err)
			assert.NotNil(t, db)

			err = db.TryConnection(timeout)
			assert.NoError(t, err)

			tx, err := db.Begin()
			assert.NoError(t, err)

			testRow := contentTables[0]
			testRow.Username = "vracine"

			queryUpdate := "UPDATE table_update_test SET username=? WHERE user_id=?"
			if testCase.Config.Driver == DRIVER_POSTGRES {
				queryUpdate = "UPDATE table_update_test SET username=$1 WHERE user_id=$2"
			}
			res, err := db.Exec(
				tx,
				queryUpdate,
				testRow.Username, testRow.UserID,
			)
			assert.NoError(t, err)
			rowAffected, err := res.RowsAffected()
			assert.NoError(t, err)
			assert.EqualValues(t, 1, rowAffected)

			err = db.Commit(tx)
			assert.NoError(t, err)

			querySelect := "SELECT * FROM table_update_test WHERE user_id=?"
			if testCase.Config.Driver == DRIVER_POSTGRES {
				querySelect = "SELECT * FROM table_update_test WHERE user_id=$1"
			}
			row := db.QueryRow(querySelect, testRow.UserID)
			rowResult := TableStructure{}
			var lastLoginNullTime sql.NullTime

			err = row.Scan(&rowResult.UserID, &rowResult.Username, &rowResult.Password, &rowResult.Email, &rowResult.CreatedOn, &lastLoginNullTime)
			assert.NoError(t, err)
			if lastLoginNullTime.Valid {
				rowResult.LastLogin = lastLoginNullTime.Time
			}

			assert.Equal(t, testRow, rowResult)
		})
	}
}

type connectorExecFailTestData struct {
	Config Conf
}

var connectorExecFailTestCases = [...]connectorExecFailTestData{
	{Config: configPostgres},
	{Config: configMysql},
}

func TestConnectorExecFail(t *testing.T) {
	for i, testCase := range connectorExecFailTestCases {
		t.Run("Case "+strconv.Itoa(i), func(t *testing.T) {
			db, err := FactoryConnector(testCase.Config)
			assert.NoError(t, err)
			assert.NotNil(t, db)

			err = db.TryConnection(timeout)
			assert.NoError(t, err)

			tx, err := db.Begin()
			assert.NoError(t, err)

			_, err = db.Exec(tx, "NOT SQL QUERY, test for r******d proof")
			assert.Error(t, err)
		})
	}
}
