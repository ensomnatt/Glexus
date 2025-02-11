package wsserver

type wsMessage struct {
  Time string `json:"time,omitempty"`
  IP string `json:"ip,omitempty"`
  Action string `json:"action"`
  UserNames []string `json:"usernames,omitempty"`
}

type httpMessage struct {
  Type string `json:"type"`
  VideoFiles []string `json:"videofiles,omitempty"`
}
