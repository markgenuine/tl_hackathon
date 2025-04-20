package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/markgenuine/tl_hackathon/internal/openai"
	"github.com/markgenuine/tl_hackathon/internal/parser"
)

// UploadHandler ...
func UploadHandler(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(60 << 20) // file limit 60 mb
		if err != nil {
			http.Error(w, "Error of parsing form, big size file: "+err.Error(), http.StatusBadRequest)
			return
		}

		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			http.Error(w, "File is not select", http.StatusBadRequest)
			return
		}

		var parsedFiles []*parser.ParsedFile

		for _, fh := range files {
			file, err := fh.Open()
			if err != nil {
				http.Error(w, "Error of open file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			parsed, purposePayments, err := parser.Parse1CClientBankExchange(file)
			if err != nil {
				http.Error(w, "Error parsing file: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if err := parser.ValidateParsedFile(parsed); err != nil {
				http.Error(w, "Error file validation: "+err.Error(), http.StatusBadRequest)
				return
			}

			result, err := openai.AnalyzePurposePayment(purposePayments, config.APIURL, config.APIKey)
			if err == nil && len(result) == len(parsed.Documents) {
				for i := range parsed.Documents {
					if result[i] == nil {
						continue
					}

					parsed.Documents[i].PaymentDetails = result[i]
				}
			}

			parsedFiles = append(parsedFiles, parsed)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(parsedFiles)
	}
}
