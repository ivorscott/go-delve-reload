package secrets

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// DockerSecrets contains secrets
type DockerSecrets struct {
	secretsDir string
	secrets    map[string]string
}

// NewDockerSecrets creates an instance of DockerSecrets
func NewDockerSecrets() (*DockerSecrets, error) {
	secretsDir := "/run/secrets"
	dockerSecrets := &DockerSecrets{secretsDir: secretsDir, secrets: map[string]string{}}
	err := dockerSecrets.readAll()
	return dockerSecrets, err
}

// GetDir returns the secretsDir
func (ds *DockerSecrets) GetDir() string {
	return ds.secretsDir
}

// Get returns one secret by secretName
func (ds *DockerSecrets) Get(secretName string) (string, error) {
	if _, ok := ds.secrets[secretName]; !ok {
		return "", fmt.Errorf("secret not exsist: %s", secretName)
	}
	return ds.secrets[secretName], nil
}

// Reads all secrets on the specified path in the secretsDir
func (ds *DockerSecrets) readAll() error {
	secretsDir := ds.GetDir()
	err := isDir(secretsDir)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(secretsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := ds.read(file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// Reads a secret from file
func (ds *DockerSecrets) read(file string) error {
	buf, err := ioutil.ReadFile(ds.GetDir() + "/" + file)
	if err != nil {
		return err
	}
	ds.secrets[file] = strings.TrimSpace(string(buf))
	return nil
}

// Checks if the given path is a directory
func isDir(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.Mode().IsDir() {
		return fmt.Errorf("is not a directory: %s", path)
	}
	return nil
}
