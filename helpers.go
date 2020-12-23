package main

import (
	"encoding/json"
	"fmt"
	"github.com/xelaj/mtproto/telegram"
	"math/rand"
	"regexp"
	"strings"
)

func ConvertBotAPIEntitiesToClientAPI(entities []MessageEntity) ([]telegram.MessageEntity, error) {
	var res []telegram.MessageEntity
	for _, v := range entities {
		if v.Type == "bold" {
			res = append(res, &telegram.MessageEntityBold{
				Offset: v.Offset,
				Length: v.Length,
			})
		} else if v.Type == "italic" {
			res = append(res, &telegram.MessageEntityItalic{
				Offset: v.Offset,
				Length: v.Length,
			})
		} else if v.Type == "underline" {
			res = append(res, &telegram.MessageEntityUnderline{
				Offset: v.Offset,
				Length: v.Length,
			})
		} else if v.Type == "strikethrough" {
			res = append(res, &telegram.MessageEntityStrike{
				Offset: v.Offset,
				Length: v.Length,
			})
		} else if v.Type == "code" {
			res = append(res, &telegram.MessageEntityCode{
				Offset: v.Offset,
				Length: v.Length,
			})
		} else if v.Type == "pre" {
			res = append(res, &telegram.MessageEntityPre{
				Offset: v.Offset,
				Length: v.Length,
			})
		}
	}
	return res, nil
}

func ParseEntitiesToBotAPI(text, mode string) ([]MessageEntity, error) {
	if mode == "HTML" {
		return ParseEntitiesToBotAPIHTML(text)
	}
	return []MessageEntity{}, fmt.Errorf("unsupported parse mode")
}

func ParseEntitiesToBotAPIHTML(text string) ([]MessageEntity, error) {
	type Tag struct {
		openingTag string
		closingTag string
		tagType    string
		regex      *regexp.Regexp
	}
	tags := []Tag{
		{
			openingTag: "<b>",
			closingTag: "</b>",
			tagType:    "bold",
		},
		{
			openingTag: "<i>",
			closingTag: "</i>",
			tagType:    "italic",
		},
		{
			openingTag: "<u>",
			closingTag: "</u>",
			tagType:    "underline",
		},
		{
			openingTag: "<s>",
			closingTag: "</s>",
			tagType:    "strikethrough",
		},
		{
			openingTag: "<code>",
			closingTag: "</code>",
			tagType:    "code",
		},
		{
			openingTag: "<pre>",
			closingTag: "</pre>",
			tagType:    "pre",
		},
		{
			openingTag: "",
			closingTag: "",
			tagType:    "mention",
			regex:      regexp.MustCompile("(@[a-zA-Z0-9_]{5,32})\\b"),
		},
	}
	canBeContainedList := map[string]bool{
		"bold":          true,
		"italic":        true,
		"underline":     true,
		"strikethrough": true,
		"code":          false,
		"pre":           false,
	}
	var entities []MessageEntity
	var seenPositions []int32
	for _, v := range tags {
		if v.openingTag != "" {
			for i := 0; i < strings.Count(text, v.openingTag); i++ {
				isOriginal := true
				openingPos := int32(strings.Index(text, v.openingTag))
				for tmp := range seenPositions {
					if int32(tmp) == openingPos {
						isOriginal = false
						break
					}
				}
				if openingPos != -1 && isOriginal {
					closingPos := int32(strings.Index(text, v.closingTag))
					if closingPos != -1 && closingPos > openingPos {
						for _, x := range entities {
							if ((openingPos > x.Offset && openingPos < x.Offset+x.Length && closingPos < x.Offset+x.Length) || (openingPos < x.Offset && closingPos > x.Offset+x.Length) || (closingPos < x.Offset+x.Length) || (openingPos > x.Offset+x.Length)) != true {
								return []MessageEntity{}, fmt.Errorf("incorrect containment")
							} else if canBeContainedList[x.Type] != true {
								fmt.Println(1)
								return []MessageEntity{}, fmt.Errorf("incorrect containment")
							}
						}
						entities = append(entities, MessageEntity{
							Type:   v.tagType,
							Offset: openingPos,
							Length: closingPos,
						})
						seenPositions = append(seenPositions, openingPos)
					} else {
						return []MessageEntity{}, fmt.Errorf("can't find closing tag for [%s]", v.openingTag)
					}
				}
			}
		} else if v.regex != nil {
			fmt.Println(v.regex.FindAllString(text, -1))
		}
	}
	return entities, nil
}

