# captionbot
A wrapper for [CaptionBot](https://www.captionbot.ai) written in Go.

Usage:
```Go
package main

import "github.com/dattish/captionbot"

func main() {
    //captions a local file
    captionbot.FileCaption("image.png")
    
    //captions an image from a url
    captionbot.URLCaption("https://example.com/image.png")
    
    //rate the caption
    captionbot.Rate(5)
}
```