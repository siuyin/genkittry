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
	"github.com/siuyin/dflt"
)

type MenuItem struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Calories    int      `json:"calories"`
	Allergens   []string `json:"allergens"`
}

func main() {
	g, model := getModel()

	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	resp, err := genkit.Generate(timeout, g, ai.WithModel(model),
		ai.WithPrompt("Invent a menu item for a pirate themed restaurant."),
		ai.WithOutputType(&MenuItem{}),
	)
	if err != nil {
		log.Fatal("generate: ", err)
	}

	var item MenuItem
	if err := resp.Output(&item); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s (%d calories, %d allergens): %s\n",
		item.Name, item.Calories, len(item.Allergens), item.Description)
	fmt.Printf("\n\n%#v\n", item)

	genData(timeout, g, model)
}

func getModel() (*genkit.Genkit, ai.Model) {
	baseURL := dflt.EnvString("BASE_URL", "http://localhost:11434/v1")
	modelName := dflt.EnvString("MODEL", "mistral") // mistral:7b, gemma3:4b also works well
	log.Printf("BASE_URL=%s MODEL=%s", baseURL, modelName)

	ctx := context.Background()
	var err error

	mySvr := &openai.OpenAI{APIKey: "Ollama", Opts: []option.RequestOption{option.WithBaseURL(baseURL)}}

	g, err := genkit.Init(ctx, genkit.WithPlugins(mySvr))
	if err != nil {
		log.Fatalf("could not initialize Genkit: %v", err)
	}
	model, err := mySvr.DefineModel(g, modelName,
		ai.ModelInfo{Supports: &ai.ModelSupports{Multiturn: true, SystemRole: true, Tools: false, Media: false}},
	)
	if err != nil {
		log.Fatal(err)
	}
	return g, model
}

func genData(ctx context.Context, g *genkit.Genkit, model ai.Model) {
	item, _, err := genkit.GenerateData[MenuItem](ctx, g, ai.WithModel(model),
		ai.WithPrompt("Invent a menu item for a pirate themed restaurant."),
	)
	if err != nil {
		log.Fatal("generate: ", err)
	}
	fmt.Printf("\nNew Item:\n%#v\n", item)
}
