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

// Define the input structure for the tool
type WeatherInput struct {
	Location string `json:"location" jsonschema_description:"Location to get weather for"`
}

func main() {
	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	g, model := getModel()
	getWeatherTool := genkit.DefineTool(
		g, "getWeather", "Gets the current weather in a given location",
		func(ctx *ai.ToolContext, input WeatherInput) (string, error) {
			// Here, we would typically make an API call or database query. For this
			// example, we just return a fixed value.
			log.Printf("Tool 'getWeather' called for location: %s", input.Location)
			return fmt.Sprintf("The current weather in %s is 30°C and sunny.", input.Location), nil
		})

	resp, err := genkit.Generate(timeout, g, ai.WithPrompt("What is the weather in Singapore?"),
		ai.WithModel(model),
		ai.WithTools(getWeatherTool),
	)
	if err != nil {
		log.Fatal("generate: ", err)
	}
	fmt.Println(resp.Text())
}

func getModel() (*genkit.Genkit, ai.Model) {
	baseURL := dflt.EnvString("BASE_URL", "http://localhost:11434/v1")
	modelName := dflt.EnvString("MODEL", "mistral") // mistral:7b.  gemma3:4b does not have tools
	log.Printf("BASE_URL=%s MODEL=%s", baseURL, modelName)

	ctx := context.Background()
	var err error

	mySvr := &openai.OpenAI{APIKey: "Ollama", Opts: []option.RequestOption{option.WithBaseURL(baseURL)}}

	g, err := genkit.Init(ctx, genkit.WithPlugins(mySvr))
	if err != nil {
		log.Fatalf("could not initialize Genkit: %v", err)
	}
	model, err := mySvr.DefineModel(g, modelName,
		ai.ModelInfo{Supports: &ai.ModelSupports{Multiturn: true, SystemRole: true, Tools: true, Media: false}},
	)
	if err != nil {
		log.Fatal(err)
	}
	return g, model
}
