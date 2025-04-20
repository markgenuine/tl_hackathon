package parser

import (
	"github.com/markgenuine/tl_hackathon/internal/openai"
)

const (
	tagVersion            string = "1.03"
	tagCoding             string = "Windows"
	tagHeader             string = "header"
	tagDocument           string = "document"
	tagNameVersion        string = "ВерсияФормата"
	tagNameInvoice        string = "СекцияРасчСчет"
	tagNameSection        string = "СекцияДокумент"
	tagNameEndDocument    string = "КонецДокумента"
	tagNamePurposePayment string = "НазначениеПлатежа"
	tagEndFile            string = "КонецФайла"
)

// ParsedFile ...
type ParsedFile struct {
	Version   string            `json:"version"`
	Encoding  string            `json:"encoding"`
	Header    map[string]string `json:"header"`
	Documents []ParsedDocument  `json:"Документы"`
}

// ParsedDocument ...
type ParsedDocument struct {
	Type           string                 `json:"type"`
	Fields         map[string]string      `json:"fields"`
	PaymentDetails *openai.PaymentDetails `json:"ДетальнаяИнформация,omitempty"`
}
