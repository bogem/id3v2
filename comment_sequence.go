// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

// commentSequence stores several unique comment frames.
// Key for commentSequence is language and description,
// so there is only one comment frame with the same language and
// description.
//
// ID3v2 Documentation: "There may be more than one comment frame in each tag,
// but only one with the same language and content descriptor."
type commentSequence struct {
	sequence    map[string]CommentFrame
	framesCache []Framer
}

func newCommentSequence() sequencer {
	return &commentSequence{
		sequence: make(map[string]CommentFrame),
	}
}

func (cs *commentSequence) AddFrame(f Framer) {
	cs.framesCache = nil

	cf := f.(CommentFrame)
	id := cf.Language + cf.Description
	cs.sequence[id] = cf
}

func (cs commentSequence) Count() int {
	return len(cs.sequence)
}

func (cs *commentSequence) Frames() []Framer {
	cache := cs.framesCache
	if len(cache) == 0 {
		cache = make([]Framer, 0, len(cs.sequence))
		for _, f := range cs.sequence {
			cache = append(cache, f)
		}
		cs.framesCache = cache
	}
	return cache
}
