# 📦 tl_hackathon

`tl_hackathon` — project for the tl_hackathon.
It is a web application built using the Go programming language.

## 🧩 Project struct

The project has the following structure:

- `backend/` — source code of the service for processing payment assignments from 1C format.
- `local_ai` — attempt to build your own AI for local use for string recognition with training.
- `1c` — client part of 1C (implemented as an extension).

## 🚀 Launch

### 🛠️ backend:
To build and run the project, you need to have the Go compiler  version 1.24+ installed.

1. Clone repository:

   ```bash
   git clone https://github.com/markgenuine/tl_hackathon.git
   cd tl_hackathon
    ```

2. Execute command:

    ``` bash
    cd backend
    go build
    ```

3. Execute application:

    ``` bash
    ./backend
    ```

### 🧠 local_AI (experement):
It's try launch local ai from opensource project [GPT4All](https://github.com/nomic-ai/gpt4all), model: `Qwen/Qwen2-1.5B-Instruct` with fun_tuning from generate `*.gguf-file`.

#### 🚀 What is used

- **GPT4All** — a platform with a convenient API and GUI for local launch of LLM.
- **Qwen2-1.5B-Instruct** — a compact yet powerful model from Alibaba, optimized for instructions.
- **GGUF** — the model format supported by `llama.cpp` is optimized for performance.
- **Fine-tuning** — additional training of the model on user data.
- **Interface** — connecting the model to your application via GPT4All API or gRPC/WebSocket.

#### 📦 Preparation

1. **Download GPT4All**:
   - [https://github.com/nomic-ai/gpt4all](https://github.com/nomic-ai/gpt4all)

2. **Training model**:
    - [https://github.com/markgenuine/tl_hackathon/local_ai/training.py](Model training)

3. **Convert model in `.gguf` file**:
   If you did fine-tuning in Hugging Face format, use `transformers` + `gguf` tools:

   Use lib: [llama.cpp](https://github.com/ggml-org/llama.cpp)
  
   ```bash
   python convert_huggingface_to_gguf.py --input_dir ./path_to_model --output_file model.gguf
   ```

4. **Start model**:

    ```bash
    ./gpt4all --model models/qwen2_finetuned.gguf
    ```

5. **Example, how use model**:

    ```bash
    curl -X POST http://localhost:4891/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d '{
     "model": "Qwen2-1.5B-Instruct",
     "messages": [{"role": "user", "content": "Hello, how are you?"}]
    }'
    ```
