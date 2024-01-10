package forms

type MessageForm struct {
	MessageType int32  `json:"message_type" binging:"required,oneof=1 2 3 4 5"`
	Subject     string `json:"subject" binging:"required"`
	Message     string `json:"message" binging:"required"`
	File        string `json:"file" binging:"required"`
}
