package walker

import (
    "net/url"
    "fmt"
)

type Asset struct {
    Type string
    URL *url.URL
}

func (a *Asset) String() string {
    return fmt.Sprintf("Asset[%s] =\t%s", a.Type, a.URL)
}
