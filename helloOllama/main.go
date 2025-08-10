package main

import (
	"context"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
)

func main() {
	ctx := context.Background()
	var err error

	mySvr := &ollama.Ollama{ServerAddress: "http://localhost:11434"}
	g, err := genkit.Init(ctx, genkit.WithPlugins(mySvr))
	if err != nil {
		log.Fatalf("could not initialize Genkit: %v", err)
	}
	model := mySvr.DefineModel(g,
		ollama.ModelDefinition{Name: "gemma3:1b", Type: "chat"},
		&ai.ModelInfo{Supports: &ai.ModelSupports{Multiturn: true, SystemRole: true, Tools: false, Media: false}},
	)

	_, err = genkit.Generate(ctx, g, ai.WithModel(model),
		ai.WithPrompt("What is the meaning of life? Respond in one paragraph."),
		ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
			fmt.Print(chunk.Text())
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}
