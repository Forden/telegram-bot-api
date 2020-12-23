package main

import (
	"fmt"
	"github.com/xelaj/mtproto/telegram"
	"math/rand"
)

type ClientAPI struct {
	client   *telegram.Client
	selfData User
}

func getPeerType(peerId int64) string {
	var MinChannelId, MaxChannelId, MinChatId, MaxUserId int64
	MinChannelId = -1002147483647
	MaxChannelId = -1000000000000
	MinChatId = -2147483647
	MaxUserId = 2147483647
	if (0 < peerId) && (peerId <= MaxUserId) {
		return "user"
	} else if peerId < 0 {
		if MinChatId <= peerId {
			return "chat"
		} else if (peerId <= MaxChannelId) && (MinChannelId <= peerId) {
			return "channel"
		}
	}
	return ""
}

func (tg ClientAPI) UpdateSelfInfo() {
	//
}

type SendMessageOpts struct {
	ParseMode                string
	Entities                 []telegram.MessageEntity
	DisableWebPagePreview    bool
	DisableNotification      bool
	ReplyToMessageId         int32
	AllowSendingWithoutReply bool
	ReplyMarkup              telegram.ReplyMarkup
}

func (tg ClientAPI) SendMessage(chatId int32, text string, opts *SendMessageOpts) (telegram.Updates, error) {
	chatType := getPeerType(int64(chatId))
	params := &telegram.MessagesSendMessageParams{}
	params.NoWebpage = opts.DisableWebPagePreview
	params.Silent = opts.DisableNotification
	params.Background = false
	params.ClearDraft = true
	params.ReplyToMsgID = opts.ReplyToMessageId
	params.Message = text
	params.RandomID = rand.Int63()
	params.ReplyMarkup = opts.ReplyMarkup
	params.Entities = opts.Entities
	if chatType == "user" {
		params.Peer = &telegram.InputPeerUser{
			UserID:     chatId,
			AccessHash: 0,
		}
	} else if chatType == "chat" {
		params.Peer = &telegram.InputPeerChat{
			ChatID: chatId,
		}
	} else if chatType == "channel" {
		params.Peer = &telegram.InputPeerChannel{
			ChannelID:  chatId,
			AccessHash: 0,
		}
	}
	if params.Peer != nil {
		return tg.client.MessagesSendMessage(params)
	}
	return nil, fmt.Errorf("unknown peer type for id %d", chatId)
}

func (tg ClientAPI) ResolveUsername(username string) int32 {
	resolved, err := tg.client.ContactsResolveUsername(username)
	if err != nil {
		return 0
	} else {
		if len(resolved.Chats) != 0 {
			return resolved.Chats[0].(*telegram.ChatObj).ID
		} else if len(resolved.Users) != 0 {
			return resolved.Users[0].(*telegram.UserObj).ID
		} else {
			return 0
		}
	}
}
