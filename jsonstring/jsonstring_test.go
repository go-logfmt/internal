package jsonstring

import (
	"bytes"
	"testing"
	"unicode"
)

func TestWriteQuoted(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		in  string
		out string
	}{
		{"\x00", `"\u0000"`},
		{"\x01", `"\u0001"`},
		{"\x02", `"\u0002"`},
		{"\x03", `"\u0003"`},
		{"\x04", `"\u0004"`},
		{"\x05", `"\u0005"`},
		{"\x06", `"\u0006"`},
		{"\x07", `"\u0007"`},
		{"\x08", `"\u0008"`},
		{"\x09", `"\t"`},
		{"\x0a", `"\n"`},
		{"\x0b", `"\u000b"`},
		{"\x0c", `"\u000c"`},
		{"\x0d", `"\r"`},
		{"\x0e", `"\u000e"`},
		{"\x0f", `"\u000f"`},
		{"\x10", `"\u0010"`},
		{"\x11", `"\u0011"`},
		{"\x12", `"\u0012"`},
		{"\x13", `"\u0013"`},
		{"\x14", `"\u0014"`},
		{"\x15", `"\u0015"`},
		{"\x16", `"\u0016"`},
		{"\x17", `"\u0017"`},
		{"\x18", `"\u0018"`},
		{"\x19", `"\u0019"`},
		{"\x1a", `"\u001a"`},
		{"\x1b", `"\u001b"`},
		{"\x1c", `"\u001c"`},
		{"\x1d", `"\u001d"`},
		{"\x1e", `"\u001e"`},
		{"\x1f", `"\u001f"`},
	}

	for _, tc := range testCases {
		buf := &bytes.Buffer{}
		_, err := WriteQuotedString(buf, tc.in)
		if err != nil {
			t.Errorf("WriteQuotedString(%q): %v", tc.in, err)
			continue
		}
		if got, want := buf.String(), tc.out; got != want {
			t.Errorf("WriteQuotedString(%q) = %#q, want %#q", tc.in, got, want)
		}

		buf.Reset()
		_, err = WriteQuotedBytes(buf, []byte(tc.in))
		if err != nil {
			t.Errorf("WriteQuotedBytes(%q): %v", tc.in, err)
			continue
		}
		if got, want := buf.Bytes(), []byte(tc.out); !bytes.Equal(got, want) {
			t.Errorf("WriteQuotedBytes(%q) = %#q, want %#q", tc.in, got, want)
		}
	}
}

func TestStringBytes(t *testing.T) {
	t.Parallel()

	// Test that WriteQuotedString and WriteQuotedBytes use the same encoding.
	var r []rune
	for i := '\u0000'; i <= unicode.MaxRune; i++ {
		r = append(r, i)
	}
	s := string(r) + "\xff\xff\xffhello" // some invalid UTF-8 too

	es := &bytes.Buffer{}
	_, err := WriteQuotedString(es, s)
	if err != nil {
		t.Fatalf("WriteQuotedString: %v", err)
	}

	esBytes := &bytes.Buffer{}
	_, err = WriteQuotedBytes(esBytes, []byte(s))
	if err != nil {
		t.Fatalf("WriteQuotedBytes: %v", err)
	}

	enc := es.String()
	encBytes := esBytes.String()
	if enc != encBytes {
		i := 0
		for i < len(enc) && i < len(encBytes) && enc[i] == encBytes[i] {
			i++
		}
		enc = enc[i:]
		encBytes = encBytes[i:]
		i = 0
		for i < len(enc) && i < len(encBytes) && enc[len(enc)-i-1] == encBytes[len(encBytes)-i-1] {
			i++
		}
		enc = enc[:len(enc)-i]
		encBytes = encBytes[:len(encBytes)-i]

		if len(enc) > 20 {
			enc = enc[:20] + "..."
		}
		if len(encBytes) > 20 {
			encBytes = encBytes[:20] + "..."
		}

		t.Errorf("encodings differ at %#q vs %#q", enc, encBytes)
	}
}

func TestRoundTripBytes(t *testing.T) {
	t.Parallel()

	// Test that UnquoteBytes reverses WriteQuotedBytes.
	var r []rune
	for i := '\u0000'; i <= unicode.MaxRune; i++ {
		r = append(r, i)
	}
	s := string(r)

	esBytes := &bytes.Buffer{}
	_, err := WriteQuotedBytes(esBytes, []byte(s))
	if err != nil {
		t.Fatalf("WriteQuotedBytes: %v", err)
	}

	decBytes, ok := UnquoteBytes(esBytes.Bytes())
	if !ok {
		t.Fatalf("UnquoteBytes not ok")
	}

	unqs := string(decBytes)
	if unqs != s {
		i := 0
		for i < len(unqs) && i < len(s) && unqs[i] == s[i] {
			i++
		}
		unqs = unqs[i:]
		s = s[i:]
		i = 0
		for i < len(unqs) && i < len(s) && unqs[len(unqs)-i-1] == s[len(s)-i-1] {
			i++
		}
		unqs = unqs[:len(unqs)-i]
		s = s[:len(s)-i]

		if len(unqs) > 20 {
			unqs = unqs[:20] + "..."
		}
		if len(s) > 20 {
			s = s[:20] + "..."
		}

		t.Errorf("decoding differs at %#q vs %#q", unqs, s)
	}
}
