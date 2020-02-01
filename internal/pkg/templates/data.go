package templates

// Data is that data sent to the template parsing
type Data struct {
	ContainerData
	Arguments map[string]string
}

// ContainerData is a part of the template parsing data
type ContainerData struct {
	Labels map[string]string
}
