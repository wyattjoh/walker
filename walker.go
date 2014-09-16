package walker

import (
    "net/url"
    "strings"
    "github.com/willf/bloom"
    "code.google.com/p/go.net/html"
    "log"
)

type Walker struct {
    RootNode *Node
    Filter *bloom.BloomFilter
}

// Creates a new walker object
func New(rawBaseUrl string) (w *Walker) {
    // Parse the base url
    baseUrl, err := url.Parse(rawBaseUrl)

    // Fail on error
    if err != nil {
        log.Fatal("baseURL (", baseUrl, ") failed to parse: ", err)
    }

    rootNode := buildNode(nil, baseUrl)

    // Create the bloom filter
    filter := bloom.New(20000, 5)

    // Create the object
    w = &Walker{rootNode, filter}

    // Return
    return
}

// Fixes the missing data from the url and substitutes the data from the
// BaseUrl
func (w *Walker) fixMissingURLData(u *url.URL) {
    // Check if host is defined
    if u.Host == "" {
        // If there isn't a host, then it's a relative link to the base domain
        u.Host = w.RootNode.URL.Host
        u.Scheme = w.RootNode.URL.Scheme
    }

    if u.Scheme == "" {
        u.Scheme = w.RootNode.URL.Scheme
    }
}

// Builds a URL from a rawURLString and adds missing data automatically
func (w *Walker) BuildURL(rawURLString string) (*url.URL) {
    // Parse out the raw url
    u, err := url.Parse(rawURLString)

    // There was an error parsing the URL
    if err != nil {
        log.Fatal(err)
    }

    // Fix the url data
    w.fixMissingURLData(u)

    return u
}

// Matches
func (w *Walker) MatchURL(u *url.URL) (bool) {
    // Check if the scheme or host is mismatching
    if u.Scheme == w.RootNode.URL.Scheme && u.Host == w.RootNode.URL.Host   {
        return true
    } else {
        return false
    }
}

func (w *Walker) AddURL(parsedURL *url.URL, node *Node) bool {
    urlAsBytes := []byte(parsedURL.String())

    if w.MatchURL(parsedURL) {
        // Test if the url has been crawled
        if w.Filter.Test(urlAsBytes) {
            // This has been crawled
            return false
        } else {
            // This has not been crawled! Add it to the list
            w.Filter.Add(urlAsBytes)
            node.AddURL(parsedURL)

            return true
        }
    } else {
        return false
    }
}

func parseBody(s string) (doc *html.Node) {
    // Parse it with the parser
    doc, err := html.Parse(strings.NewReader(s))

    // Handle error
    if err != nil {
        log.Fatal(err)
    }

    return doc
}

func (w *Walker) GetNode(n *Node) *html.Node {
    // Get the body text
    body := getBodyViaNode(n)

    // Parse it
    doc := parseBody(body)

    // Send it back!
    return doc
}

func extractAttributes(attributes []html.Attribute, attributeNames ...string) (found []bool, attributeValues []string) {
    attributeValues = make([]string, len(attributeNames))
    found = make([]bool, len(attributeNames))

    for i, attributeName := range attributeNames {
        for _, attr := range attributes {
            if attr.Key == attributeName {
                found[i] = true
                attributeValues[i] = attr.Val
                break
            }
        }
    }

    return
}

// Walks the document for a given node and updates internal structures
func (w *Walker) WalkPage(n *Node, doc *html.Node) {
    // If this is an html element
    if doc.Type == html.ElementNode {

        // If this is an anchor link
        if doc.Data == "a" {

            found, data := extractAttributes(doc.Attr, "href")

            if found[0] {
                // Build the URL
                u := w.BuildURL(data[0])

                // Add it to the node if we haven't visited it yet
                w.AddURL(u, n)
            }

        } else if doc.Data == "img" {

            found, data := extractAttributes(doc.Attr, "src")

            if found[0] {
                // Build Asset Object
                asset := &Asset{"img", w.BuildURL(data[0])}

                // Add Asset to Node
                n.AddAsset(asset)
            }

        } else if doc.Data == "link" {

            found, data := extractAttributes(doc.Attr, "href", "rel")

            if found[0] && found[1] {
                // Build Asset Object
                asset := &Asset{data[1], w.BuildURL(data[0])}

                // Add Asset to Node
                n.AddAsset(asset)
            }

        } else if doc.Data == "script" {
            found, data := extractAttributes(doc.Attr, "src")

            if found[0] {
                // Build Asset Object
                asset := &Asset{"js", w.BuildURL(data[0])}

                // Add Asset to Node
                n.AddAsset(asset)
            }
        }
    }

    // Recursively walk the DOM tree
    for c := doc.FirstChild; c != nil; c = c.NextSibling {
        w.WalkPage(n, c)
    }
}

func (w *Walker) Walk() {
    // Step 1. Walk the root
}
