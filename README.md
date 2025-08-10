# Experiments with firebase's genkit

## Ollama
There is a problem with the ollama plugin.
The bug was fixed but there was a reversion.

See: https://github.com/firebase/genkit/issues/719

In v0.6.2 the Timeout parameter is not present.

### Ollama in OpenAI compatibility mode

See: helloOllamOpenAICompat. This works well with a timeout context.

## Structured Output
Gemma3 (4b) and Mistral (7b) work well to generate structured output.

Gemma3 (1b) is intermitted. 4b is sometime intermittent as well. Retry in event of error.

## Flows
See flows/main.go .

This is an example of parameterising the theme of the menu item to be generated.
