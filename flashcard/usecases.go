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
	var builder TelegramTextBuilder

	builder.writeBold(flashcard.Name)

	if flashcard.HasExample() {
		builder.writeEmptyLine()
		builder.writeItalic(*flashcard.Example)
	}

	if flashcard.HasExplanation() {
		builder.writeEmptyLine()
		var explanationBuilder strings.Builder
		explanationBuilder.WriteString("                                          ") // added indent to make flash card wider
		explanationBuilder.WriteString("\n")
		explanationBuilder.WriteString(
			strings.TrimLeft(
				*flashcard.Explanation,
				"\n",
			),
		)
		builder.writeSpoiler(explanationBuilder.String())
	}

	builder.writeEmptyLine()
	builder.writeText("Choose: ")
	return builder.String()
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

type TelegramTextBuilder struct {
	strings.Builder
}

func (b *TelegramTextBuilder) writeBold(text string) {
	b.WriteString("*")
	b.writeText(text)
	b.WriteString("*")
}

func (b *TelegramTextBuilder) writeItalic(text string) {
	b.WriteString("_")
	b.writeText(text)
	b.WriteString("_")
}

func (b *TelegramTextBuilder) writeSpoiler(text string) {
	b.WriteString("||")
	b.writeText(text)
	b.WriteString("||")
}

func (b *TelegramTextBuilder) writeText(text string) {
	b.WriteString(shieldProhibitedSymbols(text))
}

func (b *TelegramTextBuilder) writeEmptyLine() {
	b.WriteString("\n")
	b.WriteString("\n")
}