func parseInlineKeyboard(kbData ReplyMarkupData) (InlineKeyboard, error) {
	var parsedKb InlineKeyboard
	err := json.Unmarshal([]byte(kbData.Data), &parsedKb)
	if err != nil {
		return InlineKeyboard{}, err
	} else {
		return parsedKb, nil
	}
}

func parseKeyboardToClientAPI(kbData ReplyMarkupData) (telegram.ReplyMarkup, error) {
	if kbData.Type == "remove" {
		type ReplyKeyboardRemove struct {
			Selective bool `json:"selective"`
		}
		var selectiveData ReplyKeyboardRemove
		err := json.Unmarshal([]byte(kbData.Data), &selectiveData)
		if err != nil {
			return nil, err
		} else {
			return &telegram.ReplyKeyboardHide{Selective: selectiveData.Selective}, nil
		}
	} else if kbData.Type == "force_reply" {
		type ForceReply struct {
			Selective bool `json:"selective"`
		}
		var selectiveData ForceReply
		err := json.Unmarshal([]byte(kbData.Data), &selectiveData)
		if err != nil {
			return nil, err
		} else {
			return &telegram.ReplyKeyboardForceReply{Selective: selectiveData.Selective}, nil
		}
	} else if kbData.Type == "reply_markup" {
		maxBtns := 100
		parsedBtns := 0
		type ReplyKeyboard struct {
			Keyboard [][]struct {
				Text            string `json:"text"`
				RequestContact  bool   `json:"request_contact,omitempty"`
				RequestLocation bool   `json:"request_location,omitempty"`
				RequestPoll     *struct {
					Type string `json:"type,omitempty"`
				} `json:"request_poll,omitempty"`
			} `json:"keyboard"`
			ResizeKeyboard  bool `json:"resize_keyboard,omitempty"`
			OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`
			Selective       bool `json:"selective,omitempty"`
		}
		var parsedKb ReplyKeyboard
		err := json.Unmarshal([]byte(kbData.Data), &parsedKb)
		if err != nil {
			return nil, err
		} else {
			var rows []*telegram.KeyboardButtonRow
			for _, kbRow := range parsedKb.Keyboard {
				var row telegram.KeyboardButtonRow
				for _, kbBtn := range kbRow {
					if parsedBtns+1 <= maxBtns {
						var btn telegram.KeyboardButton
						if kbBtn.RequestContact {
							btn = &telegram.KeyboardButtonRequestPhone{Text: kbBtn.Text}
						} else if kbBtn.RequestLocation {
							btn = &telegram.KeyboardButtonRequestGeoLocation{Text: kbBtn.Text}
						} else if kbBtn.RequestPoll != nil {
							btn = &telegram.KeyboardButtonRequestPoll{
								Quiz: kbBtn.RequestPoll.Type == "quiz",
								Text: kbBtn.Text,
							}
						} else {
							btn = &telegram.KeyboardButtonObj{Text: kbBtn.Text}
						}
						row.Buttons = append(row.Buttons, btn)
						parsedBtns = parsedBtns + 1
					}
				}
				rows = append(rows, &row)
			}
			return &telegram.ReplyKeyboardMarkup{
				Resize:    parsedKb.ResizeKeyboard,
				SingleUse: parsedKb.OneTimeKeyboard,
				Selective: parsedKb.Selective,
				Rows:      rows,
			}, nil
		}
	} else if kbData.Type == "inline" {
		maxBtns := 100
		parsedBtns := 0
		var parsedKb InlineKeyboard
		err := json.Unmarshal([]byte(kbData.Data), &parsedKb)
		if err != nil {
			return nil, err
		} else {
			var rows []*telegram.KeyboardButtonRow
			for a, kbRow := range parsedKb.InlineKeyboard {
				var row telegram.KeyboardButtonRow
				for b, kbBtn := range kbRow {
					if parsedBtns+1 <= maxBtns {
						var btn telegram.KeyboardButton
						if kbBtn.Url != "" {
							btn = &telegram.KeyboardButtonURL{Text: kbBtn.Text, URL: kbBtn.Url}
						} else if kbBtn.LoginUrl != nil {
							btn = &telegram.KeyboardButtonURLAuth{
								Text:     kbBtn.Text,
								FwdText:  kbBtn.LoginUrl.ForwardText,
								URL:      kbBtn.LoginUrl.Url,
								ButtonID: rand.Int31(),
							}
						} else if kbBtn.CallbackData != "" {
							btn = &telegram.KeyboardButtonCallback{
								RequiresPassword: false,
								Text:             kbBtn.Text,
								Data:             []byte(kbBtn.CallbackData),
							}
						} else if kbBtn.SwitchInlineQuery != "" {
							btn = &telegram.KeyboardButtonSwitchInline{
								SamePeer: false,
								Text:     kbBtn.Text,
								Query:    kbBtn.SwitchInlineQuery,
							}
						} else if kbBtn.SwitchInlineQueryCurrentChat != "" {
							btn = &telegram.KeyboardButtonSwitchInline{
								SamePeer: true,
								Text:     kbBtn.Text,
								Query:    kbBtn.SwitchInlineQuery,
							}
						} else if kbBtn.CallbackGame != nil {
							if a == 1 && b == 1 {
								btn = &telegram.KeyboardButtonGame{
									Text: kbBtn.Text,
								}
							} else {
								return nil, fmt.Errorf("game button must be first in first row")
							}
						} else if kbBtn.Pay {
							if a == 1 && b == 1 {
								btn = &telegram.KeyboardButtonBuy{
									Text: kbBtn.Text,
								}
							} else {
								return nil, fmt.Errorf("pay button must be first in first row")
							}
						} else {
							return nil, fmt.Errorf("can't parse keyboard")
						}
						row.Buttons = append(row.Buttons, btn)
						parsedBtns = parsedBtns + 1
					} else {
						break
					}
				}
				rows = append(rows, &row)
			}
			return &telegram.ReplyInlineMarkup{
				Rows: rows,
			}, nil
		}
	}
	return nil, fmt.Errorf("unknown keyboard type")
}

func ConvertClientAPIEntitiesToBotAPIEntities(entities []telegram.MessageEntity) ([]MessageEntity, error) {
	var res []MessageEntity
	for _, v := range entities {
		if _, ok := v.(*telegram.MessageEntityMention); ok {
			res = append(res, MessageEntity{
				Type:   "mention",
				Offset: v.(*telegram.MessageEntityMention).Offset,
				Length: v.(*telegram.MessageEntityMention).Length,
			})
		} else if _, ok = v.(*telegram.MessageEntityHashtag); ok {
			res = append(res, MessageEntity{
				Type:   "hashtag",
				Offset: v.(*telegram.MessageEntityHashtag).Offset,
				Length: v.(*telegram.MessageEntityHashtag).Length,
			})
		} else if _, ok = v.(*telegram.MessageEntityCashtag); ok {
			res = append(res, MessageEntity{
				Type:   "cashtag",
				Offset: v.(*telegram.MessageEntityCashtag).Offset,
				Length: v.(*telegram.MessageEntityCashtag).Length,
			})
		} else if _, ok = v.(*telegram.MessageEntityEmail); ok {
			res = append(res, MessageEntity{
				Type:   "email",
				Offset: v.(*telegram.MessageEntityEmail).Offset,
				Length: v.(*telegram.MessageEntityEmail).Length,
			})
		} else if _, ok = v.(*telegram.MessageEntityPhone); ok {
			res = append(res, MessageEntity{
				Type:   "phone_number",
				Offset: v.(*telegram.MessageEntityPhone).Offset,
				Length: v.(*telegram.MessageEntityPhone).Length,
			})
		} else if _, ok = v.(*telegram.MessageEntityBotCommand); ok {
			res = append(res, MessageEntity{
				Type:   "bot_command",
				Offset: v.(*telegram.MessageEntityBotCommand).Offset,
				Length: v.(*telegram.MessageEntityBotCommand).Length,
			})
		}
	}
	return res, nil
}
