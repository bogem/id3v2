package id3v2

import (
	"bytes"
	"testing"
)

func TestDecodeText(t *testing.T) {
	testCases := []struct {
		src  []byte
		from Encoding
		utf8 string
	}{
		{[]byte{0x48, 0xE9, 0x6C, 0x6C, 0xF6}, EncodingISO, "Héllö"},
		{[]byte{0xFF, 0xFE, 0x48, 0x00, 0xE9, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0xF6, 0x00}, EncodingUTF16, "Héllö"},
		{[]byte{0xFF, 0xFE}, EncodingUTF16, ""},
		{[]byte{0x00, 0x48, 0x00, 0xE9, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0xF6}, EncodingUTF16BE, "Héllö"},
	}

	for _, tc := range testCases {
		got := decodeText(tc.src, tc.from)
		if got != tc.utf8 {
			t.Errorf("Expected %q from %v encoding, got %q", tc.utf8, tc.from, got)
		}
	}
}

func TestEncodeWriteText(t *testing.T) {
	testCases := []struct {
		src      string
		to       Encoding
		expected []byte
	}{
		{"Héllö", EncodingISO, []byte{0x48, 0xE9, 0x6C, 0x6C, 0xF6}},
		{"Héllö", EncodingUTF16, []byte{0xFE, 0xFF, 0x00, 0x48, 0x00, 0xE9, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0xF6, 0x00}},
		{"Héllö", EncodingUTF16BE, []byte{0x00, 0x48, 0x00, 0xE9, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0xF6}},
	}

	buf := new(bytes.Buffer)
	bw := newBufWriter(buf)

	for _, tc := range testCases {
		buf.Reset()
		bw.Reset(buf)

		bw.EncodeAndWriteText(tc.src, tc.to)
		if err := bw.Flush(); err != nil {
			t.Fatal(err)
		}
		got := buf.Bytes()
		if !bytes.Equal(got, tc.expected) {
			t.Errorf("Expected %q to %q encoding, got %q", tc.expected, tc.to, got)
		}
		if bw.Written() != len(tc.expected) {
			t.Errorf("Expected %v size, got %v", len(tc.expected), bw.Written())
		}
	}
}

// See https://github.com/bogem/id3v2/issues/51.
func TestUnsynchronisedLyricsFrameWithUTF16(t *testing.T) {
	contentDescriptor := "Content descriptor"
	lyrics := "Lyrics"

	frame := UnsynchronisedLyricsFrame{
		Encoding:          EncodingUTF16,
		Language:          "eng",
		ContentDescriptor: contentDescriptor,
		Lyrics:            lyrics,
	}

	buf := new(bytes.Buffer)

	if _, err := frame.WriteTo(buf); err != nil {
		t.Fatal(err)
	}

	parsed, err := parseUnsynchronisedLyricsFrame(newBufReader(buf), 4)
	if err != nil {
		t.Fatal(err)
	}

	uslf := parsed.(UnsynchronisedLyricsFrame)

	if uslf.ContentDescriptor != contentDescriptor {
		t.Errorf("Expected content descriptor: %q, got: %q", contentDescriptor, uslf.ContentDescriptor)
	}

	if uslf.Lyrics != lyrics {
		t.Errorf("Expected lyrics: %q, got: %q", lyrics, uslf.Lyrics)
	}

}
