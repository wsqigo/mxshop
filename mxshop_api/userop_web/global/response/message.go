package response

type MessageResp struct {
	Id      int32  `json:"id"`
	UserId  int32  `json:"user_id"`
	Type    int32  `json:"type"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	File    string `json:"file"`
}
