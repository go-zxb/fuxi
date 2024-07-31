package gpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-zxb/fuxi/config"
	"github.com/go-zxb/fuxi/internal/model"
	"io"
	"net/http"
)

type MsgHistory struct {
	Msg string `json:"msg"`
}

type Chat struct {
	gpt config.GPT
}

func (c Chat) name() {

}

func NewChat(gpt config.GPT) *Chat {
	return &Chat{gpt: gpt}
}

func (c Chat) DeepSeekChat(query string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	fmt.Println("-------------------DeepSeek----------------------------")

	client := &http.Client{}

	msg := model.RequestModel{
		Model:          c.gpt.DeepSeek.Model,
		Temperature:    c.gpt.Temperature,
		ResponseFormat: map[string]string{"type": "json_object"},
		Messages: []*model.Message{
			{Role: "system", Content: c.gpt.Prompt},
			{Role: "user", Content: query},
		},
	}

	requestBodyJSON, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return "", errors.New("Error marshalling request body:" + err.Error())
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.gpt.DeepSeek.BaseURL, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "", errors.New("Error creating HTTP request:" + err.Error())
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.gpt.DeepSeek.ApiKey))

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return "", errors.New("Error sending HTTP request:" + err.Error())
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Error ReadAll Body:" + err.Error())
	}
	//fmt.Println(string(all))

	response := model.ChatGptResp{}
	err = json.Unmarshal(all, &response)
	if err != nil {
		fmt.Println("Error Unmarshal:", err)
		return "", errors.New("Error Unmarshal:" + err.Error())
	}

	if len(response.Choices) == 0 {
		fmt.Println(string(all))
		return "", errors.New("数据丢失")
	}
	//fmt.Println("-------------------DeepSeekChat END----------------------------")
	fmt.Println(response.Choices[0].Message.Content)

	return response.Choices[0].Message.Content, nil
}

func (c Chat) KimiChat(query string) (string, error) {
	fmt.Println("-------------------Kimi----------------------------")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	client := &http.Client{}

	msg := model.RequestModel{
		Model:          c.gpt.Kimi.Model,
		Temperature:    c.gpt.Temperature,
		ResponseFormat: map[string]string{"type": "json_object"},
		Messages: []*model.Message{
			{Role: "system", Content: c.gpt.Prompt},
			{Role: "user", Content: query},
		},
	}

	requestBodyJSON, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return "", errors.New("Error marshalling request body:" + err.Error())
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.gpt.Kimi.BaseURL, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "", errors.New("Error creating HTTP request:" + err.Error())
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.gpt.Kimi.ApiKey))

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return "", errors.New("Error sending HTTP request:" + err.Error())
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(string(all))

	response := model.ChatGptResp{}
	err = json.Unmarshal(all, &response)
	if err != nil {
		fmt.Println("Error Unmarshal:", err)
		return "", errors.New("Error Unmarshal:" + err.Error())
	}

	if len(response.Choices) == 0 {
		fmt.Println(string(all))
		return "", errors.New("数据丢失")
	}

	fmt.Println("-------------------Kimi END----------------------------")
	//fmt.Println(response.Choices[0].Message.Content)
	return response.Choices[0].Message.Content, nil
}

func (c Chat) Chat(query string) (string, error) {
	switch c.gpt.ChatGPTPlatform {
	case model.KimiType:
		return c.KimiChat(query)
	case model.DeepSeekType:
		return c.DeepSeekChat(query)
	default:
		return c.DeepSeekChat(query)
	}
}
