package sqldb

import (
	"database/sql/driver"
	"errors"
	"os"
	"strconv"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	configPostgres = Conf{
		Driver: DRIVER_POSTGRES,
		DSN:    "user=postgres password=example dbname=postgres host=localhost port=15432 sslmode=disable TimeZone=UTC",
	}
	configMysql = Conf{
		Driver: DRIVER_MYSQL,
		DSN:    "root:example@tcp(localhost:13306)/dbtest?loc=UTC&tls=false&parseTime=true",
	}

	timeout = 1
)

func TestMain(m *testing.M) {
	r := m.Run()
	os.Exit(r)
}

type factoryConnectorTestData struct {
	name       string
	config     Conf
	shouldFail bool
}

var factoryConnectorTestCases = []factoryConnectorTestData{
	{
		name:       "Valid postgres connector",
		config:     configPostgres,
		shouldFail: false,
	},
	{
		name:       "Valid mysql connector",
		config:     configMysql,
		shouldFail: false,
	},
	{
		name: "Invalid driver",
		config: Conf{
			Driver: "dummy",
			DSN:    configMysql.DSN,
		},
		shouldFail: true,
	},
	{
		name: "Invalid DSN",
		config: Conf{
			Driver: "mysql",
			DSN:    "\t\t",
		},
		shouldFail: true,
	},
}

func TestFactoryConnector(t *testing.T) {
	for _, testCase := range factoryConnectorTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			connector, err := FactoryConnector(testCase.config)
			if testCase.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, connector)
			}
		})
	}
}

type connectorTryConnectionTestData struct {
	name       string
	config     Conf
	shouldFail bool
	withMock   bool
}

var connectorTryConnectionTestCases = []connectorTryConnectionTestData{
	{
		name:       "Successs with postgres",
		config:     configPostgres,
		shouldFail: false,
		withMock:   true,
	},
	{
		name:       "Successs with MySQL",
		config:     configMysql,
		shouldFail: false,
		withMock:   true,
	},
	{
		name:       "Fail case, no database to reach",
		config:     configMysql,
		shouldFail: true,
		withMock:   false,
	},
}

func TestConnectorTryConnection(t *testing.T) {
	for _, testCase := range connectorTryConnectionTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			connector, err := FactoryConnector(testCase.config)
			assert.NoError(t, err)

			if testCase.withMock {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				defer db.Close()

				mock.ExpectPing()

				connector.DB = db
				assert.NoError(t, err)
				assert.NotNil(t, connector)
			}

			err = connector.TryConnection(timeout)
			if testCase.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type connectorCommitTestData struct {
	name          string
	config        Conf
	withInvalidTx bool
	shouldFail    bool
}

var connectorCommitTestCases = [...]connectorCommitTestData{
	{
		name:          "Successs case with postgres",
		config:        configPostgres,
		withInvalidTx: false,
		shouldFail:    false,
	},
	{
		name:          "Successs case with MySQL",
		config:        configMysql,
		withInvalidTx: false,
		shouldFail:    false,
	},
	{
		name:          "Fail case with postgres and invalid Tx",
		config:        configPostgres,
		withInvalidTx: true,
		shouldFail:    true,
	},
}

func TestTransactionControl(t *testing.T) {
	for _, testCase := range connectorCommitTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectCommit()

			connector, err := FactoryConnector(testCase.config)
			assert.NoError(t, err)
			connector.DB = db

			tx, err := connector.Begin()
			assert.NoError(t, err)
			assert.NotNil(t, tx)

			if testCase.withInvalidTx {
				tx.Rollback()
			}

			err = connector.Commit(tx)
			if testCase.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type connectorExecQueryRowTestData struct {
	name       string
	config     Conf
	query      string
	shouldFail bool
}

var connectorExecQueryRowTestCases = [...]connectorExecQueryRowTestData{
	{
		name:       "Success case with postgres",
		config:     configPostgres,
		query:      "INSERT INTO table_insert_test",
		shouldFail: false,
	},
	{
		name:       "Success case with MySQL",
		config:     configMysql,
		query:      "INSERT INTO table_insert_test",
		shouldFail: false,
	},
	{
		name:       "Fail to prepare statement error",
		config:     configMysql,
		query:      "******",
		shouldFail: true,
	},
}

func TestExecQueryRow(t *testing.T) {
	for _, testCase := range connectorExecQueryRowTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectPrepare(testCase.query)
			mock.ExpectCommit()

			connector, err := FactoryConnector(configPostgres)
			assert.NoError(t, err)
			assert.NotNil(t, connector)
			connector.DB = db

			tx, err := db.Begin()
			assert.NoError(t, err)

			_, err = connector.ExecQueryRow(
				tx,
				testCase.query,
			)

			if testCase.shouldFail {
				assert.Error(t, err)
			} else {
				err = connector.Commit(tx)
				assert.NoError(t, err)
			}

		})
	}
}

type connectorExecTestData struct {
	name       string
	config     Conf
	query      string
	shouldFail bool
}

var connectorExecTestCases = [...]connectorExecTestData{
	{
		name:       "Successs case with postgres",
		config:     configPostgres,
		query:      "INSERT INTO table_insert_test",
		shouldFail: false,
	},
	{
		name:       "Successs case with MySQL",
		config:     configMysql,
		query:      "INSERT INTO table_insert_test",
		shouldFail: false,
	},
	{
		name:       "Fail case invalid sql query",
		config:     configMysql,
		query:      "*******",
		shouldFail: true, // indicates to mocking to return an error
	},
	{
		name:       "Fail case with MySQL, exec error",
		config:     configMysql,
		query:      "INSERT INTO table_insert_test",
		shouldFail: true, // indicates to mocking to return an error
	},
}

func TestConnectorExec(t *testing.T) {
	for i, testCase := range connectorExecTestCases {
		t.Run("Case "+strconv.Itoa(i), func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectPrepare(testCase.query)
			if testCase.shouldFail {
				mock.ExpectExec(testCase.query).WithArgs(1).WillReturnError(errors.New("error"))
			} else {
				mock.ExpectExec(testCase.query).WithArgs(1).WillReturnResult(driver.ResultNoRows)
			}
			mock.ExpectCommit()

			connector, err := FactoryConnector(testCase.config)
			assert.NoError(t, err)
			assert.NotNil(t, db)
			connector.DB = db

			tx, err := connector.Begin()
			assert.NoError(t, err)

			_, err = connector.Exec(
				tx,
				"INSERT INTO table_insert_test",
				1,
			)

			if testCase.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				err = connector.Commit(tx)
				assert.NoError(t, err)
			}
		})
	}
}
