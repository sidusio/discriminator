package main

import (
	"github.com/sirupsen/logrus"

	"sidus.io/discriminator/internal/app/discriminator"
)

func main() {
	err := discriminator.Run()
	if err != nil {
		logrus.WithError(err).Fatalf("Application stopped with error")
	}
	logrus.Info("Application closed")
}
