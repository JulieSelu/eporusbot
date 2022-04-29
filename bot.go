package main

import (
	"context"
	"encoding/json"
	"github.com/andamound/telegram-esperanto-bot/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

type Request struct {
	Body string `json:"body"`
}

type TelegramResponse struct {
	Method           string `json:"method"`
	ChatId           int64  `json:"chat_id"`
	Text             string `json:"text"`
	ReplyToMessageID int    `json:"reply_to_message_id"`
}

type Response struct {
	StatusCode      uint              `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            TelegramResponse  `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

func handleUpdate(update *tgbotapi.Update) TelegramResponse {
	if update.Message == nil || !update.Message.IsCommand() {
		return TelegramResponse{}
	}

	arguments := update.Message.CommandArguments()
	arguments = strings.TrimSpace(arguments)

	if arguments == "" {
		return TelegramResponse{}
	}

	var sourceLanguage, targetLanguage string

	switch update.Message.Command() {
	case "rus":
		sourceLanguage = "ru"
		targetLanguage = "eo"
	case "epo":
		sourceLanguage = "eo"
		targetLanguage = "ru"
	}

	msgText := translate.Translate(arguments, sourceLanguage, targetLanguage)

	return TelegramResponse{
		Method:           "sendMessage",
		ChatId:           update.Message.Chat.ID,
		ReplyToMessageID: update.Message.MessageID,
		Text:             msgText,
	}
}

func Handler(ctx context.Context, request *Request) ([]byte, error) {
	var update tgbotapi.Update

	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		log.Panic("Exception during unmarshal request body")
	}

	telegramResponse := handleUpdate(&update)

	response := Response{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            telegramResponse,
	}

	marshalResponse, err := json.Marshal(&response)
	if err != nil {
		log.Panic("Exception during marshal request body")
	}

	return marshalResponse, nil
}
