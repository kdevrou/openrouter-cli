# OpenRouter CLI

A command-line interface for [OpenRouter.ai](https://openrouter.ai) - access 400+ AI models directly from your terminal.

## Features

- **Chat completions**: Send prompts to any AI model on OpenRouter
- **List models**: Browse available models with pricing and capabilities
- **Flexible input**: Accept text as arguments or from stdin pipes
- **Multiple output formats**: Pretty-printed, raw, or JSON output
- **Easy configuration**: Store API key in config file or environment variable
- **Scriptable**: Perfect for piping to other commands

## Installation

### From Source

```bash
git clone https://github.com/kdevrou/openrouter-cli.git
cd openrouter-cli
make install
```

### Prerequisites

- Go 1.20 or later
- An OpenRouter API key (get one at https://openrouter.ai)

### NixOS Development

If you're on NixOS, a `flake.nix` is included for easy development setup:

```bash
# Enter the development environment
nix flake update  # First time only
nix develop

# Now you can build and use the project
make build
```

The flake provides Go and other build tools. The binary is statically compiled (no cgo), so it runs anywhere without external dependencies.

## Quick Start

### 1. Set up your API key

Choose one of these options (listed in order of precedence):

**Option A: Command-line flag (testing only)**
```bash
openrouter chat --api-key "sk-or-v1-..." "Your prompt"
```

**Option B: Environment variable (recommended for temporary use)**
```bash
export OPENROUTER_API_KEY="sk-or-v1-..."
openrouter chat "Your prompt"
```

**Option C: Config file (recommended for regular use)**

Create your config file (path varies by OS):
- **Linux**: `~/.config/openrouter/config.yaml`
- **macOS**: `~/Library/Application\ Support/openrouter/config.yaml` or `~/.config/openrouter/config.yaml`
- **Windows**: `%APPDATA%\openrouter\config.yaml`

```yaml
api_key: "sk-or-v1-..."
default_model: "openai/gpt-4"
output_format: "pretty"
```

**Security Note**: Set restrictive permissions on your config file:
```bash
chmod 600 ~/.config/openrouter/config.yaml
```
This prevents other users on the system from reading your API key.

**Config Precedence** (highest to lowest):
1. `--api-key` command-line flag
2. `OPENROUTER_API_KEY` environment variable
3. Config file
4. Defaults (if no API key configured, auth error)

### 2. Send a message

```bash
# Basic usage
openrouter chat "What is Go?"

# With a specific model
openrouter chat -m anthropic/claude-3.5-sonnet "Explain quantum computing"

# With pipes
echo "Summarize this text" | openrouter chat
```

### 3. List available models

```bash
# View all models
openrouter list

# Filter by name
openrouter list --filter gpt

# Get raw JSON
openrouter list --json
```

## Usage

### Chat Command

Send a message to an AI model:

```bash
openrouter chat [prompt] [flags]
```

**Examples:**

```bash
# Basic message
openrouter chat "Hello, world!"

# With specific model
openrouter chat -m "openai/gpt-4-turbo-preview" "Write a haiku about Go"

# With temperature control (creativity)
openrouter chat -t 1.5 "Generate a creative story"

# With token limit
openrouter chat --max-tokens 200 "Brief explanation of recursion"

# Raw output for piping
openrouter chat --raw "Generate a UUID" | cut -d' ' -f1

# JSON output for scripting
openrouter chat --json "Hello" | jq '.choices[0].message.content'
```

**Flags:**

- `-m, --model <model>` - Model to use (default: from config)
- `-t, --temperature <value>` - Temperature 0.0-2.0 (default: 1.0)
- `--max-tokens <n>` - Maximum tokens in response (default: 4096)
- `--raw` - Output only the response text (for piping)
- `--json` - Output full API response as JSON

**Input:**

Text can be provided as:
- Command argument: `openrouter chat "Your prompt here"`
- Piped input: `echo "Your prompt" | openrouter chat`
- Interactive (if neither provided, shows error with helpful message)

### List Command

Display available models:

```bash
openrouter list [flags]
```

**Examples:**

```bash
# Show all models
openrouter list

# Filter by search term
openrouter list --filter "gpt-4"
openrouter list --filter "claude"

# Get JSON output
openrouter list --json | jq '.[] | select(.context_length > 100000)'
```

**Flags:**

- `--filter <term>` - Filter models by name or ID
- `--json` - Output as JSON instead of table

### Global Flags

These flags work with all commands:

- `--api-key <key>` - Override API key (for quick testing)
- `--config <path>` - Use custom config file path
- `--debug` - Show debug information
- `-h, --help` - Show help
- `-v, --version` - Show version

## Configuration

### Config File Format

Default location: `~/.config/openrouter/config.yaml`

Supports XDG Base Directory spec - set `$XDG_CONFIG_HOME` to use a different location.

**Example config:**

```yaml
# Your OpenRouter API key
api_key: "sk-or-v1-..."

# Default model for chat (if not specified with -m)
default_model: "openai/gpt-4"

# Default temperature for chat requests
default_temperature: 1.0

# Default max tokens for responses
default_max_tokens: 4096

# Output format: pretty | raw | json
output_format: "pretty"

# API settings
api_base_url: "https://openrouter.ai/api/v1"
timeout: 60  # seconds - request timeout for API calls
```

### Environment Variables

- `OPENROUTER_API_KEY` - Your API key (highest priority)
- `XDG_CONFIG_HOME` - Custom config directory location

## Examples

### Simple Q&A

```bash
openrouter chat "What's the capital of France?"
openrouter chat -m anthropic/claude-3-opus "Explain quantum entanglement"
```

### Piping and Integration

```bash
# Get help from Claude
echo "How do I use grep with regex?" | openrouter chat --raw

# Generate and use data
openrouter chat --raw "Generate a random hex color code" | xargs -I {} echo "Color: {}"

# Combine with other tools
openrouter chat "List 5 JavaScript tips" --raw | grep -i "tip"

# Save responses to files
openrouter chat "Write a Python function to calculate fibonacci" > fibonacci.py
```

### Working with JSON

```bash
# Get specific fields from response
openrouter chat --json "Hello" | jq '.usage.total_tokens'

# View full response with all details
openrouter chat --json "Hello" | jq '.'

# Filter models
openrouter list --json | jq '.[] | select(.context_length > 100000) | .id'

# Extract model pricing
openrouter list --json | jq '.[] | {id, pricing}'

# See token usage and model details
openrouter chat --json "Your prompt" | jq '{model: .model, tokens: .usage.total_tokens}'
```

### Batch Processing

```bash
# Process multiple prompts
for prompt in "Hello" "Hi" "Hey"; do
  openrouter chat --raw "$prompt"
done

# With different models
openrouter chat -m "openai/gpt-4" "Technical question"
openrouter chat -m "openai/gpt-3.5-turbo" "Creative writing"
```

## Building

### Build the binary

```bash
make build
# Binary will be at ./bin/openrouter
```

### Install globally

```bash
make install
# Installs to $GOPATH/bin/openrouter
```

### Development build (with debug info)

```bash
go build -o bin/openrouter cmd/openrouter/main.go
./bin/openrouter chat "Hello" --debug
```

### NixOS/Flake Development

If you're developing on NixOS, use the provided Flake:

```bash
# Set up the development environment
nix develop

# Now you have all build tools available
make build
make install

# Exit the environment
exit
```

The flake provides Go, GCC, pkg-config, and other build tools automatically.

## Troubleshooting

### "No API key found"

Make sure your API key is set:

```bash
# Check environment variable
echo $OPENROUTER_API_KEY

# Or create config file at ~/.config/openrouter/config.yaml
cat ~/.config/openrouter/config.yaml
```

### "Model not found"

Check available models:

```bash
openrouter list --filter "gpt"
```

Use the full model ID from the output.

### Timeout errors

Try increasing the timeout in your config:

```yaml
timeout: 120  # seconds
```

Or disable piping to ensure the connection stays open.

### "No input provided"

Either pass a prompt as an argument:

```bash
openrouter chat "Your prompt"
```

Or pipe input:

```bash
echo "Your prompt" | openrouter chat
```

### "Provider returned error (HTTP 429)"

This means the model you're trying to use is rate-limited. This can happen if:

- The provider (OpenAI, Anthropic, etc.) is rate-limiting requests
- You've exceeded your account's rate limit
- The free tier for that model is overloaded

**Solutions:**
- Try a different model: `openrouter list` to see alternatives
- Wait a few minutes and try again
- Upgrade your account for higher rate limits
- Use a less popular model that has more availability

### "Provider returned error (HTTP 402)"

This means your OpenRouter account has a billing issue:

- **No credits**: Add credits at https://openrouter.ai/account/billing/overview
- **No payment method**: Add a valid credit card at https://openrouter.ai/account/billing/methods
- **Spending limit reached**: Increase it in your account settings

## Future Enhancements

Planned features for future releases:

- **Streaming responses** for real-time output with Server-Sent Events
- **Conversation history** with multi-turn conversations and session management
- **Image and video input** support (base64 encoding)
- **Interactive REPL mode** for multi-turn conversations without piping
- **Model aliases** (e.g., `gpt-4` → `openai/gpt-4-turbo-preview`)
- **Automatic retries** with exponential backoff for rate-limited requests (429 errors)
- **Cost estimation** before sending requests
- **Token counting** utilities to preview costs
- **Secret service integration** for secure API key storage (macOS Keychain, Linux Secret Service)
- **Shell command completion** generation for bash/zsh
- **Configuration management** commands (`openrouter config set/get`)

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.

## Support

- Issues: https://github.com/kdevrou/openrouter-cli/issues
- Discussions: https://github.com/kdevrou/openrouter-cli/discussions
- OpenRouter Docs: https://openrouter.ai/docs
