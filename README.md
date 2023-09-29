# sqldb

This project is a module for SqlDB Connector.

## Usage

### Configuration

The `sqldb.Conf` uses a YAML tags, it's easy to load SqlDB config with configuration file in your project

```go
type Conf struct {
	Driver string `yaml:"driver"` // example: postgres, mysql
	DSN    string `yaml:"dsn"` // connection string (format depends on the driver, read the associated documentation)
}
```

Example DSN:

- `postgres:` user=postgres password=example dbname=postgres host=localhost port=5432 sslmode=disable TimeZone=UTC
- `mysql:` root:example@tcp(localhost:3306)/dbtest?loc=UTC&tls=false&parseTime=true `(WARNING: parseTime=true is require)`

### Environement variable

- PREFIX\_ + `SQLDB_DRIVER`
- PREFIX\_ + `SQLDB_DSN`

### Create new connector

To create new SqlDB Connector use this function with as configuration the structure `sqldb.FactoryConnector(c sqldb.Conf) (*sqldb.Connector, error)` and try connection with `sqldb.Connector.TryConnection(t int) err`

```go
var config = sqldb.Conf{
	Driver:     "mysql",
	Pwd:      	"root:example@tcp(localhost:13306)/dbtest?loc=UTC&tls=false&parseTime=trueword",
}

// Build Connector
pg, err = sqldb.FactoryConnector(config)
if err != nil {
	return fmt.Errorf("fail to init SqlDB connector: %w", err)
}

// Test connection
err = pg.TryConnection(10)
if err != nil {
	return fmt.Errorf("fail to ping SqlDB: %w", err)
}

```

### Use DockerManager for testing

`sqldb.DockerManager` is used to up and down SqlDB container during the test with the `sqldb.Conf`

```go
var config = sqldb.Conf{
	Driver:     "mysql",
	Pwd:      	"root:example@tcp(localhost:13306)/dbtest?loc=UTC&tls=false&parseTime=trueword",
}

dm := sqldb.DockerManager{
	ContainerName: "docker_manager_postgres",
	PathScript:    "path_script.sql",
	Config:        config,
}

// Up docker container
err := dm.Up()
if err != nil {
	panic(fmt.Errorf("fail to Up SqlDB Docker container: %w", err))
}

// Down docker container
err = dm.Down()
if err != nil {
	panic(fmt.Errorf("fail to Down SqlDB Docker container: %w", err))
}
```

## Test

- `make test`

`test` folder contains script to set up test database.
