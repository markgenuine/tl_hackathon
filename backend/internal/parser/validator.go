package parser

import (
	"fmt"
)

var requiredHeaderFields = []string{
	"ВерсияФормата",
	"ДатаНачала",
	"ДатаКонца",
	"РасчСчет",
}

var requiredDocFields = []string{
	"Номер",
	"Дата",
	"Сумма",
	"Плательщик",
	"Получатель",
	"НазначениеПлатежа",
}

// ValidateParsedFile checks that required fields are present in header and documents.
func ValidateParsedFile(p *ParsedFile) error {
	missingHeader := []string{}
	for _, field := range requiredHeaderFields {
		if _, ok := p.Header[field]; !ok {
			missingHeader = append(missingHeader, field)
		}
	}
	if len(missingHeader) > 0 {
		return fmt.Errorf("missing required fields in header: %v", missingHeader)
	}

	for i, doc := range p.Documents {
		missingDoc := []string{}
		for _, field := range requiredDocFields {
			if _, ok := doc.Fields[field]; !ok {
				missingDoc = append(missingDoc, field)
			}
		}
		if len(missingDoc) > 0 {
			return fmt.Errorf("document #%d (%s): missing fields: %v", i+1, doc.Type, missingDoc)
		}
	}

	return nil
}
