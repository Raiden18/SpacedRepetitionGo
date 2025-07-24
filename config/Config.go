package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Notion   Notion   `json:"notion"`
	Telegram Telegram `json:"telegram"`
}

type Notion struct {
	ApiKey                   string   `json:"api_key"`
	ObservedDatabases        []string `json:"observed_databases"`
	EnglishVocabularyId      string   `json:"english_vocabulary_id"`
	GreekVocabularyId        string   `json:"greek_vocabulary_id"`
	GreekPassiveVocabularyId string   `json:"greek_passive_vocabulary_id"`
}

type Telegram struct {
	ApiKey string `json:"api_key"`
	ChatId string `json:"chat_id"`
}

type NotionDataBase struct {
	Id string `json:"id"`
}

func NotionApiKey() string {
	return NewConfigFromFile().Notion.ApiKey
}

func EnglishVocabularyId() string {
	return NewConfigFromFile().Notion.EnglishVocabularyId
}

func GreekVocabularyId() string {
	return NewConfigFromFile().Notion.GreekVocabularyId
}

func GreekPassiveVocabularyId() string {
	return NewConfigFromFile().Notion.GreekPassiveVocabularyId
}

func TelegramApiKey() string {
	return NewConfigFromFile().Telegram.ApiKey
}

func TelegramChatId() int64 {
	file := NewConfigFromFile()
	id, error := strconv.ParseInt(file.Telegram.ChatId, 10, 64)
	if error != nil {
		log.Fatalln(error)
	}
	return id
}

func NewConfigFromFile() Config {
	homeDir, _ := os.UserHomeDir()
	configPath := homeDir + "/repetition/spaced_repetition.config"
	configFile, _ := os.Open(configPath)
	bytes, _ := io.ReadAll(configFile)

	var config Config

	err := json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}
	configFile.Close()
	return config
}

func NewNotionDataBase(id string) NotionDataBase {
	empty := []string{}
	return NewNotionDataBaseWithDictionary(id, empty)
}

func NewNotionDataBaseWithDictionary(id string, dictionaries []string) NotionDataBase {
	return NotionDataBase{
		Id: id,
	}
}

func GetObservedDatabasesIds() []string {
	return NewConfigFromFile().Notion.ObservedDatabases
}
