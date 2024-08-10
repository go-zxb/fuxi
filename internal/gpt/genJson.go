package gpt

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-zxb/fuxi/config"
	"github.com/go-zxb/fuxi/internal/model"
)

func GenCode(question string) (*model.CodeModel, string, error) {
	conf := config.GetConfig()
	gpt := conf.GPT
	chat := NewChat(gpt)

	str, err := chat.Chat(question + ".忽略掉创建时间,更新时间,删除时间,主键id等字段.")
	if err != nil {
		return nil, "", errors.New("❎ 请求Ai失败:" + err.Error())
	}

	fmt.Println(str)

	var code *model.CodeModel
	err = json.Unmarshal([]byte(str), &code)
	if err != nil {
		return nil, "", errors.New("❎ json.Unmarshal 失败:" + err.Error())
	}
	return code, str, err
}
