package messages

type LoginMessage struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (lm *LoginMessage) IsValid() bool {
	return true
}
