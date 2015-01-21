package minhash

func BitsSet(x uint64) (count int) {

	count = 0
	for ; x > 0; count++ {
		x &= x - 1
	}

	return count
}
func Hamming_distance(a uint64, b uint64) (distance uint8) {

	distance = 0

	var a_chunk, b_chunk uint8
	// a>0 || b>0 in this case a^b is sufficient since 0 means a=b=0 or a==b which makes their distance 0 :)
	for a^b != 0 {

		a_chunk = uint8((size - 1) & a)
		b_chunk = uint8((size - 1) & b)
		distance += hamming[a_chunk][b_chunk]
		a = a >> bit_length
		b = b >> bit_length
	}

	return
}

func Tokenize(in string, length int) (tokens []string) {

	num_tokens := (len(in)-1)/length + 1
	tokens = make([]string, num_tokens)

	i := 0
	for ; i < num_tokens-1; i++ {
		tokens[i] = in[i*length : (i+1)*length]
	}
	tokens[i] = in[i*length:]

	return tokens
}

func Tokenize_stride(in string, length int) (tokens []string) {

	if length > len(in) {
		return []string{in}
	}

	num_tokens := len(in) - length + 1
	tokens = make([]string, num_tokens)

	for i := 0; i < num_tokens; i++ {
		tokens[i] = in[i : i+length]
	}

	return tokens
}

func HammingDistance(a uint64, b uint64) int {

	x := a ^ b

	count := 0
	for ; x > 0; count++ {
		x &= x - 1
	}

	return count
}
