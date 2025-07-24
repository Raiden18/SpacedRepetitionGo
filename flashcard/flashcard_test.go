package flashcard

import (
	"testing"
)

var (
	TRUE  = true
	FALSE = false
)

func TestShouldRecallToLevel2(t *testing.T) {
	level1 := Flashcard{
		KnowLevel1:  &TRUE,
		KnowLevel2:  &FALSE,
		KnowLevel3:  &FALSE,
		KnowLevel4:  &FALSE,
		KnowLevel5:  &FALSE,
		KnowLevel6:  &FALSE,
		KnowLevel7:  &FALSE,
		KnowLevel8:  &FALSE,
		KnowLevel9:  &FALSE,
		KnowLevel10: &FALSE,
		KnowLevel11: &FALSE,
		KnowLevel12: &FALSE,
		KnowLevel13: &FALSE,
	}

	actual := level1.Recall()

	if *actual.KnowLevel1 != TRUE || *actual.KnowLevel2 != TRUE {
		t.Errorf(`Failed to recall to level 2.`)
	}
}

func TestShouldRecallToLevel13(t *testing.T) {
	level1 := Flashcard{
		KnowLevel1:  &TRUE,
		KnowLevel2:  &TRUE,
		KnowLevel3:  &TRUE,
		KnowLevel4:  &TRUE,
		KnowLevel5:  &TRUE,
		KnowLevel6:  &TRUE,
		KnowLevel7:  &TRUE,
		KnowLevel8:  &TRUE,
		KnowLevel9:  &TRUE,
		KnowLevel10: &TRUE,
		KnowLevel11: &TRUE,
		KnowLevel12: &TRUE,
		KnowLevel13: &FALSE,
	}

	actual := level1.Recall()

	if *actual.KnowLevel1 != TRUE || *actual.KnowLevel2 != TRUE || *actual.KnowLevel3 != TRUE || *actual.KnowLevel4 != TRUE || *actual.KnowLevel5 != TRUE || *actual.KnowLevel6 != TRUE || *actual.KnowLevel7 != TRUE || *actual.KnowLevel8 != TRUE || *actual.KnowLevel9 != TRUE || *actual.KnowLevel10 != TRUE || *actual.KnowLevel11 != TRUE || *actual.KnowLevel12 != TRUE || *actual.KnowLevel13 != TRUE {
		t.Errorf(`Failed to recall to level 13.`)
	}
}
