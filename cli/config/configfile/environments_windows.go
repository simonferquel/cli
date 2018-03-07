package configfile

func localEnvironment() Environment {
	return Environment{
		Docker: &DockerEnvironment{
			Host: "npipe:////./pipe/docker_engine",
		},
	}
}
