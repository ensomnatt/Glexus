package wsserver

import (
	"html/template"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type watchData struct {
  VideoDir string
}

func (s *Server) watch(w http.ResponseWriter, r *http.Request) {
  data := watchData{
    VideoDir: s.config.VideoDir,
  }

  if _, err := os.Stat(data.VideoDir); os.IsNotExist(err) {
    logrus.Fatalf("file doesn't exists")
  }

  tmpl, err := template.ParseFiles("frontend/src/pages/watch.html")
  if err != nil {
    logrus.Fatalf("failed to parse template: %v", err)
  }

  err = tmpl.Execute(w, data)
  if err != nil {
    logrus.Fatalf("failed to apply template: %v", err)
  }
}
