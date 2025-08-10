package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

func main() {
	ctx := context.Background()
	var err error

	mySvr := &openai.OpenAI{APIKey: "Ollama", Opts: []option.RequestOption{option.WithBaseURL("http://localhost:11434/v1")}}

	g, err := genkit.Init(ctx, genkit.WithPlugins(mySvr))
	if err != nil {
		log.Fatalf("could not initialize Genkit: %v", err)
	}
	model, err := mySvr.DefineModel(g, "gemma3:1b",
		ai.ModelInfo{Supports: &ai.ModelSupports{Multiturn: true, SystemRole: true, Tools: false, Media: false}},
	)
	if err != nil {
		log.Fatal(err)
	}
	timeout, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	_, err = genkit.Generate(timeout, g, ai.WithModel(model),
		ai.WithPrompt("What is the meaning of life?"),
		ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
			fmt.Print(chunk.Text())
			return nil
		}),
	)
	if err != nil {
		log.Fatal("generate:", err)
	}

	fmt.Println()
}
