package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	Notion struct {
		ApiKey                   string `yaml:"api_key"`
		ObservedDatabases        string `yaml:"observed_databases"`
		EnglishVocabularyId      string `yaml:"english_vocabulary_id"`
		GreekVocabularyId        string `yaml:"greek_vocabulary_id"`
		GreekPassiveVocabularyId string `yaml:"greek_passive_vocabulary_id"`
	} `yaml:"notion"`

	Telegram struct {
		ApiKey string `yaml:"api_key"`
		ChatId string `yaml:"chat_id"`
	} `yaml:"telegram"`
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
	id, err := strconv.ParseInt(file.Telegram.ChatId, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return id
}

var (
	configOnce sync.Once
	configData Conf
)

func NewConfigFromFile() Conf {
	configOnce.Do(func() {
		configPath := secretsPath()
		bytes, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("read config file %s: %v", configPath, err)
		}

		if err := yaml.Unmarshal(bytes, &configData); err != nil {
			log.Fatalf("parse config file %s: %v", configPath, err)
		}
	})

	return configData
}

func secretsPath() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("resolve config path: runtime caller unavailable")
	}

	return filepath.Join(filepath.Dir(currentFile), "..", "secrets.yaml")
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

func GetObservedDatabasesId() string {
	return NewConfigFromFile().Notion.ObservedDatabases
}
