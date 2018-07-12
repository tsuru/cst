package scan

import (
	"time"

	"github.com/optiopay/klar/clair"
	"github.com/optiopay/klar/docker"
	"github.com/sirupsen/logrus"
)

// Clair is a struct that implements a Scanner interface.
type Clair struct {
	Address string
	Name    string
	Timeout time.Duration
}

// Scan analyzes a container image on CoreOS Clair security engine.
func (c *Clair) Scan(image string) Result {

	log := logrus.
		WithField("clair.address", c.Address).
		WithField("image", image)

	log.Info("initializing scan on CoreOS Clair")

	defer log.Info("finishing scan on CoreOS Clair")

	dockerImage, err := docker.NewImage(&docker.Config{
		ImageName: image,
	})

	if err != nil {
		return c.makeErrorResult(err)
	}

	log.
		WithField("docker.registry", dockerImage.Registry).
		WithField("docker.image", dockerImage.Name).
		WithField("docker.tag", dockerImage.Tag).
		Info("fetching manifest and layers from image's repository")

	err = dockerImage.Pull()

	if err != nil {
		return c.makeErrorResult(err)
	}

	var vulns []*clair.Vulnerability

	for _, apiVersion := range []int{1, 3} {
		clairClient := clair.NewClair(c.Address, apiVersion, c.Timeout)

		vulns, err = clairClient.Analyse(dockerImage)

		if err == nil {
			break
		}

		log.
			WithField("clair.api", apiVersion).
			WithError(err).
			Warn("failed to analyze using that CoreOS Clair API version")

		continue
	}

	if err != nil {
		return c.makeErrorResult(err)
	}

	vulnerabilities := make([]clair.Vulnerability, len(vulns))

	for index, vulnerability := range vulns {
		vulnerabilities[index] = *vulnerability
	}

	log.Info("successful to get vulnerabilities on CoreOS Clair")

	return Result{
		Scanner:         c.Name,
		Vulnerabilities: vulnerabilities,
	}
}

func (c *Clair) makeErrorResult(err error) Result {

	return Result{
		Scanner: c.Name,
		Error:   "could not analyze that image on CoreOS Clair",
	}
}
