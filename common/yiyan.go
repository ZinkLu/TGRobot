package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const YIYAN_URL = "https://international.v1.hitokoto.cn"

func GetYiYan() (*YiYan, error) {
	yiyan := &YiYan{}
	res, err := http.Get(YIYAN_URL)
	if err != nil {
		return yiyan, err
	}
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, yiyan)
	if err != nil {
		return yiyan, err
	}
	return yiyan, err

}

type YiYan struct {
	Id          int    `json:"id"`
	Uuid        string `json:"uuid"`
	Hitokoto    string `json:"hitokoto"`
	Type        string `json:"type"`
	From        string `json:"from"`
	From_who    string `json:"from_who"`
	Creator     string `json:"creator"`
	Creator_uid int    `json:"creator_uid"`
	Reviewer    int    `json:"reviewer"`
	Commit_from string `json:"commit_from"`
	Created_at  string `json:"created_at"`
	Length      int    `json:"length"`
}

func (yiyan *YiYan) Quote() string {
	return fmt.Sprintf("%s -- %s《%s》", yiyan.Hitokoto, yiyan.From_who, yiyan.From)
}
