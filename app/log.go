package main

import "github.com/sirupsen/logrus"

var Logger = logrus.New()

func init() {
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}
