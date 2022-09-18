# aai - ask ai cli

A command line tool that helps you find a command you need.
It uses AI to suggest a command based on your query.

It is advised **not to use the suggestions blindly**,
but rather to read the documentation and understand
what the command does before running it.

OpenAI is the only provider currently supported.

### Usage examples

```bash
$ aai "kill process at port 8080"
kill $(lsof -t -i:8080)
```

```bash
$ aai "find all yaml files in subdirs"
find . -name "*.yaml"
```

## Getting started
### Install:

#### Brew tap
aai can be installed with homebrew using
```bash
brew tap TomaszDomagala/ask-ai-cli
brew install aai
```
#### Release
Download the latest release [here](https://github.com/TomaszDomagala/ask-ai-cli/releases/latest).

#### From source

```bash
git clone git@github.com:TomaszDomagala/ask-ai-cli.git
cd ask-ai-cli
make build
```

### Setup
aai will search for config.yaml file in the following locations: `/etc/aai/`, `$HOME/.aai/`, `.`. Before using it, we need to create the config file.

Set the OpenAI API key and create the file in `/etc/aai/`, you can create one [here](https://beta.openai.com/account/api-keys).
```bash
aai config set --create --openai-apikey sk-XXX
```

You can permanently set OpenAI request options with flags, for example, you can change the model that is used to generate suggestions.
```bash
aai config set --openai-model code-davinci-002
