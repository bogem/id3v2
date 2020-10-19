// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func TestPictureFrame(t *testing.T) {
	t.Parallel()

	// create ID3 "manually"
	buffer := make([]byte, 10 + 10 + 5 + 10 + 20 + 10 + 25)
	pos := 0

	// ID3 header
	buffer[pos] = 'I'	// Magic
	pos++
	buffer[pos] = 'D'
	pos++
	buffer[pos] = '3'
	pos++
	buffer[pos] = 3		// Version (3.0)
	pos++
	buffer[pos] = 0
	pos++
	buffer[pos] = 0		// Flags
	pos++

	buffer[pos] = byte((len(buffer) >> 21) & 0x7F)		// Size
	pos++
	buffer[pos] = byte((len(buffer) >> 14) & 0x7F)
	pos++
	buffer[pos] = byte((len(buffer) >>  7) & 0x7F)
	pos++
	buffer[pos] = byte((len(buffer) >>  0) & 0x7F)
	pos++


	// some tag before
	buffer[pos] = 'P'
	pos++
	buffer[pos] = 'R'
	pos++
	buffer[pos] = 'I'
	pos++
	buffer[pos] = 'V'
	pos++

	// PRIV size
	buffer[pos] = 0
	pos++
	buffer[pos] = 0
	pos++
	buffer[pos] = 0
	pos++
	buffer[pos] = 5
	pos++

	// PRIV flags
	buffer[pos] = 0
	pos++
	buffer[pos] = 0
	pos++

	// PRIV data
	for i := 0; i < 5; i++ {
		buffer[pos] = byte(rand.Uint32())
		pos++
	}




	// APIC tag
	buffer[pos] = 'A'
	pos++
	buffer[pos] = 'P'
	pos++
	buffer[pos] = 'I'
	pos++
	buffer[pos] = 'C'
	pos++

	// APIC size
	buffer[pos] = 0
	pos++
	buffer[pos] = 0
	pos++
	buffer[pos] = 0
	pos++
	buffer[pos] = 5
	pos++

	// APIC flags
	buffer[pos] = 0
	pos++
	buffer[pos] = 0
	pos++

	// APIC data (INVALID)
	for i := 0; i < 5; i++ {
		buffer[pos] = byte(rand.Uint32())	// data
		pos++
	}




	// some tag after
	buffer[pos] = 'G'
	pos++
	buffer[pos] = 'R'
	pos++
	buffer[pos] = 'I'
	pos++
	buffer[pos] = 'D'
	pos++

	// GRID size
	buffer[pos] = 0
	pos++
	buffer[pos] = 0
	pos++
	buffer[pos] = 0
	pos++
	buffer[pos] = 25
	pos++

	// GRID flags
	buffer[pos] = 0
	pos++
	buffer[pos] = 0
	pos++

	// GRID data
	buffer[pos] = '8'				// owner ID
	pos++
	buffer[pos] = 0
	pos++
	buffer[pos] = 0xF3				// symbol
	pos++
	for i := 0; i < 25 - 3; i++ {
		buffer[pos] = byte(rand.Uint32())	// data
		pos++
	}



	// parse these 3 frames
	tag, err := ParseReader(bytes.NewReader(buffer), Options{Parse: true})
fmt.Println("Got tag:", tag)

	if err != nil {
		t.Fatalf("Received error parsing frames: %v", err)
	}

	allFrames := tag.AllFrames()
	if len(allFrames) != 2 {
		t.Fatalf("Expected exactly 2 tags, received %d instead", len(allFrames))
	}

	{
		priv := tag.GetFrames("PRIV")
		if priv == nil {
			t.Fatalf("Expected a PRIV frame")
		}
		if len(priv) != 1 {
			t.Fatalf("Expected exactly one PRIV frame got %d instead", len(priv))
		}
		if priv[0].Size() != 5 {
			t.Fatalf("Expected a PRIV frame to be 5 bytes, it is %d instead", priv[0].Size())
		}
	}

	{
		grid := tag.GetFrames("GRID")
		if grid == nil {
			t.Fatalf("Expected a GRID frame")
		}
		if len(grid) != 1 {
			t.Fatalf("Expected exactly one GRID frame got %d instead", len(grid))
		}
		if grid[0].Size() != 25 {
			t.Fatalf("Expected a GRID frame to be 25 bytes, it is %d instead", grid[0].Size())
		}
	}
}
