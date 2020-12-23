package main

type Message struct {
	MessageId            int64            `json:"message_id,omitempty"`
	From                 *User            `json:"from,omitempty"`
	Date                 int32            `json:"date"`
	Chat                 *Chat            `json:"chat,omitempty"`
	SenderChat           *Chat            `json:"sender_chat,omitempty"`
	ForwardFrom          *User            `json:"forward_from"`
	ForwardFromChat      *Chat            `json:"forward_from_chat"`
	ForwardFromMessageId int64            `json:"forward_from_message_id"`
	ForwardSignature     string           `json:"forward_signature,omitempty"`
	ForwardSenderName    string           `json:"forward_sender_name,omitempty"`
	ForwardDate          int32            `json:"forward_date,omitempty"`
	ReplyToMessage       *Message         `json:"reply_to_message,omitempty"`
	ViaBot               bool             `json:"via_bot,omitempty"`
	EditDate             int32            `json:"edit_date,omitempty"`
	MediaGroupId         string           `json:"media_group_id,omitempty"`
	AuthorSignature      string           `json:"author_signature,omitempty"`
	Text                 string           `json:"text,omitempty"`
	Entities             *[]MessageEntity `json:"entities,omitempty"`
	ReplyMarkup          *InlineKeyboard  `json:"reply_markup,omitempty"`
}

type Chat struct {
	Id               int64            `json:"id"`
	Type             string           `json:"type"`
	Title            string           `json:"title,omitempty"`
	Username         string           `json:"username,omitempty"`
	FirstName        string           `json:"first_name,omitempty"`
	LastName         string           `json:"last_name,omitempty"`
	Photo            *ChatPhoto       `json:"photo,omitempty"`
	Bio              string           `json:"bio,omitempty"`
	Description      string           `json:"description,omitempty"`
	InviteLink       string           `json:"invite_link,omitempty"`
	PinnedMessage    *Message         `json:"pinned_message,omitempty"`
	Permissions      *ChatPermissions `json:"permissions,omitempty"`
	SlowModeDelay    int32            `json:"slow_mode_delay,omitempty"`
	StickerSetName   string           `json:"sticker_set_name,omitempty"`
	CanSetStickerSet bool             `json:"can_set_sticker_set,omitempty"`
	LinkedChatId     int64            `json:"linked_chat_id,omitempty"`
	Location         *ChatLocation    `json:"location,omitempty"`
}

type ChatPhoto struct {
	SmallFileId       string `json:"small_file_id"`
	SmallFileUniqueId string `json:"small_file_unique_id"`
	BigFileId         string `json:"big_file_id"`
	BigFileUniqueId   string `json:"big_file_unique_id"`
}

type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`
	CanSendMediaMessages  bool `json:"can_send_media_messages,omitempty"`
	CanSendPolls          bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
}

type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

type Location struct {
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int64   `json:"live_period,omitempty"`
	Heading              int64   `json:"heading,omitempty"`
	ProximityAlertRadius int64   `json:"proximity_alert_radius,omitempty"`
}

type InlineKeyboard struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text     string `json:"text"`
	Url      string `json:"url,omitempty"`
	LoginUrl *struct {
		Url                string `json:"url"`
		ForwardText        string `json:"forward_text,omitempty"`
		BotUsername        string `json:"bot_username,omitempty"`
		RequestWriteAccess bool   `json:"request_write_access,omitempty"`
	} `json:"login_url,omitempty"`
	CallbackData                 string      `json:"callback_data,omitempty"`
	SwitchInlineQuery            string      `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string      `json:"switch_inline_query_current_chat,omitempty"`
	CallbackGame                 interface{} `json:"callback_game,omitempty"`
	Pay                          bool        `json:"pay,omitempty"`
}

type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int32  `json:"offset"`
	Length   int32  `json:"length"`
	URL      string `json:"url,omitempty"`
	User     *User  `json:"user,omitempty"`
	Language string `json:"language,omitempty"`
}

type User struct {
	ID                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`
	Username                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
}

type ReplyMarkupData struct {
	Type string `json:"keyboard_type"`
	Data string `json:"data"`
}
