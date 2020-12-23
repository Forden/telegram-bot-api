package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/xelaj/mtproto/telegram"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route
type BotAPI struct {
	cache    map[string]ClientAPI
	schemas  map[string]interface{}
	authData struct {
		AppId   int
		AppHash string
	}
	mx sync.Mutex
}

func NewRouter(cfg MyConfig) *mux.Router {
	r := mux.NewRouter()
	api := BotAPI{
		cache: make(map[string]ClientAPI),
	}
	api.authData.AppId = cfg.Auth.AppId
	api.authData.AppHash = cfg.Auth.AppHash

	s := r.PathPrefix("/bot{token:\\d{0,10}\\:[a-zA-Z0-9_\\-]{35}}").Subrouter()
	for _, route := range generateRoutes(&api) {
		s.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(api.AuthorizeHandler(route.HandlerFunc))
	}
	return r
}

func (api BotAPI) parseRequest(r *http.Request, parseStruct interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return fmt.Errorf("can't read body")
	} else {
		err = json.Unmarshal(body, &parseStruct)
		if err != nil {
			log.Error(err)
			return fmt.Errorf("can't parse body to struct")
		}
		return nil
	}

}

func (api BotAPI) AuthorizeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		token := vars["token"]
		var newClient *ClientAPI
		var fromCache = false
		if client, ok := api.cache[token]; ok {
			newClient = &client
			fmt.Println("from cache")
			fromCache = true
		} else {
			fmt.Println("new client")
			newClient = &ClientAPI{}
			newClient.client, err = telegram.NewClient(telegram.ClientConfig{
				SessionFile:    "./authData/" + token + "session.json",
				ServerHost:     "149.154.167.50:443",
				PublicKeysFile: "./authData/keys.pem",
				AppVersion:     "0.1",
				AppID:          api.authData.AppId,
				AppHash:        api.authData.AppHash,
			})
			if err != nil {
				log.Panicf("new session error: %s", err)
			}
			_, err := newClient.client.AuthImportBotAuthorization(1, int32(api.authData.AppId), api.authData.AppHash, token)
			if err != nil {
				log.Panicf("login error: %s", err)
				http.Error(w, "can't login with this token: [%s]", 401)
				return
			}
		}
		tmp, _ := strconv.ParseInt(strings.Split(token, ":")[0], 10, 64)
		newClient.selfData.ID = tmp

		if fromCache != true {
			api.mx.Lock() // possible race condition
			defer api.mx.Unlock()
			api.cache[token] = *newClient
		}

		h.ServeHTTP(w, r)
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	})
}

func (api BotAPI) getMe(w http.ResponseWriter, r *http.Request) {
	var response APIResponse
	var user User
	vars := mux.Vars(r)
	token := vars["token"]
	client := api.cache[token]
	user.ID = client.selfData.ID
	user.IsBot = true
	users := []telegram.InputUser{&telegram.InputUserObj{
		UserID:     int32(user.ID),
		AccessHash: int64(0),
	}}
	fmt.Println(1)
	usersInfo, err := client.client.UsersGetUsers(users) // softlock by mtproto lib
	fmt.Println(2)
	if err != nil {
		log.Errorf("error getting self %s", err)
	} else {
		fmt.Println(usersInfo[0].(*telegram.UserObj))
		userInfo := usersInfo[0].(*telegram.UserObj)
		user.FirstName = userInfo.FirstName
		user.LastName = userInfo.LastName
		user.LanguageCode = userInfo.LangCode
	}
	response.Result = user
	resp, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	_, _ = w.Write(resp)
}

func (api BotAPI) sendMessage(w http.ResponseWriter, r *http.Request) {
	var response APIResponse
	var req SendMessage
	response.Ok = true
	err := api.parseRequest(r, &req)
	if err != nil {
		response.Ok = false
		response.Description = fmt.Sprintf("%s", err)
		response.ErrorCode = 400
	} else {
		vars := mux.Vars(r)
		token := vars["token"]
		client := api.cache[token]
		var opts SendMessageOpts
		opts.ParseMode = req.ParseMode
		var entities []MessageEntity
		if opts.ParseMode != "" {
			entities, err = ParseEntitiesToBotAPI(req.Text, opts.ParseMode)
			if err != nil {
				log.Errorf("error parsing entities: %s", err)
				response.Ok = false
				response.Description = fmt.Sprintf("error parsing entities: %s", err)
				response.ErrorCode = 400
			}
		}
		opts.ReplyMarkup, err = parseKeyboardToClientAPI(*req.ReplyMarkup)
		if err != nil {
			log.Errorf("error parsing keyboard: %s", err)
			response.Ok = false
			response.Description = fmt.Sprintf("error parsing keyboard: %s", err)
			response.ErrorCode = 400
		}
		if response.Ok {
			opts.DisableWebPagePreview = req.DisableWebPagePreview
			opts.DisableNotification = req.DisableNotification
			opts.ReplyToMessageId = int32(req.ReplyToMessageID)
			opts.AllowSendingWithoutReply = req.AllowSendingWithoutReply
			opts.Entities, err = ConvertBotAPIEntitiesToClientAPI(entities)
			if err != nil {
				log.Error(err)
			}
			upds, err := client.SendMessage(int32(req.ChatID), req.Text, &opts)
			if err != nil {
				log.Errorf("error [%s] while sending message\n", err)
				response.Ok = false
				response.Description = fmt.Sprintf("%s", err)
			} else {
				res := upds.(*telegram.UpdateShortSentMessage)
				var sentMessage Message
				clientEntitites, _ := ConvertClientAPIEntitiesToBotAPIEntities(res.Entities)
				entities = append(entities, clientEntitites...)
				response.Ok = true
				response.Result = sentMessage
				sentMessage.Entities = &entities
				sentMessage.MessageId = int64(res.ID)
				sentMessage.Date = res.Date
				sentMessage.Text = req.Text
				sentMessage.Chat = &Chat{
					Id: req.ChatID,
				}
				if req.ReplyMarkup.Type == "inline" {
					parsedInlineKb, err := parseInlineKeyboard(*req.ReplyMarkup)
					if err == nil {
						sentMessage.ReplyMarkup = &parsedInlineKb
					}
				}
				response.Result = sentMessage
			}
		}
	}
	resp, _ := json.Marshal(response)
	_, _ = w.Write(resp)
}

//goland:noinspection GoUnusedParameter
func (api BotAPI) testMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	client := api.cache[token]
	a, _ := client.client.MessagesGetMessages([]telegram.InputMessage{&telegram.InputMessageID{ID: 178}, &telegram.InputMessageID{ID: 177}})
	for _, v := range a.(*telegram.MessagesMessagesObj).Messages {
		fmt.Println(v.(*telegram.MessageObj).ReplyMarkup.(*telegram.ReplyInlineMarkup).Rows[0].Buttons[0])
	}
}

func generateRoutes(api *BotAPI) Routes {
	return Routes{
		Route{
			"getMe",
			"POST",
			"/getMe",
			api.getMe,
		},
		Route{
			"sendMessage",
			"POST",
			"/sendMessage",
			api.sendMessage,
		},
		Route{
			"test",
			"POST",
			"/test",
			api.testMethod,
		},
	}

}
