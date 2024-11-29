package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// pair details ordering of an element.
type pair struct {
	index int
	value []byte
}

// compare decides the direction of the greater pair value.
func compare(a pair, b pair) int {
	for ii := range len(a.value) {
		switch {
		case len(b.value) < ii:
			return 1
		case a.value[ii] < b.value[ii]:
			return 1
		case a.value[ii] > b.value[ii]:
			return -1
		}
	}
	if len(a.value) < len(b.value) {
		return -1
	} else {
		return 0
	}
}

// calculate computes the final index of input element i.
func calculate(input []pair, i int, c chan pair) {
	a := input[i]
	count := 0
	for j, b := range input {
		diff := compare(a, b)
		if diff == -1 || diff == 0 && i > j {
			count += 1
		}
	}
	c <- pair{value: a.value, index: count}
}

// deorr does offset computation on a logical group of lines.
func deorr(input io.Reader, output io.Writer) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return err
	}
	pairs := []pair{}
	for i, v := range strings.Split(string(data), "\n") {
		if v == "" {
			continue
		}
		pair := pair{
			index: i,
			value: []byte(v + "\n"),
		}
		pairs = append(pairs, pair)
	}
	c := make(chan pair)
	for index := range pairs {
		go calculate(pairs, index, c)
	}
	sorted := make([]pair, len(pairs))
	for range pairs {
		pair := <-c
		sorted[pair.index] = pair
	}
	for _, v := range sorted {
		_, err := output.Write(v.value)
		if err != nil {
			return err
		}
	}
	return nil
}

// config contains processed customizations.
type config struct {
	input  *os.File
	output *os.File
}

// setup gathers values for the sorting processes.
func setup() (c config, err error) {
	flagInput := flag.String("input-file", "numbers.in", "unsorted records")
	flagOutput := flag.String("output-file", "numbers.out", "sorted destination")
	flag.Parse()
	c.input, err = os.Open(*flagInput)
	if err != nil {
		return config{}, err
	}
	c.output, err = os.Create(*flagOutput)
	if err != nil {
		c.input.Close()
		return config{}, err
	}
	return c, nil
}

// main orchestrates the command.
func main() {
	c, err := setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting setup: %+v\n", err)
		return
	}
	defer func() {
		c.input.Close()
		c.output.Close()
	}()
	err = deorr(c.input, c.output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sorting deorr: %+v\n", err)
		return
	}
}
