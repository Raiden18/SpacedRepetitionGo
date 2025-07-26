package flashcard

import (
	"testing"
)

var (
	TRUE  = true
	FALSE = false
)

func TestShouldForgetAsMap(t *testing.T) {
	knowLevels := map[int]*bool{
		1:  &TRUE,
		2:  &TRUE,
		3:  &TRUE,
		4:  &TRUE,
		5:  &TRUE,
		6:  &TRUE,
		7:  &TRUE,
		8:  &TRUE,
		9:  &TRUE,
		10: &TRUE,
	}

	ForgetAsMap(knowLevels)

	expected := map[int]*bool{
		1:  &FALSE,
		2:  &FALSE,
		3:  &FALSE,
		4:  &FALSE,
		5:  &FALSE,
		6:  &FALSE,
		7:  &FALSE,
		8:  &FALSE,
		9:  &FALSE,
		10: &FALSE,
	}

	for i := range knowLevels {
		if *knowLevels[i] != *expected[i] {
			t.Errorf(`Failed to forget level %d. Expected %v, got %v`, i, *expected[i], *knowLevels[i])
			return
		}
	}
}

func TestShoulMemrozieAsMap(t *testing.T) {
	knowLevels := map[int]*bool{
		1: &FALSE,
		2: &FALSE,
		3: &FALSE,
		4: &FALSE,
		5: &FALSE,
	}

	MemorizeAsMap(knowLevels)

	expected := map[int]*bool{
		1: &TRUE,
		2: &FALSE,
		3: &FALSE,
		4: &FALSE,
		5: &FALSE,
	}

	for i := range knowLevels {
		if *knowLevels[i] != *expected[i] {
			t.Errorf(`Failed to memorize level %d. Expected %v, got %v`, i, *expected[i], *knowLevels[i])
			return
		}
	}
}

func TestShouldRecallFromLevel1To3AsMap(t *testing.T) {
	knowLevels := map[int]*bool{
		1: &TRUE,
		2: &TRUE,
		3: &FALSE,
		4: &FALSE,
		5: &FALSE,
		6: &FALSE,
		7: &FALSE,
	}

	RecallAsMap(knowLevels)

	expected := map[int]*bool{
		1: &TRUE,
		2: &TRUE,
		3: &TRUE,
		4: &FALSE,
		5: &FALSE,
		6: &FALSE,
		7: &FALSE,
	}

	for i := range knowLevels {
		if *knowLevels[i] != *expected[i] {
			t.Errorf(`Failed to recall level %d. Expected %v, got %v`, i, *expected[i], *knowLevels[i])
			return
		}
	}
}
