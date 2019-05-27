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
  
  + Service1: SaveAllArticles \ doArticleStreaming 
  
    - Server: receive articles which are streaming from client and save those received articles into a json file for future reference.    
    - Client: provide a text file of articles to server
    
    

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

