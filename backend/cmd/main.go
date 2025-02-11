package main

import (
	wsserver "glexus/backend/internal/WSserver"
  logrus "github.com/sirupsen/logrus"
)

func main() {
  logrus.SetLevel(logrus.DebugLevel)
  wsserver := wsserver.NewServer("0.0.0.0:6969")
  logrus.Info("server started")
  err := wsserver.Start()
  if err != nil {
    logrus.Fatalf("ListenAndServe: %v", err)
  }
}
