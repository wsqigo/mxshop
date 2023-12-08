package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	tmp := fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(tmp), nil
}

type UserInfoResp struct {
	ID       int32    `json:"id"`
	NickName string   `json:"nick_name"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
