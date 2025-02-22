package wsserver

type wsMessage struct {
  Time float64 `json:"time,omitempty"`
  IP string `json:"ip,omitempty"`
  Action string `json:"action"`
  UserNames []string `json:"usernames,omitempty"`
}

type httpMessage struct {
  Type string `json:"type"`
  VideoFiles []string `json:"videofiles,omitempty"`
  VideoDir string `json:"videodir,omitempty"`
}
