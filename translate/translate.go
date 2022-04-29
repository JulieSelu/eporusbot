package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type RequestTranslate struct {
	FolderId           string   `json:"folderId"`
	SourceLanguageCode string   `json:"sourceLanguageCode"`
	TargetLanguageCode string   `json:"targetLanguageCode"`
	Texts              []string `json:"texts"`
}

type TextTranslation struct {
	Text string `json:"text"`
}

type TranslationsStruct struct {
	Translations []TextTranslation `json:"translations"`
}

func Translate(text string, sourceLanguage string, targetLanguage string) string {
	url := os.Getenv("YandexTranslateUrl")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	data := &RequestTranslate{
		FolderId:           os.Getenv("YandexCloudFolderId"),
		SourceLanguageCode: sourceLanguage,
		TargetLanguageCode: targetLanguage,
		Texts:              []string{text},
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Panic("Exception during convert data to json")
	}
	log.Print(jsonData)
	log.Print(string(jsonData))

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Panic("Exception during create request")
	}

	yandexAuthorizationHead := fmt.Sprintf("Api-Key %s", os.Getenv("YandexCloudApiKey"))
	request.Header.Set("Authorization", yandexAuthorizationHead)

	response, err := client.Do(request)
	if err != nil {
		log.Panic("Error during translate")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	log.Print(string(body))

	var v TranslationsStruct

	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Panic("Error")
	}

	return v.Translations[0].Text
}
