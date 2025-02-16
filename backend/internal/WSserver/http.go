package wsserver

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *Server) sendVideoFiles(w http.ResponseWriter, _ *http.Request) {
  msg := httpMessage{
    Type: "sendVideoFiles",
    VideoFiles: s.config.VideoFiles,
  }

  err := json.NewEncoder(w).Encode(msg)
  if err != nil {
    logrus.Errorf("failed to encode json: %v", err)
  }
}

func (s *Server) sendVideoDir(w http.ResponseWriter, _ *http.Request) {
  msg := httpMessage{
    Type: "sendVideoDir",
    VideoDir: s.config.VideoDir,
  }

  err := json.NewEncoder(w).Encode(msg)
  if err != nil {
    logrus.Errorf("failed to encode json: %v", err)
  }
}
