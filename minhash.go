package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var len_buckets = 101 // A prime number(Since the number of bands I chose is 100)
type ar_int []int
type ar_ar_int []ar_int
type pair struct {
	rand1 int
	rand2 int
}

func InitializeArrayBuckets(a int) (b []ar_ar_int) {
	array_buckets := make([]ar_ar_int, (len_buckets * a))
	return array_buckets
}

func HashMinhash(x int, vari int, cons int, n int) (result int) {
	return (x*vari + cons) % len_buckets
}

func GenerateHashFuncs(num int) (arr_hash_funcs []pair) {
	num_of_funcs := make([]pair, num)
	for i := 0; i < num; i++ {
		num_of_funcs[i].rand1 = rand.Intn(num)
		num_of_funcs[i].rand2 = rand.Intn(num)
	}
	return num_of_funcs
}

func InitializeMatrix(docs []string, arr_shingles [][]string) (mat [][]int, rows map[string]int) {
	k := 0
	for _, shingles := range arr_shingles {
		for _, shingle := range shingles {
			if rows[shingle] == 0 {
				rows[shingle] = k
				k++
			}
		}
	}
	row := len(rows)
	col := len(docs)
	init_matrix := make([][]int, row)
	e := make([]int, row*col)
	for i := range init_matrix {
		init_matrix[i] = e[i*col : (i+1)*col]
	}

	return init_matrix, rows
}

func ConstructShingles(doc string, k int, h bool) (shingles []string) {
	shingles = make([]string, 1)

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println(GenerateHashFuncs(10))
}
