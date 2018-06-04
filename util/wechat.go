package util

import (
	"github.com/astaxie/beego"
	"sort"
	"crypto/sha1"
	"fmt"
	"encoding/xml"
	"time"
	"strconv"
)

var accountId, token, appId string

// 消息类型
const (
	MSGTYPE_TEXT       = "text"
	MSGTYPE_IMAGE      = "image"
	MSGTYPE_VOICE      = "voice"
	MSGTYPE_VIDEO      = "video"
	MSGTYPE_SHORTVIDEO = "shortvideo"
	MSGTYPE_LOCATION   = "location"
	MSGTYPE_LINK       = "link"
	MSGTYPE_EVENT      = "event"
)

// 事件推送类型
const (
	EVENT_SUBSCRIBE   = "subscribe"
	EVENT_UNSUBSCRIBE = "unsubscribe"
	EVENT_SCAN        = "scan"
	EVENT_LOCATION    = "location"
	EVENT_CLICK       = "click"
)

// 基础数据
type BaseData struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
}

type CDATA struct {
	Content string `xml:",cdata"`
}

// 文本消息
type TextMsg struct {
	Base    BaseData
	Content string `xml:"Content"`
	MsgId   string `xml:"MsgId"`
}

type ReplyTextMsg struct {
	XMLName      string `xml:"xml"`
	ToUserName   CDATA  `xml:"ToUserName"`
	FromUserName CDATA  `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      CDATA  `xml:"MsgType"`
	Content      CDATA  `xml:"Content"`
}

func init() {
	accountId = beego.AppConfig.String("wechat_account_id")
	appId = beego.AppConfig.String("wechat_app_id")
	token = beego.AppConfig.String("wechat_token")
}

// 验证签名
func AuthVerify(signature, timestamp, nonce, echostr string) (string, bool) {
	// 将参数排序和拼接
	str := sort.StringSlice{token, timestamp, nonce}
	sort.Sort(str)
	sortStr := ""
	for _, v := range str {
		sortStr += v
	}

	// 进行sha1加密
	sh := sha1.New()
	sh.Write([]byte(sortStr))
	encryptStr := fmt.Sprintf("%x", sh.Sum(nil))

	// 将本地计算的签名和微信传递过来的签名进行对比
	if encryptStr == signature {
		return echostr, true
	}

	return "Invalid Signature.", false
}

// 解析事件推送基础数据
func ParseBase(base []byte) (BaseData, error) {
	baseEvent := BaseData{}

	err := xml.Unmarshal(base, &baseEvent)

	return baseEvent, err
}

// 解析文本消息
func ParseMsg(msg []byte) (TextMsg, error) {
	textMsg := TextMsg{}

	err := xml.Unmarshal(msg, &textMsg)

	return textMsg, err
}

// 回复文本消息
func ReplayTextMsg(openid, msg string) ([]byte, error) {
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	replay := ReplyTextMsg{
		ToUserName:   CDATA{openid},
		FromUserName: CDATA{accountId},
		CreateTime:   currentTime,
		MsgType:      CDATA{"text"},
		Content:      CDATA{msg},
	}

	return xml.MarshalIndent(replay, "", "  ")
}
