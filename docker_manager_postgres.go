package sqldb

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/Aloe-Corporation/docker"
)

const (
	IMAGE_DOCKER_POSTGRES = "postgres"
)

var _ docker.Manager = (*DockerManagerPostgres)(nil)

// DockerManager implement docker.Manager, it's used for up and down postgres container.
type DockerManagerPostgres struct {
	ContainerName string
	PathScript    string
	Config        Conf
}

// Up is used to up postgres container.
func (dm *DockerManagerPostgres) Up() error {
	regexPassword := regexp.MustCompile(`password=([^\s]+)`)
	matchedPassword := regexPassword.MatchString(dm.Config.DSN)
	if !matchedPassword {
		return errors.New("password not found")
	}
	pwd := regexPassword.FindAllStringSubmatch(dm.Config.DSN, -1)[0][1]

	regexPort := regexp.MustCompile(`port=(\d+)`)
	matchedPort := regexPort.MatchString(dm.Config.DSN)
	if !matchedPort {
		return errors.New("port not found")
	}
	port := regexPort.FindAllStringSubmatch(dm.Config.DSN, -1)[0][1]

	err := docker.Run(IMAGE_DOCKER_POSTGRES,
		dm.ContainerName,
		"-p",
		port+":5432",
		"-v",
		dm.PathScript+":/docker-entrypoint-initdb.d/"+filepath.Base(dm.PathScript)+":ro",
		"-e",
		"POSTGRES_PASSWORD="+pwd,
		"-d",
	)
	if err != nil {
		return fmt.Errorf("can't up PostgreSQL. error : %w", err)
	}
	return nil
}

// Down is used to down postgres container.
func (dm *DockerManagerPostgres) Down() error {
	if err := docker.Down(dm.ContainerName); err != nil {
		return fmt.Errorf("fail to down container: %w", err)
	}

	return nil
}
