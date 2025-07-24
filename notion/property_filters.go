package notion

import (
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
