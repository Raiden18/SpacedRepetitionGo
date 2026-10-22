package flashcard

import (
	"sort"
	"spacedrepetitiongo/telegram"
)

var (
	FORGOTTEN_VALUE = false
	MEMORIZED_VALUE = true
)

type Flashcard struct {
	Id string
	// Image is what the bot sends: a local file path once the image job has
	// downloaded it, or the remote URL until then.
	Image *string
	// ImageSource is the cover URL as it currently is in Notion.
	ImageSource *string
	BoxId       string
	Name        string
	Example     *string
	Explanation *string
	// Queue is QueueRevise or QueueMemorize; Position orders the box's cards
	// within the queue.
	Queue    string
	Position int
	// Previous and Next are the neighbouring cards' ids, filled in when the
	// card is loaded for display.
	Next       *string
	Previous   *string
	KnowLevels map[int]*bool
}

func (flashcard Flashcard) RemoveFromChat(bot telegram.Bot, id int) Flashcard {
	bot.DeleteMessage(id)
	return flashcard
}

// SendToTelegram sends the card with the buttons for its queue. A nil card
// sends nothing.
func (flashcard *Flashcard) SendToTelegram(bot telegram.Bot) {
	if flashcard == nil {
		return
	}
	if flashcard.Queue == QueueMemorize {
		flashcard.ToTelegramMessageToMemorize().SendToTelegram(bot)
	} else {
		flashcard.ToTelegramMessageToRevise().SendToTelegram(bot)
	}
}

func (flashcard *Flashcard) ToTelegramMessageToRevise() *FlashcardTelegramMessageToRevise {
	if flashcard == nil {
		return nil
	}
	message := NewFlashcardTelegramMessageToRevise(*flashcard)
	return &message
}

func (flashcard *Flashcard) ToTelegramMessageToMemorize() *FlashcardTelegramMessageToMemorize {
	if flashcard == nil {
		return nil
	}
	message := NewFlashcardTelegramMessageToMemorize(*flashcard)
	return &message
}

func (flashcard Flashcard) HasExplanation() bool {
	return flashcard.Explanation != nil && *flashcard.Explanation != ""
}

func (flashcard Flashcard) HasExample() bool {
	return flashcard.Example != nil && *flashcard.Example != ""
}

func (flashcard Flashcard) HasImage() bool {
	return flashcard.Image != nil && *flashcard.Image != ""
}

func HasImage(flashcard Flashcard) bool {
	return flashcard.HasImage()
}

func (flashcard Flashcard) HasImageSource() bool {
	return flashcard.ImageSource != nil && *flashcard.ImageSource != ""
}

func (flashcard Flashcard) Recall() Flashcard {
	RecallAsMap(flashcard.KnowLevels)
	return flashcard
}

func (flashcard Flashcard) Forget() Flashcard {
	ForgetAsMap(flashcard.KnowLevels)
	return flashcard
}

func (flashcard Flashcard) Memorize() Flashcard {
	MemorizeAsMap(flashcard.KnowLevels)
	return flashcard
}

func ForgetAsMap(knowLevels map[int]*bool) {
	for key, value := range knowLevels {
		if value != nil {
			knowLevels[key] = &FORGOTTEN_VALUE
		}
	}
}

func MemorizeAsMap(knowLevels map[int]*bool) {
	ForgetAsMap(knowLevels)
	knowLevels[1] = &MEMORIZED_VALUE
}

func RecallAsMap(knowLevels map[int]*bool) {

	keys := make([]int, 0, len(knowLevels))
	for k := range knowLevels {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		value := knowLevels[key]
		if value != nil && !*value {
			knowLevels[key] = &MEMORIZED_VALUE
			break
		}
	}
}
