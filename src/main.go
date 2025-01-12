package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/slack-go/slack"
)

// Article represents the structure of a single article.
type Article struct {
	ID          interface{} `json:"id"`
	Title       string      `json:"title"`
	Path        string      `json:"path"`
	PublishedAt string      `json:"published_at"`
}

// Response represents the structure of the API response.
type Response struct {
	Articles []Article `json:"articles"`
}

// fetchParameter fetches a parameter value from AWS SSM Parameter Store.
func fetchParameter(name string, withDecryption bool) (string, error) {
	// セッション作成
	sess, err := session.NewSession(&aws.Config{Region: aws.String("ap-northeast-1")})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := ssm.New(sess)
	res, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(withDecryption),
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch parameter %s: %w", name, err)
	}

	return *res.Parameter.Value, nil
}

// fetchArticles fetches articles from the Zenn API.
func fetchArticles(url string) ([]Article, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return response.Articles, nil
}

// getRandomArticle selects a random article from the list.
func getRandomArticle(articles []Article) Article {
	rand.Seed(time.Now().UnixNano())
	return articles[rand.Intn(len(articles))]
}

// postToSlack posts a message to a Slack channel.
func postToSlack(token, channel, message string) error {
	client := slack.New(token)
	_, _, err := client.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		return fmt.Errorf("error posting to Slack: %w", err)
	}
	return nil
}

// handler is the entry point for the AWS Lambda function.
func handler(ctx context.Context) error {
	// Fetch parameters from AWS SSM Parameter Store
	zennAPIURL, err := fetchParameter("ZENN_API_URL", false)
	if err != nil {
		return fmt.Errorf("failed to fetch ZENN_API_URL: %w", err)
	}

	slackToken, err := fetchParameter("SLACK_TOKEN", true)
	if err != nil {
		return fmt.Errorf("failed to fetch SLACK_TOKEN: %w", err)
	}

	slackChannel, err := fetchParameter("SLACK_CHANNEL", false)
	if err != nil {
		return fmt.Errorf("failed to fetch SLACK_CHANNEL: %w", err)
	}

	// Fetch articles from Zenn API
	articles, err := fetchArticles(zennAPIURL)
	if err != nil {
		return fmt.Errorf("failed to fetch articles: %w", err)
	}

	// Select a random article
	randomArticle := getRandomArticle(articles)

	// Format the Slack message
	message := fmt.Sprintf(
		"*おすすめの記事*\nTitle: %s\nPath: https://zenn.dev%s\nPublishedAt: %s",
		randomArticle.Title, randomArticle.Path, randomArticle.PublishedAt,
	)

	// Post the message to Slack
	if err := postToSlack(slackToken, slackChannel, message); err != nil {
		return fmt.Errorf("failed to post to Slack: %w", err)
	}

	log.Println("Message posted successfully!")
	return nil
}

func main() {
	lambda.Start(handler)
}
