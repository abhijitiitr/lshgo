package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var len_buckets = 101 // A prime number(Since the number of bands I chose is 100)
type pair struct {
	rand1 int
	rand2 int
}

func InitializeArrayBuckets(num_of_bands int) [][]string {
	array_buckets := MakeTwoDArrayCustom(num_of_bands, len_buckets)
	return array_buckets
}

func HashMinhash(x, vari, cons, n int) (result int) {
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

func ConstructShingles(doc string, k int) (new_doc string, shingles []string) {
	shingles = make([]string, 1)
	new_doc = strings.ToLower(doc)
	substr_set := mapset.NewSet()
	new_doc = strings.Replace(new_doc, " ", "", -1)
	for i := 0; i < len(new_doc); i++ {
		if i+k < len(new_doc) {
			substr := new_doc[i:(i + k)]
			substr_set.Add(substr)
		}
	}
	data := substr_set.ToSlice()

	for i := range data {
		shingles = append(shingles, data[i].(string))
	}
	return new_doc, shingles
}

func ConstructSetShingles(docs []string, k int) (new_docs [][]string, shingles [][]string) {
	shingles = make([][]string, 1)
	for i := 0; i < len(docs); i++ {
		doc := docs[i]
		doc, sh := ConstructShingles(doc, k)
		docs[i] = doc
		shingles = append(shingles, sh)
	}
	return new_docs, shingles
}

func SortDocumentShingles(docs []string, shingles [][]string) (matrix [][]int) {
	matrix, rows := InitializeMatrix(docs, shingles)

	for col := 0; col < len(docs); col++ {
		for key, val := range rows {
			if strings.Contains(docs[col], key) {
				matrix[val][col] = 1
			}
		}
	}
	return matrix
}

func ComputeMinhashSignatures(matrix [][]int, n int) (signed_matrix [][]int) {
	hash_funcs := GenerateHashFuncs(n)
	hash_values := make([][]int, 1)
	val := make([]int, len(matrix))
	for _, hash_func := range hash_funcs {
		for i := 0; i < len(matrix); i++ {
			val[i] = HashMinhash(i, hash_func.rand1, hash_func.rand2, n)
		}
		hash_values = append(hash_values, val)
	}
	signed_matrix = MakeTwoDArray(n, len(matrix))
	signed_matrix = Add(signed_matrix, 100000)
	for i, row := range matrix {
		for j, _ := range row {
			if matrix[i][j] != 0 {
				for k := 0; k < n; k++ {
					hi := hash_values[k]
					signed_matrix[k][i] = Min(signed_matrix[k][i], hi[j])
				}
			}
		}
	}
	return signed_matrix
}
func EuclideanDistance(a, b []int) float64 {
	sum := ComputeSquaredDiffSum(a, b)
	return math.Sqrt(float64(sum))
}

func ComputeSquaredDiffSum(a, b []int) int {
	sum := 0
	for i, _ := range a {
		sum += (a[i] - b[i]) * (a[i] - b[i])
	}
	return sum
}

func ComputeDotProduct(a, b []int) int {
	sum := 0
	for i, _ := range a {
		sum += (a[i] * b[i])
	}
	return sum
}

func CosineDistance(a, b []int) float64 {
	prod_ab := ComputeDotProduct(a, b)
	zeros := make([]int, len(a))
	a_zero := EuclideanDistance(a, zeros)
	b_zero := EuclideanDistance(b, zeros)
	return float64(prod_ab) / (a_zero * b_zero)
}

func ApplyLsh(bands, rows int, signature_matrix [][]int) map[pair]float64 {
	array_buckets := InitializeArrayBuckets(bands)
	candidates := make(map[pair]float64)
	for i := 0; i < bands; i++ {
		buckets := array_buckets[i]
		band := signature_matrix[i:(i + rows)]
		for col := 0; col < len(signature_matrix[0]); col++ {
			sum := 0
			for k := 0; k < col; k++ {
				sum += band[k][col]
			}
			key := sum % len(buckets)
			buckets[key] += "_" + string(col)
		}
		i += rows
		for _, item := range buckets {
			item1 := strings.Split(item, "_")
			if len(item) > 1 {
				first_item, _ := strconv.Atoi(item1[0])
				second_item, _ := strconv.Atoi(item1[1])
				new_pair := pair{rand1: first_item, rand2: second_item}
				if candidates[new_pair] == 0 {
					first := signature_matrix[:][first_item]
					second := signature_matrix[:][second_item]
					similarity := CosineDistance(first, second)
					if similarity > 0.8 {
						candidates[new_pair] = similarity
					}
				}
			}
		}
	}
	return candidates
}

func Min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
func MakeTwoDArray(a, b int) (arr [][]int) {
	arr = make([][]int, a)
	e := make([]int, a*b)
	for i := range arr {
		arr[i] = e[i*b : (i+1)*b]
	}
	return arr
}

func MakeTwoDArrayCustom(a, b int) (arr [][]string) {
	arr = make([][]string, a)
	for i, val := range arr {
		for j, _ := range val {
			arr[i][j] = ""
		}
	}
	return arr
}

func Add(matrix [][]int, val int) [][]int {
	for i, row := range matrix {
		for j, col := range row {
			matrix[i][j] = col + val
		}
	}
	return matrix
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println(GenerateHashFuncs(10))
	fmt.Println(ConstructShingles("I am coding in Golang", 3))
}
