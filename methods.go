package main

type APIResponse struct {
	Ok          bool                `json:"ok"`
	Description string              `json:"description,omitempty"`
	Result      interface{}         `json:"result,omitempty"`
	ErrorCode   int64               `json:"error_code,omitempty"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

type ResponseParameters struct {
	MigrateToChatId int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      int64 `json:"retry_after,omitempty"`
}

type SendMessage struct {
	ChatID                   int64            `json:"chat_id"`
	Text                     string           `json:"text"`
	ParseMode                string           `json:"parse_mode,omitempty"`
	DisableWebPagePreview    bool             `json:"disable_web_page_preview,omitempty"`
	DisableNotification      bool             `json:"disable_notification,omitempty"`
	ReplyToMessageID         int64            `json:"reply_to_message_id,omitempty"`
	AllowSendingWithoutReply bool             `json:"allow_sending_without_reply,omitempty"`
	ReplyMarkup              *ReplyMarkupData `json:"reply_markup,omitempty"`
	Entities                 *[]MessageEntity `json:"entities,omitempty"`
}
