package sqldb

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Aloe-Corporation/docker"
)

const (
	IMAGE_DOCKER_MYSQL = "mysql"
)

var _ docker.Manager = (*DockerManagerMySQL)(nil)

// DockerManager implements docker.Manager, it's used to up and down postgres container.
type DockerManagerMySQL struct {
	ContainerName string
	PathScript    string
	Config        Conf
}

// Up is used for up postgres container.
func (dm *DockerManagerMySQL) Up() error {
	regexPassword := regexp.MustCompile(`:(.*)@`)
	matchedPassword := regexPassword.MatchString(dm.Config.DSN)
	if !matchedPassword {
		return errors.New("password not found")
	}
	pwd := regexPassword.FindAllStringSubmatch(dm.Config.DSN, -1)[0][1]

	regexPort := regexp.MustCompile(`\(.*:(\d+)\)`)
	matchedPort := regexPort.MatchString(dm.Config.DSN)
	if !matchedPort {
		return errors.New("port not found")
	}
	port := regexPort.FindAllStringSubmatch(dm.Config.DSN, -1)[0][1]

	err := docker.Run(IMAGE_DOCKER_MYSQL,
		dm.ContainerName,
		"-p",
		port+":3306",
		"-v",
		dm.PathScript+":/docker-entrypoint-initdb.d/schema.sql:ro",
		"-e",
		"MYSQL_ROOT_PASSWORD="+pwd,
		"-d",
	)
	if err != nil {
		return fmt.Errorf("can't up MySQL. error : %w", err)
	}
	return nil
}

// Down is used for down postgres container.
func (dm *DockerManagerMySQL) Down() error {
	if err := docker.Down(dm.ContainerName); err != nil {
		return fmt.Errorf("fail to down container: %w", err)
	}

	return nil
}
