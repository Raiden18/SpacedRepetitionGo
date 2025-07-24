package utils

import (
	notionApi "github.com/jomei/notionapi"
)

func RichTextToString(richText []notionApi.RichText) string {
	var content string = ""
	for _, rt := range richText {
		if rt.Text != nil {
			content += rt.Text.Content
		}
	}
	return content
}
