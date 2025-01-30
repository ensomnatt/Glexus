package wsserver

type wsMessage struct {
  Time string `json:"time,omitempty"`
  IP string `json:"ip,omitempty"`
  Action string `json:"action"`
  UserNames []string `json:"usernames,omitempty"`
}
