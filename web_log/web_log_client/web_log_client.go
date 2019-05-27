package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"grpc_web_log/web_log/web_log_pb"
	"log"
	"os"

	"google.golang.org/grpc"
)

// text files with articles
const txtFile = "conf/articles.txt"

// variables with articleID
var specifiedArticleID = "dff4bf31-6423-779d-65dc-72c67c96e5ef"
var updateArticleID = "925ee90d-fdf1-9e1a-3146-3fb4331b9023"
var removeArticleID = "7a566f6f-f931-609e-425c-d5b48fa7e4f5"

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := web_log_pb.NewWebLogServiceClient(conn)

	doArticleStreaming(c)

	doAllArticles(c)

	doSpecifiedArticle(c)

	doUpdateSpecifiedArticle(c)

	doRemoveSpecifiedArticle(c)
}

// gRPC client for doArticleStreaming: provide a text file with articles to server
func doArticleStreaming(c web_log_pb.WebLogServiceClient) {
	fmt.Println("\nStarting to do a Article Streaming RPC...")
	stream, err := c.SaveAllArticles(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GetArticles: %v", err)
	}

	fileHandle, err := os.Open(txtFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fileHandle.Close()

	// Create a new Scanner for the file
	fileScanner := bufio.NewScanner(fileHandle)

	var article bytes.Buffer
	request := &web_log_pb.SaveAllArticlesRequest{
		Article: article.String(),
	}
	for fileScanner.Scan() {
		if fileScanner.Text() == "" {
			request = &web_log_pb.SaveAllArticlesRequest{
				Article: article.String(),
			}
			fmt.Printf("Sending req:\n%v\n", request)
			stream.Send(request)
			article.Reset()
		} else {
			if article.String() == "" {
				// separate title and content to different lines
				article.WriteString(fileScanner.Text() + "\n")
			} else {
				article.WriteString(fileScanner.Text())
			}
		}
	}
	request = &web_log_pb.SaveAllArticlesRequest{
		Article: article.String(),
	}
	stream.Send(request)
	fmt.Printf("Sending req:\n%v\n", request)

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reciving response from GetArticles: %v", err)
	}

	fmt.Printf("GetAllArticle response:\n%v\n", res.Result)
}

// gRPC client for doAllArticles: request a list of current articleIDs and titles which are saved in server end
func doAllArticles(c web_log_pb.WebLogServiceClient) {
	fmt.Println("\nStarting to do a All Articles RPC...")

	req := &web_log_pb.GetAllArticlesRequest{}

	res, err := c.GetAllArticles(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling All Articles Rpc: %v", err)
	}
	log.Printf("Response from GetAllArticles: %v\n", res.Result)

}

// gRPC client for doSpecifiedArticle: request to get an article title and content with a specified articleID
func doSpecifiedArticle(c web_log_pb.WebLogServiceClient) {
	fmt.Println("\nStarting to do a Specified Article RPC...")
	req := &web_log_pb.GetSpecifiedArticleRequest{
		ArticleID: specifiedArticleID,
	}
	fmt.Println("req", req)
	res, err := c.GetSpecifiedArticle(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Specified Article Rpc: %v", err)
	}
	log.Printf("Response from GetSpecifiedArticle:\n %v\n %v\n %v\n", res.ArticleID, res.Title, res.Content)
}

// gRPC client for doUpdateSpecifiedArticle: request to update an article title and content with a specified articleID
func doUpdateSpecifiedArticle(c web_log_pb.WebLogServiceClient) {
	fmt.Println("\nStarting to do a Update Specified Article RPC...")
	req := &web_log_pb.UpdateSpecifiedArticleRequest{
		ArticleID: updateArticleID,
		Title:     "update title",
		Content:   "update content",
	}
	fmt.Println("req", req)
	res, err := c.UpdateSpecifiedArticle(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling UpdateSpecified Article Rpc: %v", err)
	}
	log.Printf("Response from UpdateSpecifiedArticle: %v\n", res.Result)
}

// gRPC client for doRemoveSpecifiedArticle: request to remove an article title and content with a specified articleID
func doRemoveSpecifiedArticle(c web_log_pb.WebLogServiceClient) {
	fmt.Println("\nStarting to do a Remove Specified Article RPC...")
	req := &web_log_pb.RemoveSpecifiedArticleRequest{
		ArticleID: removeArticleID,
	}
	fmt.Println("req", req)
	res, err := c.RemoveSpecifiedArticle(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling RemoveSpecified Article Rpc: %v", err)
	}
	log.Printf("Response from RemoveSpecifiedArticle: %v\n", res.Result)
}
