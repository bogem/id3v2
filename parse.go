// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/bogem/id3v2/bbpool"
	"github.com/bogem/id3v2/lrpool"
	"github.com/bogem/id3v2/util"
)

const frameHeaderSize = 10

var errBlankFrame = errors.New("id or size of frame are blank")

type frameHeader struct {
	ID       string
	BodySize int64
}

func parseTag(reader io.Reader, opts Options) (*Tag, error) {
	if reader == nil {
		return nil, errors.New("reader is nil")
	}

	header, err := parseHeader(reader)
	if err == errNoTag {
		return newTag(reader, 0, 4), nil
	}
	if err != nil {
		return nil, errors.New("error by parsing tag header: " + err.Error())
	}
	if header.Version < 3 {
		err = errors.New(fmt.Sprint("unsupported version of ID3 tag: ", header.Version))
		return nil, err
	}

	t := newTag(reader, tagHeaderSize+header.FramesSize, header.Version)
	if opts.Parse {
		err = t.parseAllFrames(opts)
	}

	return t, err
}

func newTag(reader io.Reader, originalSize int64, version byte) *Tag {
	tag := &Tag{
		frames:    make(map[string]Framer),
		sequences: make(map[string]*sequence),

		reader:       reader,
		originalSize: originalSize,
		version:      version,
	}

	if file, ok := reader.(*os.File); ok {
		tag.file = file
	}

	return tag
}

func (t *Tag) parseAllFrames(opts Options) error {
	// Size of frames in tag = size of whole tag - size of tag header.
	framesSize := t.originalSize - tagHeaderSize

	// Convert descriptions, specified by user in opts.ParseFrames, to IDs.
	// Use map for speed.
	parseIDs := make(map[string]bool, len(opts.ParseFrames))
	for _, description := range opts.ParseFrames {
		parseIDs[t.CommonID(description)] = true
	}

	for framesSize > 0 {
		// Parse frame header.
		header, err := parseFrameHeader(t.reader)
		if err == io.EOF || err == errBlankFrame || err == util.ErrInvalidSizeFormat {
			break
		}
		if err != nil {
			return err
		}
		id := header.ID
		bodySize := header.BodySize

		// Substitute the size of the whole frame from framesSize.
		framesSize -= frameHeaderSize + bodySize

		// If user set opts.ParseFrames, take it into consideration.
		if len(parseIDs) > 0 {
			if !parseIDs[id] {
				_, err = io.CopyN(ioutil.Discard, t.reader, bodySize)
				continue
			}
		}

		// Limit t.file by header.BodySize.
		bodyRd := lrpool.Get(t.reader, bodySize)
		defer lrpool.Put(bodyRd)

		// Parse frame body.
		frame, err := parseFrameBody(id, bodyRd)
		if err != nil && err != io.EOF {
			return err
		}

		// Add frame to tag.
		t.AddFrame(id, frame)

		if err == io.EOF {
			break
		}
	}

	return nil
}

func parseFrameHeader(rd io.Reader) (frameHeader, error) {
	var header frameHeader

	// Limit rd by frameHeaderSize.
	bodyRd := lrpool.Get(rd, frameHeaderSize)
	defer lrpool.Put(bodyRd)

	fhBuf := bbpool.Get()
	defer bbpool.Put(fhBuf)

	_, err := fhBuf.ReadFrom(bodyRd)
	if err != nil {
		return header, err
	}
	data := fhBuf.Bytes()

	id := string(data[:4])
	bodySize, err := util.ParseSize(data[4:8])
	if err != nil {
		return header, err
	}

	if id == "" || bodySize == 0 {
		return header, errBlankFrame
	}

	header.ID = id
	header.BodySize = bodySize
	return header, nil

}

func parseFrameBody(id string, rd io.Reader) (Framer, error) {
	if id[0] == 'T' {
		return parseTextFrame(rd)
	}

	if parseFunc, exists := parsers[id]; exists {
		return parseFunc(rd)
	}

	return parseUnknownFrame(rd)
}
