package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"slices"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{
		"deorr",
		"-input-file",
		"records/numbers-10.in",
		"-output-file",
		"records/numbers-10.out",
	}
	main()
	output, err := os.Open("records/numbers-10.out")
	if err != nil {
		t.Logf("Failed to open the output results")
		t.FailNow()
	}
	actual, err := io.ReadAll(output)
	if err != nil {
		t.Logf("Failed to read the actual results")
		t.FailNow()
	}
	expected, err := os.Open("records/numbers-10.expected")
	if err != nil {
		t.Logf("Failed to open the expected results")
		t.FailNow()
	}
	sorted, err := io.ReadAll(expected)
	if err != nil {
		t.Logf("Failed to open the sorted results")
		t.FailNow()
	}
	if string(actual) != string(sorted) {
		t.Logf("Expected:\n%s", string(actual))
		t.Logf("Actual:\n%s", string(sorted))
		t.FailNow()
	}
}

func BenchmarkMain(b *testing.B) {
	tests := map[string]struct {
		n int
	}{
		"100": {
			n: 100,
		},
		"1000": {
			n: 1000,
		},
		"10000": {
			n: 10000,
		},
		"100000": {
			n: 100000,
		},
	}
	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				generated := generate(tt.n)
				input := bytes.Buffer{}
				tee := io.TeeReader(generated, &input)
				blob, err := io.ReadAll(tee)
				if err != nil {
					b.Log("Failed to generate the input")
					b.FailNow()
				}
				lines := strings.Split(string(blob), "\n")
				slices.Sort(lines)
				expected := strings.Join(lines, "\n")
				output := bytes.Buffer{}
				b.StartTimer()
				err = deorr(&input, &output)
				if err != nil {
					b.Log("Failed to perform the sort")
					b.FailNow()
				}
				b.StopTimer()
				actual := strings.Join(lines, "\n")
				if actual != expected {
					b.Log("Failed to sort in an order")
					b.FailNow()
				}
			}
		})
	}
}

// generateASCII generates a random ASCII sequence.
func generateASCII(length int) string {
	const asciiChars = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	result := strings.Builder{}
	for range length {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(asciiChars))))
		result.WriteByte(asciiChars[num.Int64()])
	}
	return result.String()
}

// generateBase16 generates a random base16 sequence.
func generateBase16(length int) string {
	const hexChars = "0123456789ABCDEF"
	result := strings.Builder{}
	for range length {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(hexChars))))
		result.WriteByte(hexChars[num.Int64()])
	}
	return result.String()
}

// generate generates n random records to sort.
func generate(n int) *bytes.Buffer {
	buff := bytes.Buffer{}
	for i := range n {
		_, _ = buff.WriteString(
			fmt.Sprintf(
				"%s  %s  %s\r\n",
				generateASCII(10),
				fmt.Sprintf("%032d", i),
				generateBase16(64),
			),
		)
	}
	return &buff
}
