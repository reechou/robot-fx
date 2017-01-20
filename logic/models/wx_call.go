package models

import (
	"time"
)

type WXCallCheck struct {
	Signature string
	Timestamp string
	Nonce     string
	Echostr   string
}

type WXCallRequest struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	Event        string
	EventKey     string
	Ticket       string
	MsgId        int
}

type WXCallResponse struct {
	ToUserName   string `xml:"xml>ToUserName"`
	FromUserName string `xml:"xml>FromUserName"`
	CreateTime   string `xml:"xml>CreateTime"`
	MsgType      string `xml:"xml>MsgType"`
	Content      string `xml:"xml>Content"`
	MsgId        int    `xml:"xml>MsgId"`
}
