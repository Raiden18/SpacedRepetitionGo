package flashcard

import (
	"spacedrepetitiongo/telegram"
	"strings"
)

func SendToTelegram(bot telegram.Bot, message FlashcardTelegramMessage) {
	flashcard := message.GetFlashCard()
	text := asTextMessage(flashcard)
	if flashcard.HasImage() {
		bot.SendPhotoMessage(
			text,
			*flashcard.Image,
			message.GetButtons(),
		)
	} else {
		bot.SendTextMessage(
			text,
			message.GetButtons(),
		)
	}
}

func asTextMessage(flashcard Flashcard) string {
	var textBuffer strings.Builder

	textBuffer.WriteString("*" + shieldProhibitedSymbols(flashcard.Name) + "*")

	if flashcard.HasExample() {
		textBuffer.WriteString("\n")
		textBuffer.WriteString("\n")
		textBuffer.WriteString("_" + shieldProhibitedSymbols(*flashcard.Example) + "_")
	}

	if flashcard.HasExplanation() {
		textBuffer.WriteString("\n")
		textBuffer.WriteString("\n")
		textBuffer.WriteString("||" + shieldProhibitedSymbols(strings.TrimLeft(*flashcard.Explanation, "\n")) + "||")
	}

	textBuffer.WriteString("\n")
	textBuffer.WriteString("\n")
	textBuffer.WriteString("Choose: ")
	return textBuffer.String()
}

func shieldProhibitedSymbols(from string) string {
	replacer := strings.NewReplacer(
		"{", "\\{",
		"}", "\\}",
		"|", "\\|",
		"#", "\\#",
		"<", "\\<",
		">", "\\>",
		"`", "\\`",
		"~", "\\~",
		"[", "\\[",
		"]", "\\]",
		"*", "\\*",
		"!", "\\!",
		"(", "\\(",
		")", "\\)",
		"=", "\\=",
		".", "\\.",
		"_", "\\_",
		"-", "\\-",
		"+", "\\+",
		"\\", "\\\\",
	)
	return replacer.Replace(from)
}
