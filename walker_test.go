package walker

import (
    "testing"
    "net/url"
)

func TestCreateWalker(t *testing.T) {
    w := New("https://digitalocean.com/")

    if w.RootNode.URL.Host != "digitalocean.com" {
        t.Error(w.RootNode.URL, "was supposed to have a Host of digitalocean.com")
    }
}

func TestURLBuilding(t *testing.T) {
    w := New("https://digitalocean.com/")

    var u *url.URL

    u = w.BuildURL("https://google.ca")

    if u.Scheme != "https" {
        t.Error(u, "should have a scheme of https")
    }

    if u.Host != "google.ca" {
        t.Error(u, "should have a host of google.ca")
    }
}

func TestfixMissingURLData(t *testing.T) {
    // Create a new domain search
    w := New("https://digitalocean.com/")

    u, err := url.Parse("http://google.ca")

    if err != nil {
        t.Fatalf("URL parse yielded an error that is unrecoverable")
    }

    w.fixMissingURLData(u)

    if u.Scheme != "http" {
        t.Error(u, "should still have the same scheme")
    }

    if u.Host != "google.ca" {
        t.Error(u, "should still ahve the same host")
    }

    u, err = url.Parse("/about")

    if err != nil {
        t.Fatalf("URL parse yielded an error that is unrecoverable")
    }

    w.fixMissingURLData(u)

    if u.Scheme != "https" {
        t.Error(u, "should have been corrected to add https")
    }

    if u.Host != "digitalocean.com" {
        t.Error(u, "should have had its Host changed to digitalocean.com")
    }

    u, err = url.Parse("//digitalocean.com/")

    w.fixMissingURLData(u)

    if u.Scheme == "https" {
        t.Error(u, "should not have changed the protocol")
    }

    if u.Host != "digitalocean.com" {
        t.Error(u, "should not have modified the Host")
    }
}

func TestURLMatching(t *testing.T) {
    // Create a new domain search
    w := New("https://digitalocean.com/")

    // Define a list of urls that should pass validation
    passURLs := []string { "https://digitalocean.com/", "https://digitalocean.com", "https://digitalocean.com/about", "about", "/about" }

    var u *url.URL
    var onDomain bool

    for _, rawURL := range passURLs {
        u = w.BuildURL(rawURL)
        onDomain = w.MatchURL(u)

        if !onDomain {
            t.Error(rawURL, "should be on the domain but it isn't?")
            break
        }
    }

    failURLs := []string { "http://digitalocean.com", "http://google.ca", "http://23" }

    for _, rawURL := range failURLs {
        u = w.BuildURL(rawURL)
        onDomain = w.MatchURL(u)

        if onDomain {
            t.Error(u, "should not be on the domain but it is? Host:", u.Host, "Scheme:", u.Scheme)
            break
        }
    }

}

func TestAddURL(t *testing.T) {
    // Create a new domain search
    w := New("https://digitalocean.com/")

    // Define a list of urls that should pass validation
    passURLs := []string { "https://digitalocean.com/", "https://digitalocean.com", "https://digitalocean.com/about", "contact", "/contact/us" }

    var u *url.URL
    var urlAdded bool

    for _, rawURL := range passURLs {
        u = w.BuildURL(rawURL)
        urlAdded = w.AddURL(u, w.RootNode)

        if !urlAdded {
            t.Error(rawURL, "should be added to the walker.")
            break
        }
    }

    for _, rawURL := range passURLs {
        u = w.BuildURL(rawURL)
        urlAdded = w.AddURL(u, w.RootNode)

        if urlAdded {
            t.Error(rawURL, "should not have been added to the walker (duplicate added).")
            break
        }
    }
}

func TestNodeBuilding(t *testing.T) {
    // Create a new domain search
    w := New("https://digitalocean.com/")

    parentURL := w.BuildURL("/company/about/")

    parentNode := w.BuildNode(nil, parentURL)

    if parentNode.ParentNode != parentNode {
        t.Error(parentNode, "was supposed to be parented with itself as it is the first element")
    }

    childURL := w.BuildURL("about/us")
    childNode := w.BuildNode(parentNode, childURL)

    if childNode.ParentNode != parentNode {
        t.Error(childNode, "was supposed to be the child of", parentNode)
    }
}

func TestGetBodyViaNode(t *testing.T) {
    // Create a new domain search
    w := New("https://digitalocean.com/")

    u := w.BuildURL("about")

    node := w.BuildNode(nil, u)

    body := getBodyViaNode(node)

    if len(body) <= 0 {
        t.Error(body, "should have had html in it!")
    }
}

var sampleHtmlBody string = "<div class=\"sidebar clear\"><a href=\"/\"><img src=\"img/wyatt.jpg\" alt=\"Wyatt Johnson\" class=\"sidebar__profile-img\"></a><h1 class=\"sidebar__profile-name\">Wyatt Johnson</h1><div class=\"sidebar__profile-bio clear\"><p>Devops and web for <a href=\"https://bigpixel.ca/\">Big Pixel Creative</a> <br>Astrophysics/CompSci Graduate <br>JavaScript and HTML5 <br>software engineer.</p></div><div class=\"sidebar__profile-social\"><a href=\"//twitter.com/wyattjoh\" class=\"twitter-follow-button\" data-show-count=\"true\">Follow @wyattjoh</a></div><nav class=\"sidebar__profile-nav\"><ul><li><a href=\"/\">Posts</a></li><li><a href=\"/about\">About</a></li><li><a href=\"//wyattjoh.ca/feed.xml\">RSS</a></li><li><a href=\"//github.com/wyattjoh\">GitHub</a></li></ul></nav></div><div class=\"main\" role=\"main\"><div class=\"posts\"><ul class=\"posts__list\"><li><a href=\"/fix-grub-after-windows-installation\"><h2 class=\"posts__title\">Fix GRUB after Windows Installation</h2></a></li></ul></div></div>"

func TestParseBody(t *testing.T) {
    parsedBody := parseBody(sampleHtmlBody)

    if parsedBody == nil {
        t.Error("Should have generated a parsed body response")
    }
}

func TestGetNode (t *testing.T) {
    // Create a new domain search
    w := New("https://digitalocean.com/")

    u := w.BuildURL("/about/company")

    node := w.BuildNode(nil, u)

    doc := w.GetNode(node)

    if doc == nil {
        t.Error("The doc", doc, "should not have been nil")
    }
}

func TestWalkPage(t *testing.T) {
    // Create a new domain search
    w := New("https://digitalocean.com/")

    u := w.BuildURL("/about/company")

    node := w.BuildNode(nil, u)

    doc := w.GetNode(node)

    w.WalkPage(node, doc)

    if node.Children.Len() == 0 {
        t.Error("There should have been at least one link on", u)
    }

    if node.Assets.Len() == 0 {
        t.Error("There should have been at least on asset on", u)
    }
}

func BenchmarkWalkPage(b *testing.B) {
    w := New("http://wyattjoh.ca/")

    u := w.BuildURL("/")

    for i := 0; i < b.N; i++ {
        node := w.BuildNode(nil, u)

        doc := parseBody(sampleHtmlBody)

        w.WalkPage(node, doc)
    }
}
