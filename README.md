# discriminator
[![GoDoc](https://godoc.org/github.com/sidusIO/discriminator?status.svg)](https://godoc.org/github.com/sidusIO/discriminator)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/sidusio/discriminator)](https://goreportcard.com/report/github.com/sidusio/discriminator)
[![Build Status](https://travis-ci.com/sidusIO/discriminator.svg?branch=master)](https://travis-ci.com/sidusIO/discriminator)
[![Image size](https://images.microbadger.com/badges/image/sidusio/discriminator.svg)](https://microbadger.com/images/sidusio/discriminator "Get your own image badge on microbadger.com")
[![docker hub](https://images.microbadger.com/badges/version/sidusio/discriminator.svg)](https://hub.docker.com/r/sidusio/discriminator "Get your own version badge on microbadger.com")

Automatically add labels to docker containers.

## Getting Started
This application works by looking for a labels (default: `io.sidus.discriminator`) associated value (instructions)
and specified templates to add or remove certain labels from that container.

WARNING: In order to modify the labels in a container, discriminator has to create a new one and remove the old one.

WARNING: This application is in beta, use at own risk.

### Docker
The easiest way to run the application is through docker [`sidusio/discriminator`](https://hub.docker.com/r/sidusio/discriminator)

When you run the application in docker you have to mount in:
1. Your docker socket: `/var/run/docker.sock:/var/run/docker.sock`
2. A directory with your templates to `/templates`

Example run command:
```
docker run -v /var/run/docker.sock:/var/run/docker.sock -v yourTemplatesDirectory:/templates sidusio/discriminator
``` 

### Configuration
The application is configured through environment variables

| Enviornment Variable                     | Default value          | Description                                                |
|:-----------------------------------------|:-----------------------|:-----------------------------------------------------------|
| DISCRIMINATOR_TEMPLATES_PATH             | /templates             | Directory with your templates                              |
| DISCRIMINATOR_TEMPLATES_EXTENSION        | .tmpl                  | The extension of your templates                            |
| DISCRIMINATOR_CONTAINERS_LABEL           | io.sidus.discriminator | The label to look at for instructions                      |
| DISCRIMINATOR_INCLUDE_STOPPED_CONTAINERS | false                  | Whether to run the application on stopped containers       |
| DISCRIMINATOR_RUN_INTERVAL               | 5m                     | How often the application should go through the containers |
| DISCRIMINATOR_LOG_LEVEL                  | info                   | debug/info/warn/error                                      |
| DISCRIMINATOR_LOG_FORMAT                 | text                   | text/json                                                  |


### Templates
Templates are called by instructions to modify the labels of the container.

Templates are parsed with [`text/template`](https://golang.org/pkg/text/template/).

A row starting with a `+` is a label (key and value) that should be added/overwritten ex: `+my.label=value`

A row starting with a `-` is a label (only key) that should be removed ex: `-my.label`

All other rows will be discarded.

The following data is sent to the template parser and can be used in the template:
```golang
type Data struct {
	ContainerData struct {
		Labels map[string]string
		Name   string
	}
	Arguments map[string]string
}
```
Where arguments are the ones specified in the instruction.

### Instructions
Instruction are specified with the application label (default `io.sidus.discriminator`).

An example instruction could be `testtemplate()` which would apply the template named `testtemplate` (without any file extensions).

You can also pass arguments to instructions: `template(argument: value, argument2: value)`

Instructions can be chained and will the be applied from left to right:
`template1() | template2(arg: value)  | template3()`

## Contributing
Contributions are welcome!

For example, you could contribute by expanding this section.
