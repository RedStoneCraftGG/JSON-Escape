# JSON-Escape
A Unicode escape sequence encoder/decoder for JSON files. Compatible with Minecraft Bedrock JSON UI.

# Usage

## Python

```bash
py main.py <encode/decode> <input> <output>
```

Example

```bash
py main.py encode file_name.json output_name.json
```

## GoLang

```bash
go run main.go -m <encode/decode> -i <input> -o <output>
```

Example

```bash
go run main.go -m e -i file_name.json -o output_name.json
```

Known Issue (GoLang only): During encoding or decoding, the JSON data will be sorted alphabetically. However, all functions will remain fully operational.

## HTML (Web)

Simply open the HTML file in your browser.
