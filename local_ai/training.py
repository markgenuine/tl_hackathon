import json
import random
from datasets import Dataset
from transformers import AutoModelForCausalLM, AutoTokenizer, Trainer, TrainingArguments, DataCollatorForSeq2Seq
from peft import LoraConfig, get_peft_model
import numpy as np

with open("dataset.jsonl", "r", encoding="utf-8") as f:
    raw_data = [json.loads(line) for line in f]

random.shuffle(raw_data)

split_idx = int(0.8 * len(raw_data))
train_data = raw_data[:split_idx]

train_dataset = Dataset.from_list(train_data)

model_name = "Qwen/Qwen2-1.5B-Instruct"
tokenizer = AutoTokenizer.from_pretrained(model_name, trust_remote_code=True)
model = AutoModelForCausalLM.from_pretrained(model_name, trust_remote_code=True, device_map="auto")

lora_config = LoraConfig(
    r=8,
    lora_alpha=32,
    lora_dropout=0.05,
    bias="none",
    task_type="CAUSAL_LM"
)
model = get_peft_model(model, lora_config)

def tokenize_function(example):
    user_msg = f"{example['instruction']} {example['input']}".strip()
    assistant_msg = example["output"].strip()

    full_prompt = (
        f"<|im_start|>user\n{user_msg}<|im_end|>\n"
        f"<|im_start|>assistant\n{assistant_msg}<|im_end|>"
    )

    tokenized = tokenizer(
        full_prompt,
        truncation=True,
        max_length=1024,
        padding="max_length",
        return_tensors="pt"
    )

    tokenized = {k: v.squeeze(0) for k, v in tokenized.items()}
    tokenized["labels"] = tokenized["input_ids"].clone()

    return tokenized

tokenized_train = train_dataset.map(tokenize_function)

data_collator = DataCollatorForSeq2Seq(tokenizer, model=model)

training_args = TrainingArguments(
    output_dir="./qwen2_fine_tuned",
    per_device_train_batch_size=4,
    gradient_accumulation_steps=4,
    learning_rate=2e-4,
    num_train_epochs=3,
    logging_dir="./logs",
    logging_steps=10,
    save_steps=100,
    fp16=True, 
)

trainer = Trainer(
    model=model,
    args=training_args,
    train_dataset=tokenized_train,
    tokenizer=tokenizer,
    data_collator=data_collator
)

trainer.train()

trainer.save_model("./qwen2_fine_tuned")
tokenizer.save_pretrained("./qwen2_fine_tuned") 

merged_model = model.merge_and_unload()
merged_model.save_pretrained("./qwen2_fine_tuned")