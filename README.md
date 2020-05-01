# [id3v2](https://pkg.go.dev/github.com/bogem/id3v2)

**Fast, simple and powerful ID3 decoding and encoding library written in Go.**

id3v2 supports
- versions 2.3 and 2.4
- all available encodings
- all text frames, unsynchronised lyrics/text, comments, attached pictures, UFID and TXXX frames

## Installation

```
go get -u github.com/bogem/id3v2
```

## Usage example

```go
package main

import (
	"fmt"
	"log"

	"github.com/bogem/id3v2"
)

func main() {
	tag, err := id3v2.Open("file.mp3", id3v2.Options{Parse: true})
	if err != nil {
 		log.Fatal("Error while opening mp3 file: ", err)
 	}
	defer tag.Close()

	fmt.Println(tag.Artist())
	fmt.Println(tag.Title())

	tag.SetArtist("Aphex Twin")
	tag.SetTitle("Xtal")

	comment := id3v2.CommentFrame{
		Encoding:    id3v2.EncodingUTF8,
		Language:    "eng",
		Description: "My opinion",
		Text:        "I like this song!",
	}
	tag.AddCommentFrame(comment)

	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}
```

## Read multiple frames

```go
pictures := tag.GetFrames(tag.CommonID("Attached picture"))
for _, f := range pictures {
	pic, ok := f.(id3v2.PictureFrame)
	if !ok {
		log.Fatal("Couldn't assert picture frame")
	}

	// Do something with picture frame.
	// For example, print the description:
	fmt.Println(pic.Description)
}
```

## Work with encodings

For example, if you set comment frame with custom encoding and write it, you may do this:

```go
tag := id3v2.NewEmptyTag()
comment := id3v2.CommentFrame{
	Encoding:    id3v2.EncodingUTF16,
	Language:    "ger",
	Description: "Tier",
	Text:        "Der Löwe",
}
tag.AddCommentFrame(comment)

_, err := tag.WriteTo(w)
if err != nil {
	log.Fatal(err)
}
```

`Text` field will be automatically encoded with UTF-16BE with BOM and written to w.

By default, if version of tag is 4 then UTF-8 is used for methods like
`SetArtist`, `SetTitle`, `SetGenre` and etc, otherwise ISO-8859-1.
