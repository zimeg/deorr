# > deorr

This sort hopes to use a parallel sorting algorithm to go fast. Ideally `O(n)`.

**Outline**:

- [Sorting in parallel](#sorting-in-parallel)
  - [The mathematical foundations](#the-mathematical-foundations)
  - [The algorithmic computations](#the-algorithmic-computations)
  - [The concurrent twist](#the-concurrent-twist)
- [Benchmarking](#benchmarking)

## Sorting in parallel

Ordering in linear time takes a different approach to calculating the sorted
positions of elements in a list, and takes advantage of concurrent computation.

### The mathematical foundations

Remark that in a sorted list, the index of each element is equal to the number
of elements that value is greater than or sometimes equal to.

Noted are indices in a sorted list:

```
[0, 1, 1, 2, 3, 5, 8]
       ^     ^
       i=2   i=4
```

Calculating the index, or count, for a unique value only requires the greater
than comparison, shown above as `3` being greater than 4 other elements.

Counts for duplicate values must consider the intial order of the list in the
calculation. If elements are equal, the latter value should have a higher count.

The following example demonstrates this:

```
[3, 8, 1, 0, 2, 5, 1]
       ^           ^
       i=1         i=2
```

### The algorithmic computations

This comparison method can be translated into code by keeping track of these
counts to permute the input as follows:

```go
func main() {
	input := []int{2, 0, 4, 3, 3}
	count := []int{0, 0, 0, 0, 0}
	sorts := []int{0, 0, 0, 0, 0}

	// i is the index, a is the value
	for i, a := range input {
		for j, b := range input {
			if a > b || (a == b && i > j) {
				count[i] += 1
			}
		}
	}

	for i, k := range count {
		sorts[k] = input[i]
	}

	fmt.Println(input) // [2, 0, 4, 3, 3]
	fmt.Println(count) // [1, 0, 4, 2, 3]
	fmt.Println(sorts) // [0, 2, 3, 3, 4]
}
```

The `count` calculates the number of elements each `input` value is greater
than, which is equivalent to knowing the index of a value in the sorted `sorts`.

As is, this is `O(n^2)` with the nested loop over `input`, which isn't great.

### The concurrent twist

However, notice that the inside loop does not depend on calculations from other
iterations! That means it can be done in parallel, with the results being joined
together when finished:

```go
type Pair struct {
	index int
	value int
}

func calculate(input []int, i int, c chan Pair) {
	a := input[i]
	count := 0
	for j, b := range input {
		if a > b || (a == b && i > j) {
			count += 1
		}
	}
	c <- Pair{value: a, index: count}
}

func main() {
	input := []int{2, 0, 4, 3, 3}
	sorts := []int{0, 0, 0, 0, 0}

	c := make(chan Pair)
	for i, _ := range input {
		go calculate(input, i, c) // Perform measurements in a thread
	}

	for _, _ = range input {
		pair := <-c
		sorts[pair.index] = pair.value
	}

	fmt.Println(input) // [2, 0, 4, 3, 3]
	fmt.Println(sorts) // [0, 2, 3, 3, 4]
}
```

Here, the `main` function is using a channel `c` for receiving `Pair` values
from the `calculate` function. That same `calculate` function is being invoked
multiple times as a goroutine - a slick way to run code concurrent.

In theory, having `len(input)` threads allows the `calculate` executions to run
at the same time, iterating over the inner `input` loop in a parallel but offset
manner. In practice, that many threads is infrequently found.

Both of these examples can be explored on the Go Playground - either
[sequentially][go_1] or [concurrently][go_2].

## Benchmarking

These benchmarks are done with the rules of [the sort benchmark][tsb] in mind,
with the measurements being made for a [JouleSort][joulesort] since energy is
interesting and storing 100TB for the other sorts seems challenging.

[go_1]: https://go.dev/play/p/HjtS35PaHdJ
[go_2]: https://go.dev/play/p/dKi9P5Z_nmx
[joulesort]: http://csl.stanford.edu/~christos/publications/2007.jsort.sigmod.pdf
[tsb]: http://sortbenchmark.org
