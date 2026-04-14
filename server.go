package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func handleSearch(w http.ResponseWriter, r *http.Request) {
	// Context Declaration
	ctx := r.Context()

	// Loading in the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning : No .env, by extension no valid API keys found to initialise tool !")
		return
	}

	// Reading in the data string
	rawBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Body.Close()

	// Client cretion
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		http.Error(w, "empty API_KEY", http.StatusBadRequest)
		return
	}
	authSetting := option.WithAPIKey(key)
	client, err := genai.NewClient(ctx, authSetting)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	// Prompt Creation
	collectedText := string(rawBytes)
	prompt := "You are an the summarizer medium for an AI summarizer tool that is termianl based. Summarize the following text, short enough so that it can be displayed effectievely in the terminal. Return nothing except the summarized info. If the result is a code, return only the code. Do not return images of any kind, including markdown images. The content to be summarized is : " + collectedText

	// Calling the AI
	model := client.GenerativeModel("gemini-1.5-flash")
	aiResponse, err := model.GenerateContent(ctx, genai.Text(prompt))

	// Unpack content,
	// transmit info back
	for _, cand := range aiResponse.Candidates {
		if cand.Content == nil {
			return
		}

		for _, part := range cand.Content.Parts {
			fmt.Fprint(w, part)
		}
	}

}

func main() {
	http.HandleFunc("/search", handleSearch)
	fmt.Println("Proxy server : Request acknowledged")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
