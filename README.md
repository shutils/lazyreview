# lazyreview
Terminal UI for code review with ai

## Description

`lazyreview` is a code review tool that uses GPT-4o to generate code reviews.

## Image

![image](https://private-user-images.githubusercontent.com/100141359/405655010-246c98b7-f4fa-42b3-b3e2-06830aa42539.PNG?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3Mzc1NTU5NzMsIm5iZiI6MTczNzU1NTY3MywicGF0aCI6Ii8xMDAxNDEzNTkvNDA1NjU1MDEwLTI0NmM5OGI3LWY0ZmEtNDJiMy1iM2UyLTA2ODMwYWE0MjUzOS5QTkc_WC1BbXotQWxnb3JpdGhtPUFXUzQtSE1BQy1TSEEyNTYmWC1BbXotQ3JlZGVudGlhbD1BS0lBVkNPRFlMU0E1M1BRSzRaQSUyRjIwMjUwMTIyJTJGdXMtZWFzdC0xJTJGczMlMkZhd3M0X3JlcXVlc3QmWC1BbXotRGF0ZT0yMDI1MDEyMlQxNDIxMTNaJlgtQW16LUV4cGlyZXM9MzAwJlgtQW16LVNpZ25hdHVyZT05YzdhNzc3MmRiY2U3OWRkNGFiMmNlMGUzMmMxOTUwMjVkYTE5NjJlYTNkN2VkYzM1ZmU4MjJhMjFmNWM4ZTI4JlgtQW16LVNpZ25lZEhlYWRlcnM9aG9zdCJ9.KNTP2O5jrDyt33JKyKA5xvi7JeBJMm7G5xxsQsilKNo)

## Installation

```sh
go install github.com/shutils/lazyreview@latest
```

If you want use unstable version, run the following command:

```sh
go install github.com/shutils/lazyreview@dev
```

## Usage

### Commands

```sh
lazyreview --config <config-file>
```

### Configuration

This application is configured with a toml file.

for example:

```toml
type = "azure"
key = "<your-key>"
endpoint = "<your-endpoint>"
version = "<your-version>"
model = "<your-model>"
target = "."
output = "__dev/reviews.json"
ignores = [".git"]
prompt = '''
You are a code reviewer. Please review the user's code based on the following points.

1. Code quality
2. Code readability
3. Code efficiency
4. Code security
5. Code maintainability
6. Code scalability
7. Typos and bugs

Please provide appropriate suggestions in Markdown format when answering.
'''
max_tokens = 4000
glamour = "dark"
# collector = "git diff --name-only"
# previewer = "git diff --unified=20"
```

#### type

The type field specifies the type of endpoint. Currently, `azure` and `openai` are supported.
If you set the type to `azure`, you must provide the `key`, `endpoint`, `version`, and `model` fields.
If you want to use the `openai` endpoint, you must provide the `key` and `model` fields.

#### key

The key field is the API key for the endpoint.

#### endpoint

The endpoint field is the endpoint for the API.

#### version

The version field is the deployments version of the model.

#### model

The model field is the model name.
for example, `gpt-4o-mini`

#### target

The target field is the target directory to review.

#### output

The output field is the output file to save the reviews.

#### ignores

The ignores field is the list of directories to ignore.

#### prompt

The prompt field is the prompt to use for the model.
If you not provide the prompt field, the default prompt will be used.

#### max_tokens

The max_tokens field is the maximum number of tokens to use for the model.

#### glamour

The glamour field is the glamour style to use for the output.
for example, `dark` or `light`

If you not provide the glamour field, the output will be plain text.

#### collector

The collector field is the command to use to collect the items to review.
for example, `git diff --name-only`

If you not provide the collector field, all files under the target will be reviewed.

#### previewer

The previewer field is the previewer to use for the output.
for example, `git diff --unified=20`

If you not provide the previewer field, the output will be plain text.
