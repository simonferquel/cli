// +build !windows

package configfile

func localEnvironment() Environment {
	return Environment{
		Docker: &DockerEnvironment{
			Host: "/var/run/docker.sock",
		},
	}
}
