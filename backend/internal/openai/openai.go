package openai

const (
	batchSize = 50

	// PromptParams ...
	PromptParams = `
1. НомерСчета: Нужно извлечь номер счета покупателя, если он не является банковским счетом. Если указан только банковский счет, то вернуть пустую строку.
2. ДатаСчета: Это дата, связанная с параметром "Номер счета покупателя". Если "Номер счета покупателя" пуст, то и "Дата счета покупателя" должна быть пустой.
3. НомерДоговора: Номер договора указан после слова "ДБС".
4. ДатаДоговора: Дата договора должна быть извлечена только в том случае, если найден номер договора. Дата договора указывается после номера договора (обычно после слова "от"). Если номер договора не найден, вернуть пустую строку.
5. СтавкаНДС: Если в строке указано "НДС не обл." или "не облагается", значение ставки НДС должно быть "Без НДС". Если ставка НДС не указана, вернуть пустую строку.
Если любой параметр отсутствует, вернуть пустую строку. Ответ должен быть строго в том же порядке и количестве, как и исходные строки с номером.
Результат вывести в формате JSON.
Строки:
`
)

// PaymentDetails ...
type PaymentDetails struct {
	AccountNumber  string `json:"НомерСчета,omitempty"`
	InvoiceDate    string `json:"ДатаСчета,omitempty"`
	ContractNumber string `json:"НомерДоговора,omitempty"`
	ContractDate   string `json:"ДатаДоговора,omitempty"`
	VATRate        string `json:"СтавкаНДС,omitempty"`
}

// Message ...
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest ...
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// ChatResponse ...
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
