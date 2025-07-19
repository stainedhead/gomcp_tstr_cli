# Provider Interface

## Class Structure

```
┌───────────────────────────────┐
│        <<interface>>          │
│          Provider             │
├───────────────────────────────┤
│ + Name(): string              │
│ + Chat(): (*ChatResponse, err)│
│ + ChatStream(): (<-chan, err) │
│ + ValidateConfig(): error     │
│ + Close(): error              │
└───────────────┬───────────────┘
                │
                │
┌───────────────┴───────────────┐
│       OllamaProvider          │
├───────────────────────────────┤
│ - endpoint: string            │
│ - model: string               │
│ - client: *http.Client        │
│ - logger: Logger              │
├───────────────────────────────┤
│ + NewOllamaProvider()         │
│ + Name(): string              │
│ + ValidateConfig(): error     │
│ + Chat(): (*ChatResponse, err)│
│ + ChatStream(): (<-chan, err) │
│ + convertRequest()            │
│ + Close(): error              │
└───────────────────────────────┘
```

## Implemented Providers
- **Ollama**: Local model execution
- **Planned**: OpenAI, AWS Bedrock, Google AI, Anthropic

---
