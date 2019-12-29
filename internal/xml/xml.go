package xml

import (
	"bytes"
	"encoding/xml"
)

type tagPair struct {
	Start string
	End   string
}

var TAGS = map[string]tagPair{"paragraph": tagPair{"<p>", "</p>"}, "bold": tagPair{"<b>", "</b>"}, "body": tagPair{"<div class=\"bbc-body\">", "</div>"}, "crosshead": tagPair{"<h2 class=\"bbc-crosshead\">", "</h2>"}, "list-unordered": tagPair{"<ul>", "</ul>"}, "list-ordered": tagPair{"<ol>", "</ol>"}, "listItem": tagPair{"<li>", "</li>"}}

func XmlToHtml(xmlStr string) string {
	xmlBuffer := bytes.NewBufferString(xmlStr)
	decoder := xml.NewDecoder(xmlBuffer)

	lessons := ""
	listType := "unordered"
	omit := false

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch v := token.(type) {
		case xml.StartElement:
			name := v.Name.Local
			if name == "list" {
				for _, attr := range v.Attr {
					if attr.Name.Local == "type" {
						listType = attr.Value
						break
					}
				}
				name = name + "-" + listType
			}
			tag, ok := TAGS[name]
			if ok {
				omit = false
				lessons += tag.Start
			} else {
				omit = true
			}
		case xml.EndElement:
			name := v.Name.Local
			if name == "list" {
				name = name + "-" + listType
			}
			tag, ok := TAGS[name]
			if ok {
				lessons += tag.End
			}
		case xml.CharData:
			if !omit {
				lessons += string(v)
			}
		}
	}

	return lessons
}
