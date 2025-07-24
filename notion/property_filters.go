package notion

import (
	"strconv"

	notionApi "github.com/jomei/notionapi"
)

func And(filters ...notionApi.Filter) notionApi.AndCompoundFilter {
	return filters
}

func PropertyCheckBox(name string, equals bool) notionApi.PropertyFilter {
	var checkboxFilterCondition notionApi.CheckboxFilterCondition
	if equals {
		checkboxFilterCondition.Equals = true
	} else {
		checkboxFilterCondition.DoesNotEqual = true
	}
	return notionApi.PropertyFilter{
		Property: name,
		Checkbox: &checkboxFilterCondition,
	}
}

func KnowLevelName(level int) string {
	return "Know Level " + strconv.Itoa(level)
}

func KnowLevel(level int, equals bool) notionApi.PropertyFilter {
	return PropertyCheckBox(KnowLevelName(level), equals)
}

func Show(equals bool) notionApi.PropertyFilter {
	return PropertyCheckBox("Show", equals)
}
