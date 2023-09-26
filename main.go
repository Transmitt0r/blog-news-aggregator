package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/transmitt0r/blogger/newsapi"
	"github.com/transmitt0r/blogger/scraper"
	"golang.org/x/sync/errgroup"
)

// main is the entry point of the blogger application.
// It retrieves top headlines from NewsAPI.org, scrapes the articles, and summarizes them using OpenAI's GPT-3 model.
// The summarized text is then used to create a blog post.
func main() {
	// Parse command line flags
	var (
		newsapiKey = flag.String("newsapi-key", "", "NewsAPI.org API key")
		openaiKey  = flag.String("openai-key", "", "OpenAI API key")
	)
	flag.Parse()

	// Check if required flags are provided
	if *newsapiKey == "" {
		panic("newsapi-key is required")
	}
	if *openaiKey == "" {
		panic("openai-key is required")
	}

	// Create NewsAPI client
	newsapi := newsapi.NewClient(*newsapiKey)

	// Create OpenAI LLM chat client
	llm, err := openai.NewChat(openai.WithToken(*openaiKey), openai.WithModel("gpt-3.5-turbo-16k"))
	if err != nil {
		panic(err)
	}

	// Create MapReduce chains for summarization and blog post creation
	mapChain := chains.NewLLMChain(llm, prompts.NewPromptTemplate(
		"Summarize: {{.article}}?",
		[]string{"article"},
	))
	reduceChain := chains.NewLLMChain(llm, prompts.NewPromptTemplate(
		"You are writing a post to summarize recent events and trends in the tech industry for linkedin. Here are the articles: {{.articles}}. This weeks post is:",
		[]string{"articles"},
	))
	blogChain := chains.NewMapReduceDocuments(mapChain, reduceChain)

	// Get top headlines from NewsAPI
	articles, err := newsapi.GetTopHeadlines("chatgpt")
	if err != nil {
		panic(err)
	}
	if len(articles) == 0 {
		panic("no articles found")
	}

	// Scrape articles and create Document objects
	results := make(chan schema.Document, len(articles))
	group, _ := errgroup.WithContext(context.Background())
	for _, a := range articles {
		a := a
		group.Go(func() error {
			fmt.Printf("Scraping %s\n", a.URL)
			article, err := scraper.Scrape(a.URL)
			if err != nil {
				return err
			}
			result := schema.Document{
				PageContent: article[:1000],
				Metadata: map[string]interface{}{
					"article": a.Title,
					"url":     a.URL,
				},
			}
			results <- result
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
	close(results)

	// Collect Document objects into a slice
	bigResult := []schema.Document{}
	for r := range results {
		bigResult = append(bigResult, r)
	}
	fmt.Printf("Got %d results\n", len(bigResult))

	// Call MapReduce chains to create blog post
	blogChain.MaxNumberOfConcurrent = 5
	r, err := blogChain.Call(context.Background(), map[string]any{"input_documents": bigResult}, chains.WithMaxLength(2000))
	if err != nil {
		panic(err)
	}
	fmt.Println(r["text"])
}
