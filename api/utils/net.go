package utils

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	logger2 "geekai/logger"
	"io"
	"net/http"
	"net/url"
)

var logger = logger2.GetLogger()

// ReplyChunkMessage 回复客户片段端消息
func ReplyChunkMessage(client *types.WsClient, message interface{}) {
	msg, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Error for decoding json data: %v", err.Error())
		return
	}
	err = client.Send(msg)
	if err != nil {
		logger.Errorf("Error for reply message: %v", err.Error())
	}
}

// ReplyMessage 回复客户端一条完整的消息
func ReplyMessage(ws *types.WsClient, message interface{}) {
	ReplyChunkMessage(ws, types.ReplyMessage{Type: types.WsContent, Content: message})
	ReplyChunkMessage(ws, types.ReplyMessage{Type: types.WsEnd})
}

func ReplyContent(ws *types.WsClient, message interface{}) {
	ReplyChunkMessage(ws, types.ReplyMessage{Type: types.WsContent, Content: message})
}

// ReplyErrorMessage 向客户端发送错误消息
func ReplyErrorMessage(ws *types.WsClient, message interface{}) {
	ReplyChunkMessage(ws, types.ReplyMessage{Type: types.WsErr, Content: message})
}

func DownloadImage(imageURL string, proxy string) ([]byte, error) {
	var client *http.Client
	if proxy == "" {
		client = http.DefaultClient
	} else {
		proxyURL, _ := url.Parse(proxy)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	}
	request, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

func GetBaseURL(strURL string) string {
	u, err := url.Parse(strURL)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}
