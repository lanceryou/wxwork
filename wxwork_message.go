package wxwork

import (
	"context"
	"fmt"
	"github.com/lanceryou/wxwork/http_client"
	"time"
)

type WxWorkMessage struct {
	opt   WxWorkMessageOptions
	cache map[int64]messageCache
}

type messageCache struct {
	expireTime int64
	token      string
}

func (c messageCache) ValidToken() bool {
	return time.Now().Unix() < c.expireTime
}

const (
	ValidToken = 42001
)

func (w *WxWorkMessage) SendMessage(ctx context.Context, applicationName string, targets string, message string) (err error) {
	company := matchCompanyInfo(applicationName, w.opt.companyInfoList)
	if company == nil {
		err = fmt.Errorf("match company fail applicationName:%v", applicationName)
		return
	}

	token, err := w.getToken(ctx, company)
	if err != nil {
		return
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
	body := &wxWorkRequest{
		MsgType: "text",
		AgentID: company.AgentID,
		ToUser:  targets,
		Text: Text{
			Content: message,
		},
	}

	if !body.Valid() {
		err = fmt.Errorf("request error:%v", *body)
		return
	}

	var reply wxWorkReply
	if err = w.opt.client.Post(ctx, url, body, &reply); err != nil {
		return
	}

	if reply.ErrCode == 0 {
		return
	}

	if reply.ErrCode == ValidToken {
		delete(w.cache, company.AgentID)
		return w.SendMessage(ctx, applicationName, targets, message)
	}
	return
}

func (w *WxWorkMessage) getToken(ctx context.Context, info *CompanyInfo) (token string, err error) {
	cache, ok := w.cache[info.AgentID]
	if ok && cache.ValidToken() {
		return cache.token, nil
	}

	// 重新获取token
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", info.CompanyID, info.ApplicationID)

	var reply wxWorkTokenReply
	if err = w.opt.client.Get(ctx, url, &reply); err != nil {
		return
	}

	if reply.ErrCode != 0 {
		err = fmt.Errorf("errmsg:%v", reply.ErrMsg)
		return
	}

	w.cache[info.AgentID] = messageCache{
		token:      reply.Token,
		expireTime: time.Now().Add(time.Duration(reply.ExpireTTL) * time.Second).Unix(),
	}
	return reply.Token, nil
}

func matchCompanyInfo(applicationName string, companyInfoList []CompanyInfo) *CompanyInfo {
	for _, company := range companyInfoList {
		if company.ApplicationName == applicationName {
			return &company
		}
	}
	return nil
}

func NewWxWorkMessage(opts ...WxWorkMessageOption) *WxWorkMessage {
	opt := WxWorkMessageOptions{
		client: http_client.NewHttpClient(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return &WxWorkMessage{
		opt:   opt,
		cache: make(map[int64]messageCache),
	}
}
