import json
import sys
import argparse
from pathlib import Path

def decode(input_str):
    try:
        return json.loads(f'"{input_str}"')
    except json.JSONDecodeError:
        return input_str

def remove_comments(json_str):
    result = []
    lines = json_str.split('\n')
    in_multiline_comment = False
    
    for line in lines:
        if not line.strip():
            result.append(line)
            continue
            
        if in_multiline_comment:
            if '*/' in line:
                in_multiline_comment = False
                if '*/' in line:
                    result.append(line[line.find('*/')+2:])
            continue
            
        if '//' in line:
            line = line[:line.find('//')]
            
        if '/*' in line:
            in_multiline_comment = True
            if '/*' in line:
                result.append(line[:line.find('/*')])
            continue
            
        result.append(line)
    
    return '\n'.join(result)

def encode(input_str):
    result = ""
    i = 0
    in_string = False
    word_buffer = ""
    in_comment = False
    comment_type = None
    in_number = False
    
    while i < len(input_str):
        char = input_str[i]
        
        if not in_string:
            if char == '/' and i + 1 < len(input_str):
                next_char = input_str[i + 1]
                if next_char == '/':
                    in_comment = True
                    comment_type = '//'
                    i += 2
                    continue
                elif next_char == '*':
                    in_comment = True
                    comment_type = '/*'
                    i += 2
                    continue
            
            if in_comment:
                if comment_type == '//' and char == '\n':
                    in_comment = False
                    comment_type = None
                    i += 1
                    continue
                elif comment_type == '/*' and char == '*' and i + 1 < len(input_str) and input_str[i + 1] == '/':
                    in_comment = False
                    comment_type = None
                    i += 2
                    continue
                i += 1
                continue
        
        if char == '"' and (i == 0 or input_str[i-1] != '\\'):
            in_string = not in_string
            word_buffer = ""
            in_number = False
            result += char
            i += 1
            continue
            
        if in_string:
            if char == '\\' and i + 1 < len(input_str) and input_str[i + 1] in '"\\/bfnrt':
                result += char + input_str[i + 1]
                i += 2
                continue
            result += f"\\u{hex(ord(char))[2:].zfill(4)}"
            i += 1
            continue

        if char in '{}[],: ' or char.isspace():
            word_buffer = ""
            in_number = False
            result += char
            i += 1
            continue

        # Handle negative numbers
        if char == '-' and i + 1 < len(input_str) and input_str[i + 1] in '0123456789':
            in_number = True
            result += char
            i += 1
            continue

        if char in '0123456789.':
            if not in_number and char in '0123456789':
                in_number = True
            if in_number:
                result += char
                i += 1
                continue

        if char.isalpha():
            word_buffer += char
            if word_buffer in ['true', 'false', 'null']:
                result += word_buffer
                word_buffer = ""
                i += 1
                continue
            if len(word_buffer) > 4:
                result += word_buffer
                word_buffer = ""
                i += 1
                continue
            i += 1
            continue

        result += f"\\u{hex(ord(char))[2:].zfill(4)}"
        i += 1
        
    return result

def process(input_path, output_path, mode='decode'):
    try:
        with open(input_path, 'r', encoding='utf-8') as file:
            if mode == 'decode':
                content = file.read()
                clean_json = remove_comments(content)
                data = json.loads(clean_json)
            else:
                content = file.read()
                clean_json = remove_comments(content)
                data = clean_json

        with open(output_path, 'w', encoding='utf-8') as file:
            if mode == 'decode':
                json.dump(data, file, indent=2, ensure_ascii=False)
            else:
                file.write(encode(data))

        print(f"Output: {output_path}")

    except Exception as e:
        print(f"Error: {str(e)}")
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description='Unicode Escape Sequence Encoder/Decoder')
    parser.add_argument('mode', choices=['encode', 'decode'], help='Operation mode: encode/decode')
    parser.add_argument('input_path', type=str, help='Input file path')
    parser.add_argument('output_path', type=str, help='Output file path')

    args = parser.parse_args()

    if not Path(args.input_path).exists():
        print(f"Error: '{args.input_path}' not found")
        sys.exit(1)

    process(args.input_path, args.output_path, args.mode)

if __name__ == "__main__":
    main()