package flashcard

import (
	"fmt"
	"spacedrepetitiongo/image"
	"spacedrepetitiongo/openai"
	"spacedrepetitiongo/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FlashcardTelegramMessageGptGennerated struct {
	Flashcard Flashcard
}

func GenerateFromGPT(
	language string,
	word string,
	openAi *openai.OpenAiClient,
) FlashcardTelegramMessageGptGennerated {
	englishMeaning := openAi.Ask(
		"You are a Greek language tutor",
		fmt.Sprintf(
			"Give english translation for the word \"%s\""+
				"The english must be simple, natural, and suitable"+
				"Do not add any text before or after the translation."+
				"Verbs must be in nominative for \"I\", nouns in singular form",
			word,
		),
	)

	urlImage := openAi.CreateImage(
		fmt.Sprintf(
			"Create a realistic, concrete, and simple image that visually represents the phrase \"%s\". Avoid abstract concepts. The image should be clear, direct, and easy to understand, suitable for language learners.",
			englishMeaning,
		),
	)

	image.ConvertBase64ToImage(
		urlImage,
		"/root/repetition/images/",
		"gpt_generated",
	)

	url, _ := image.FindFileByNameWithoutExt(
		"/root/repetition/images/",
		"gpt_generated",
	)
	return FlashcardTelegramMessageGptGennerated{
		Flashcard: Flashcard{
			Id:          "GPT_GENERATED",
			Image:       &url,
			BoxId:       "NO",
			Name:        englishMeaning,
			Example:     nil,
			Explanation: &word,
			KnowLevels:  make(map[int]*bool),
		},
	}
}

func (message FlashcardTelegramMessageGptGennerated) GetButtons() *tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			saveToNotion(message.Flashcard),
		),
	}
	buttons := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &buttons
}

func (message FlashcardTelegramMessageGptGennerated) GetFlashCard() Flashcard {
	return message.Flashcard
}

func (message *FlashcardTelegramMessageGptGennerated) SendToTelegram(bot telegram.Bot) {
	if message == nil {
		return
	}
	SendToTelegram(bot, message)
}

func saveToNotion(flashcard Flashcard) tgbotapi.InlineKeyboardButton {
	return telegram.NewCallbackButton(
		"Add to notion ✅",
		Parameter(
			SaveToNotiondKey(),
			flashcard.Id,
		),
	)
}

func SaveToNotiondKey() string {
	return "save_to_notion"
}
