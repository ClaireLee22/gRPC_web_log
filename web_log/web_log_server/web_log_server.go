package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"grpc_web_log/web_log/web_log_pb"
	"grpc_web_log/weblogger"
	"io"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type configuration struct {
	Port string `json:"port"`
}

type server struct{}

// Article is a single article
/*
struct fields must be exported (begins with a capital letter)
or they won't be encoded
*/
type Article struct {
	ArticleID string `json:"articleID"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}

// Articles is a slice with multiple articles
type Articles []Article

// file path
const (
	accessLogFilePath = "logger/access.log"
	errorLogFilePath  = "logger/error.log"
	savedJSONFile     = "conf/saveArticles.json"
	confFile          = "conf/conf.json"
)

// web loggers
var (
	accessWebLogger weblogger.Weblogger
	errorWebLogger  weblogger.Weblogger
)

// Get environment variables from json file
func (config *configuration) getEnvVariables() error {
	// method 1: Decode json file
	file, err := os.Open(confFile)
	if err != nil {
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.ClientIP = "ClientIP is NOT existed."
		errorWebLogger.FatalPrintln(getCurrentRPCmethod(pc), "Open config file error.", err)
		return err
	}
	decoder := json.NewDecoder(file)
	decoderErr := decoder.Decode(&config)
	if decoderErr != nil {
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.ClientIP = "ClientIP is NOT existed."
		errorWebLogger.FatalPrintln(getCurrentRPCmethod(pc), "Decode config file error.", decoderErr)
		return decoderErr
	}

	// method 2: Unmarshal json file
	// jsonData, err := ioutil.ReadFile(confFile)
	// if err != nil {
	// 	return err
	// }
	// unmarshalErr := json.Unmarshal(jsonData, &config)
	// if unmarshalErr != nil {
	// 	return unmarshalErr
	// }
	return nil
}

// Get gRPC client IP address for weblogger
func getClientIP(ctx context.Context) string {
	/* method 1: using peer package get ip address */
	// get tcpAddrn w/o port
	// tcpAddr, ok := pr.Addr.(*net.TCPAddr)
	// addr = tcpAddr.IP.String()
	var addr string
	if pr, ok := peer.FromContext(ctx); ok {
		addr = pr.Addr.String()
	}
	return addr

	/* method 2: replace [::1] to 127.0.0.1 */
	// if client and server at the same machine, it will return [::1]:port
	// https://razil.cc/post/2018/09/go-grpc-get-requestaddr/
	// pr, ok := peer.FromContext(ctx)
	// if !ok {
	// 	return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	// }
	// if pr.Addr == net.Addr(nil) {
	// 	return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	// }

	// addSlice := strings.Split(pr.Addr.String(), ":")
	// if addSlice[0] == "[" {
	// 	return "127.0.0.1"
	// }
	// return addSlice[0]

	/* method 3: use metadata get IP address : localhost:50051 */
	// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
	// md, _ := metadata.FromIncomingContext(ctx)
	// fmt.Println("md[:authority]", md[":authority"][0])
	// return md[":authority"][0]
}

// Get current calling RPC method for weblogger
func getCurrentRPCmethod(pc uintptr) string {
	fn := runtime.FuncForPC(pc).Name() // output: main.(*server).GetAllArticles
	currentMethod := strings.Split(fn, ".")
	rpcMethod := currentMethod[len(currentMethod)-1]
	return rpcMethod
}

// Generate UUID for each article
// https://yourbasic.org/golang/generate-uuid-guid/
func generateUUID() string {
	// byte generator
	b := make([]byte, 16)
	// reads 16 cryptographically secure pseudorandom numbers from rand.Reader and writes them to a byte slice.
	_, err := rand.Read(b)
	if err != nil {
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.FatalPrintln(getCurrentRPCmethod(pc), "Generate UUID error.", err)
	}
	// The slice should now contain random bytes
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

// Read saved articles from saveArticles.json file which had saved articles from the client and return json-encoded data
func getJSONData() []byte {
	// if the json file is not existed, create a new file
	if _, err := os.Stat(savedJSONFile); os.IsNotExist(err) {
		os.Create(savedJSONFile)
		// it is necessary to have a object in json file or it will raise error
		wtiteBrackets := []byte("[]")
		writeErr := ioutil.WriteFile(savedJSONFile, wtiteBrackets, 0644)
		if writeErr != nil {
			pc, _, _, _ := runtime.Caller(0)
			errorWebLogger.FatalPrintln(getCurrentRPCmethod(pc), "Write to json file error.", writeErr)
		}
	}

	jsonData, err := ioutil.ReadFile(savedJSONFile)
	if err != nil {
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.FatalPrintln(getCurrentRPCmethod(pc), "Write to json file error.", err)
	}
	return jsonData
}

// Decode json data and store in currentArticles struct
func getCurrentArticles(jsonData []byte) Articles {
	var currentArticles Articles
	unmarshalErr := json.Unmarshal(jsonData, &currentArticles)
	if unmarshalErr != nil {
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.FatalPrintln(getCurrentRPCmethod(pc), "Umarshal json file error.", unmarshalErr)
	}
	return currentArticles
}

// Write currentArticles struct to json file
func (currentArticles *Articles) write2jsonFile() {
	jsonFile, _ := json.MarshalIndent(&currentArticles, "", "  ")
	// Permissions: 1 – execute, 2 – write, 4 – read
	writeErr := ioutil.WriteFile(savedJSONFile, jsonFile, 0644)
	if writeErr != nil {
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.FatalPrintln(getCurrentRPCmethod(pc), "Write to json file error.", writeErr)
	}
}

// Check whether if request ID is existed
func isIDExist(currentArticles Articles, reqID string) (bool, string) {
	for idx, article := range currentArticles {
		if article.ArticleID == reqID {
			return true, strconv.Itoa(idx)
		}
	}
	return false, "not exist"
}

// gRPC service for SaveAllArticles
func (*server) SaveAllArticles(stream web_log_pb.WebLogService_SaveAllArticlesServer) error {
	fmt.Println("SaveAllArticles function was invoked with a streaming request")
	accessWebLogger.ClientIP = getClientIP(stream.Context())
	errorWebLogger.ClientIP = getClientIP(stream.Context())
	// return pc, filename, line, ok
	pc, _, _, _ := runtime.Caller(0)

	jsonData := getJSONData()
	// current articles in the json file
	currentArticles := getCurrentArticles(jsonData)

	var readArticles bytes.Buffer // save readed articles as string buffer
	var logBuffer bytes.Buffer
	for {
		req, err := stream.Recv()
		// req == nil and len(req.String()) != 0
		// read empty file will send empty_req but empty_req != nil but len(empty_req.String()) == 0
		if req != nil && len(req.String()) != 0 {
			s := strings.Split(req.GetArticle(), "\n")
			articleID := generateUUID()
			// setup value for each field in Article struct
			inputArticle := Article{ArticleID: articleID, Title: s[0], Content: s[1]}
			readArticles.WriteString("articleID: " + articleID + "\n")
			readArticles.WriteString("title: " + s[0] + "\n\n")
			currentArticles = append(currentArticles, inputArticle)
			logBuffer.WriteString(req.GetArticle())
		}

		var result bytes.Buffer // server response
		if err == io.EOF {
			if len(readArticles.String()) == 0 {
				pc, _, _, _ := runtime.Caller(0)
				errorWebLogger.ErrorPrintln(getCurrentRPCmethod(pc), "It is an empty file.")
				result.WriteString("It is an empty file. NO new article is saved")
			} else {
				result.WriteString("All new articles have been saved")
			}
			accessWebLogger.AccessPrintln(getCurrentRPCmethod(pc), logBuffer.String())

			// Save json file
			currentArticles.write2jsonFile()
			return stream.SendAndClose(
				&web_log_pb.SaveAllArticlesResponse{
					Result: result.String(),
				})
		}
		if err != nil {
			pc, _, _, _ := runtime.Caller(0)
			errorWebLogger.FatalPrintln(getCurrentRPCmethod(pc), "Error while reading client stream.", err)
		}
	}
}

// gRPC service for GetAllArticles
func (*server) GetAllArticles(ctx context.Context, req *web_log_pb.GetAllArticlesRequest) (*web_log_pb.GetAllArticlesResponse, error) {
	fmt.Println("GetArticles function was invoked with a streaming request")
	accessWebLogger.ClientIP = getClientIP(ctx)
	errorWebLogger.ClientIP = getClientIP(ctx)
	pc, _, _, _ := runtime.Caller(0)

	jsonData := getJSONData()
	// current articles in the json file
	currentArticles := getCurrentArticles(jsonData)

	var result bytes.Buffer // server response (using string buffer to concate strings)
	if len(currentArticles) == 0 {
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.ErrorPrintln(getCurrentRPCmethod(pc), "No article is available now.")
		result.WriteString("No article is available now.")
	} else {
		for _, article := range currentArticles {
			result.WriteString("\narticleID: " + article.ArticleID + "\ntitle: " + article.Title + "\n")
		}
	}

	res := &web_log_pb.GetAllArticlesResponse{
		Result: result.String(),
	}
	accessWebLogger.AccessPrintln(getCurrentRPCmethod(pc), "")
	return res, nil
}

// gRPC service for GetSpecifiedArticle
func (*server) GetSpecifiedArticle(ctx context.Context, req *web_log_pb.GetSpecifiedArticleRequest) (*web_log_pb.GetSpecifiedArticleResponse, error) {
	fmt.Printf("GetSpecifiedArticle function was invoked with %v\n", req)
	accessWebLogger.ClientIP = getClientIP(ctx)
	errorWebLogger.ClientIP = getClientIP(ctx)
	pc, _, _, _ := runtime.Caller(0)

	jsonData := getJSONData()

	// current articles in the json file
	currentArticles := getCurrentArticles(jsonData)

	title := ""
	content := ""
	isExist, idx := isIDExist(currentArticles, req.ArticleID)
	i, _ := strconv.Atoi(idx)

	if isExist {
		title = currentArticles[i].Title
		content = currentArticles[i].Content
	} else {
		title = "title NOT exist"
		content = "content NOT exist"
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.ErrorPrintln(getCurrentRPCmethod(pc), "articleID is NOT existed.")
	}

	// Create response
	res := &web_log_pb.GetSpecifiedArticleResponse{
		ArticleID: req.ArticleID,
		Title:     title,
		Content:   content,
	}

	var logBuffer bytes.Buffer
	logBuffer.WriteString("articleID=" + req.GetArticleID())
	accessWebLogger.AccessPrintln(getCurrentRPCmethod(pc), logBuffer.String())
	return res, nil
}

// gRPC service for UpdateSpecifiedArticle
func (*server) UpdateSpecifiedArticle(ctx context.Context, req *web_log_pb.UpdateSpecifiedArticleRequest) (*web_log_pb.UpdateSpecifiedArticleResponse, error) {
	fmt.Printf("UpdateSpecifiedArticle function was invoked with %v\n", req)
	accessWebLogger.ClientIP = getClientIP(ctx)
	errorWebLogger.ClientIP = getClientIP(ctx)
	pc, _, _, _ := runtime.Caller(0)

	jsonData := getJSONData()
	// current articles in the json file
	currentArticles := getCurrentArticles(jsonData)

	isExist, idx := isIDExist(currentArticles, req.ArticleID)
	i, _ := strconv.Atoi(idx) // string to int

	var result bytes.Buffer
	if isExist {
		currentArticles[i].Title = req.Title
		currentArticles[i].Content = req.Content
		// Save json file
		currentArticles.write2jsonFile()
		result.WriteString("The article with aricleID " + req.ArticleID + " has been updated")
	} else {
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.ErrorPrintln(getCurrentRPCmethod(pc), "articleID is NOT existed.")
		result.WriteString("The article with aricleID " + req.ArticleID + " is NOT existed")
	}

	// Create response
	res := &web_log_pb.UpdateSpecifiedArticleResponse{
		Result: result.String(),
	}

	var logBuffer bytes.Buffer
	logBuffer.WriteString("articleID=" + req.GetArticleID())
	accessWebLogger.AccessPrintln(getCurrentRPCmethod(pc), logBuffer.String())
	return res, nil
}

// gRPC service for RemoveSpecifiedArticle
func (*server) RemoveSpecifiedArticle(ctx context.Context, req *web_log_pb.RemoveSpecifiedArticleRequest) (*web_log_pb.RemoveSpecifiedArticleResponse, error) {
	fmt.Printf("RemoveSpecifiedArticle function was invoked with %v\n", req)
	accessWebLogger.ClientIP = getClientIP(ctx)
	errorWebLogger.ClientIP = getClientIP(ctx)
	pc, _, _, _ := runtime.Caller(0)

	jsonData := getJSONData()
	// current articles in the json file
	currentArticles := getCurrentArticles(jsonData)

	isExist, idx := isIDExist(currentArticles, req.ArticleID)
	i, _ := strconv.Atoi(idx) // string to int

	var result bytes.Buffer
	if isExist {
		// Remove request article
		currentArticles = append(currentArticles[:i], currentArticles[i+1:]...)
		// Save json file
		currentArticles.write2jsonFile()
		result.WriteString("The article with articleID " + req.ArticleID + " has been removed")
	} else {
		result.WriteString("The article with articleID " + req.ArticleID + " is NOT existed")
		pc, _, _, _ := runtime.Caller(0)
		errorWebLogger.ErrorPrintln(getCurrentRPCmethod(pc), "articleID is NOT existed.")
	}

	// Create response
	res := &web_log_pb.RemoveSpecifiedArticleResponse{
		Result: result.String(),
	}

	var logBuffer bytes.Buffer
	logBuffer.WriteString("articleID=" + req.GetArticleID())
	accessWebLogger.AccessPrintln(getCurrentRPCmethod(pc), logBuffer.String())
	return res, nil
}

// main function
func main() {
	fmt.Println("Server(go) is on !")
	// read conf.json file
	config := configuration{}
	readConfigErr := config.getEnvVariables()

	if readConfigErr != nil {
		errorWebLogger.ClientIP = "Client IP is NOT existed."
		errorWebLogger.ServerFatalPrintln("Failed to read config file.", readConfigErr)
	}
	// init weblogger
	accessWebLogger.InitWebLogger(accessLogFilePath)
	errorWebLogger.InitWebLogger(errorLogFilePath)

	lis, err := net.Listen("tcp4", "0.0.0.0:"+config.Port)

	// another way to get port
	// port := "0.0.0.0:" + os.Getenv("port")
	// lis, err := net.Listen("tcp4", port)

	if err != nil {
		errorWebLogger.ServerFatalPrintln("Failed to listen.", err)
	}

	s := grpc.NewServer()
	web_log_pb.RegisterWebLogServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		errorWebLogger.ServerFatalPrintln("Failed to serve.", err)
	}

	// Disable logger
	// ioutil.Discard: a writer on which all calls succeed without doing anything
	accessWebLogger.Logger.SetOutput(ioutil.Discard)
	errorWebLogger.Logger.SetOutput(ioutil.Discard)
}
