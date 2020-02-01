package docker

// Container is a simplistic representation of
// a docker container
type Container struct {
	Name   string
	ID     string
	Labels map[string]string
}
