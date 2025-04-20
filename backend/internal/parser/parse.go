package parser

import (
	"bufio"
	"io"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

// Parse1CClientBankExchange parses the 1CClientBankExchange file format v1.03
func Parse1CClientBankExchange(r io.Reader) (*ParsedFile, []string, error) {
	decoder := charmap.Windows1251.NewDecoder()
	utf8Reader := decoder.Reader(r)
	scanner := bufio.NewScanner(utf8Reader)

	parsed := &ParsedFile{
		Version:   tagVersion,
		Encoding:  tagCoding,
		Header:    make(map[string]string),
		Documents: make([]ParsedDocument, 0),
	}

	var purposePayments []string
	var currentDoc *ParsedDocument
	mode := tagHeader

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, tagNameVersion+"="):
			parsed.Version = strings.TrimPrefix(line, tagNameVersion+"=")
			parsed.Header[tagNameVersion] = parsed.Version

		case line == tagNameInvoice:
			mode = tagHeader

		case strings.HasPrefix(line, tagNameSection+"="):
			mode = tagDocument
			docType := strings.TrimPrefix(line, tagNameSection+"=")
			currentDoc = &ParsedDocument{
				Type:   docType,
				Fields: make(map[string]string),
			}

		case line == tagNameEndDocument:
			if currentDoc != nil {
				parsed.Documents = append(parsed.Documents, *currentDoc)
				if val, exists := currentDoc.Fields[tagNamePurposePayment]; exists {
					purposePayments = append(purposePayments, val)
				} else {
					purposePayments = append(purposePayments, "")
				}

				currentDoc = nil
			}

		case line == tagEndFile:

		default:
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key, val := parts[0], parts[1]

			switch mode {
			case tagHeader:
				parsed.Header[key] = val
			case tagDocument:
				if currentDoc != nil {
					currentDoc.Fields[key] = val
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return parsed, purposePayments, nil
}
