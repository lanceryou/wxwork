package wxwork

type wxWorkRequest struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentID int64  `json:"agentid"`
	Safe    int    `json:"safe"`
	Text    Text   `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

func (w wxWorkRequest) Valid() bool {
	return len(w.ToUser) != 0 || len(w.ToParty) != 0 || len(w.ToTag) != 0
}

type wxWorkReply struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

type wxWorkTokenReply struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Token     string `json:"access_token"`
	ExpireTTL int    `json:"expires_in"`
}
