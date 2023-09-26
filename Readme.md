# Blogger news aggregator
The Blogger application is a tool designed to automatically generate summarized blog posts from top headlines sourced from NewsAPI.org. It scrapes the articles associated with these headlines and then uses OpenAI's GPT-3 model to produce a concise summary. The summarized text is then used to create a blog post.

## Features:
1. **NewsAPI Integration**: Fetches top headlines from NewsAPI.org.
2. **Article Scraper**: Scrapes the content of articles linked in the headlines.
3. **OpenAI GPT-3 Summarization**: Uses OpenAI's GPT-3 model to summarize the scraped articles.
4. **Blog Post Creation**: Combines the summarized articles to create a comprehensive blog post.

## How to Use:

### 1. Setup:
Ensure you have the required dependencies installed. Navigate to the project directory and run:
```bash
go mod tidy
```

### 2. Command Line Flags:
- **-newsapi-key**: Your NewsAPI.org API key.
- **-openai-key**: Your OpenAI API key.

Example:
```bash
go run main.go -newsapi-key=YOUR_NEWSAPI_KEY -openai-key=YOUR_OPENAI_KEY
```

### 3. Output:
The program will scrape the articles, summarize them, and then print the final blog post to the console.

## Important Notes:
- Ensure you provide both the newsapi-key and openai-key when running the program. They are mandatory for the application to function.
- The scraper is set to scrape the first 1000 characters of each article. Adjust this limit if needed.
- The MapReduce chains for summarization and blog post creation have a concurrency limit set to 5. Adjust this value based on your system's capabilities.

## Troubleshooting:
- If you encounter the error newsapi-key is required or openai-key is required, ensure you have provided the respective keys when running the program.
- If you see the error no articles found, it means that no articles were returned for the specified query from NewsAPI.
- Any other errors are likely related to scraping or the OpenAI API. Check the error message for more details.

## Contributions:
Feel free to contribute to this project by submitting pull requests or opening issues for any bugs or feature requests.