// Package rndbytes is a fast pseudo random []byte generator. It uses a rand.Source64
// together with sync.Mutex which is safe for concurrent use by multiple goroutines.
// The source is seeded at package init.
// For random numbers suitable for security-sensitive work, see the crypto/rand package.
package rndbytes

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 64 / letterIdxBits   // # of letter indices fitting in 64 bits
)

type bytesSource struct {
	lk       sync.Mutex
	rng      rand.Source64
	alphaNum string
}

var src *bytesSource

func (s *bytesSource) uint64() uint64 {
	s.lk.Lock()
	n := s.rng.Uint64()
	s.lk.Unlock()
	return n
}

func init() {

	// https://stackoverflow.com/questions/12321133/how-to-properly-seed-random-number-generator#54491783
	var b [8]byte
	_, err := cryptorand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}

	data := []byte("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_")
	dataLen := len(data)
	// shuffle is not really needed, but does not hurt either
	rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:])))).Shuffle(
		dataLen, func(i, j int) {
			data[i], data[j] = data[j], data[i]
		})

	_, err = cryptorand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	src = &bytesSource{
		rng:      rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))).(rand.Source64),
		alphaNum: string(data),
	}
}

// GetBytes returns []byte of n length filled with random characters from
// "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_" set.
// When allowFirstDash param is false the first byte is guaranteed NOT to be a '-' (dash)
func GetBytes(n int, allowFirstDash bool) []byte {
	b := make([]byte, n)

	for i, cache, remain := n-1, src.uint64(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.uint64(), letterIdxMax
		}
		idx := cache & letterIdxMask
		b[i] = src.alphaNum[idx]
		i--
		cache >>= letterIdxBits
		remain--
	}

	if allowFirstDash == false && b[0] == '-' {
		cache := src.uint64()
		idx := cache & letterIdxMask
		b[0] = src.alphaNum[idx]
		if b[0] == '-' {
			delta := (cache >> letterIdxBits) & letterIdxMask
			if delta == 0 {
				delta = 1
			}
			b[0] = src.alphaNum[(idx+delta)&letterIdxMask]
		}
	}

	return b
}

// GetInt returns a random int from seeded rand.Source64
func GetInt() int {
	return int(src.uint64())
}
