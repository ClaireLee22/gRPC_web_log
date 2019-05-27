# gRPC_web_log

## Project Overview
### Project Description

Use gRPC to build a client/server system in Go.

  | RPC method  | Client request |
  | :---  | :---  |
  | SaveAllArticles  | doArticleStreaming  |
  | GetAllArticles | doAllArticles  |
  | GetSpecifiedArticle | doSpecifiedArticle |
  | UpdateSpecifiedArticle| doUpdateSpecifiedArticle  |
  | RemoveSpecifiedArticle | doRemoveSpecifiedArticle |
  
  + __Service1__: SaveAllArticles | doArticleStreaming 
  
    - Client: provide a text file with multiple articles to Server
    - Server: receive articles which are streaming from client and save those received articles into a json file for future reference.    
   
  
  + __Service2__:  GetAllArticles | doAllArticles
  
    - Client: request to show all current articles with their articleID and title
    - Server: provide a list of articleIDs and titles which are saved in the json file to Client
    
  + __Service3__: GetSpecifiedArticle | doSpecifiedArticle
  
    - Client: request to show the article's title and content by given a articleID
    - Server: provide a article's title and content with the specified articleID
    
  + __Service4__: UpdateSpecifiedArticle | doUpdateSpecifiedArticle
  
    - Client: request to update the article's title and content by given a articleID
    - Server: update the article's title and content with the specified articleID and sent a response to confirm the update 
  
  + __Service5__: RemoveSpecifiedArticle | doRemoveSpecifiedArticle
  
    - Client: request to remove a article by given a articleID
    - Server: remove the article with the specified articleID and sent a response to confirm that the article has been removed 
    
    

## Getting Started
### Prerequisites

This project requires **Python 3.6** and the following Python libraries installed:

- [NumPy](http://www.numpy.org/)
- [matplotlib](http://matplotlib.org/)
- [scikit-learn](http://scikit-learn.org/stable/)
- [tensorflow](https://www.tensorflow.org/install/pip)


### Run

In a terminal or command window, run one of the following commands:

```bash
python Doc2Doc_comparison.py
```  

