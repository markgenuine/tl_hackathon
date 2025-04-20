package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	model       = "unsloth/gemma-3-27b-it"
	temperature = 0.2
	maxTokens   = 10000
	userText    = "user"
	startPrompt = "Извлеки из строк следующие параметры:"
)

// AnalyzePurposePayment ...
func AnalyzePurposePayment(purposePayments []string, apiUrl, apiKey string) ([]*PaymentDetails, error) {
	allResults := make([]*PaymentDetails, 0, len(purposePayments))

	for i := 0; i < len(purposePayments); i += batchSize {
		end := i + batchSize
		if end > len(purposePayments) {
			end = len(purposePayments)
		}
		chunk := purposePayments[i:end]

		results, err := analyzeChunk(chunk, apiUrl, apiKey)
		if err != nil {
			return nil, fmt.Errorf("error of parsing chank %d-%d: %w", i, end, err)
		}
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

func analyzeChunk(chunk []string, apiUrl, apiKey string) ([]*PaymentDetails, error) {
	var joined strings.Builder
	for _, txt := range chunk {
		joined.WriteString(fmt.Sprintf("%s\n", txt))
	}

	prompt := startPrompt + PromptParams + joined.String()

	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{Role: userText, Content: prompt},
		},
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return nil, fmt.Errorf("error of decode chatResp: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, errors.New("empty response from model")
	}

	content := chatResp.Choices[0].Message.Content

	start := strings.Index(content, "[")
	end := strings.LastIndex(content, "]")
	if start >= 0 && end >= 0 && end > start {
		content = content[start : end+1]
	}

	var result []*PaymentDetails
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("error of parsing to JSON from content field: %w\nraw: %s", err, content)
	}

	return result, nil
}
