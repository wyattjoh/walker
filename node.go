package walker

import (
    "net/url"
    "container/list"
    "log"
    "net/http"
    "io/ioutil"
)

type Node struct {
    ParentNode *Node
    URL *url.URL
    Assets *list.List
    Children *list.List
}

func (n *Node) Init(parentNode *Node, nodeURL *url.URL) *Node {
    // If there is a parent node
    if parentNode != nil {
        // Assign it
        n.ParentNode = parentNode
    } else {
        // Otherwise..

        // Set parent to self
        n.ParentNode = n
    }

    n.URL = nodeURL
    n.Assets = list.New()
    n.Children = list.New()

    return n
}

func buildNode(parentNode *Node, nodeURL *url.URL) *Node {
    return new(Node).Init(parentNode, nodeURL)
}

func (w *Walker) BuildNode(parentNode *Node, nodeURL *url.URL) *Node {
    return buildNode(parentNode, nodeURL)
}

func (n *Node) AddURL(parsedURL *url.URL) {
    n.Children.PushBack(parsedURL)
}

func (n *Node) AddAsset(asset *Asset) {
    n.Assets.PushBack(asset)
}

func getBodyViaNode(n *Node) (s string) {
    // GET the URL
    resp, err := http.Get(n.URL.String())

    // Handle error
    if err != nil {
        log.Fatal(err)
    }

    // Close the reader when we are done...
    defer resp.Body.Close()

    // Read the body
    body, err := ioutil.ReadAll(resp.Body)

    // Handle error
    if err != nil {
        log.Fatal(err)
    }

    // Turn the body into a string
    s = string(body)

    return
}
