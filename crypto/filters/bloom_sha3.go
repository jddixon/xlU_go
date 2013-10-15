// A Bloom filter for sets of SHA3 digests.  A Bloom filter uses a set
// of k hash functions to determine set membership.  Each hash function
// produces a value in the range 0..M-1.  The filter is of size M.  To
// add a member to the set, apply each function to the new member and
// set the corresponding bit in the filter.  For M very large relative
// to k, this will normally set k bits in the filter.  To check whether
// x is a member of the set, apply each of the k hash functions to x
// and check whether the corresponding bits are set in the filter.  If
// any are not set, x is definitely not a member.  If all are set, x
// may be a member.  The probability of error (the false positive rate)
// is f = (1 - e^(-kN/M))^k, where N is the number of set members.
//
// This class takes advantage of the fact that SHA3 digests are good-
// quality pseudo-random numbers.  The k hash functions are the values
// of distinct sets of bits taken from the 20-byte SHA3 hash.  The
// number of bits in the filter, M, is constrained to be a power of
// 2; M == 2^m.  The number of bits in each hash function may not
// exceed floor(m/k).
//
// This class is designed to be thread-safe, but this has not been
// exhaustively tested.
package filters

import (
	// "code.google.com/p/go.crypto/sha3"
	"encoding/hex"
	"fmt"		// DEBUG
	"math"
	"sync"
)

type BloomSHA3 struct {
    m	uint		// protected final int m
    k	uint		// protected final int k
    count uint

	Filter	[]uint64
	ks		*KeySelector
	wordOffset	[]uint
	bitOffset	[]byte

    // convenience variables
    filterBits	int
    filterWords	int

	mu	sync.Mutex
}

 // Creates a filter with 2^m bits and k 'hash functions', where
 //each hash function is a portion of the 256-bit SHA3 hash.

 // @param m determines number of bits in filter, defaults to 20
 //  @param k number of hash functions, defaults to 8
func NewBloomSHA3( m, k uint) (b3 *BloomSHA3, err error) {

    // XXX need to devise more reasonable set of checks
    if  m < 2 || m > 20 {
		err = MOutOfRange
		// XXX what is this based on??
		if err == nil && ( k < 1 || ( k * m > 256 )) {
            // too many hash functions for filter size
			err = TooManyHashFunctions
		}
    }
	if err == nil {
		var ks *KeySelector

		filterBits := 1 << m
		filterWords :=  (filterBits + 31)/32     // round up
		b3 := &BloomSHA3 {
			m:	m,
			k:	k,
			filterBits : filterBits,
			filterWords : filterWords,
			Filter : make([]uint64, filterWords),
			wordOffset : make([]uint, k),
			bitOffset  : make([]byte, k),
		}
		b3. doClear()									// no lock
		// offsets into the filter
		ks, err = NewKeySelector(m, k, b3.bitOffset, b3.wordOffset)
		if err == nil {
			b3.ks = ks
		} else {
			b3 = nil 
		}
			
		// DEBUG
		fmt.Printf(
		"NewBloomSHA3: m = %d, k = %d, filterBits = %d, filterWords = %d\n", 
			m, k, filterBits, filterWords)
		// END
	}
	return
}

// Creates a filter of 2^m bits, with the number of 'hash functions"
// k defaulting to 8.
func NewNewBloomSHA3 (m uint) (*BloomSHA3, error) {
    return NewBloomSHA3 (m, 8)
}

// Creates a filter of 2^20 bits with k defaulting to 8.
// XXX Doubtful that this makes sense with 256 bit hash!
 
func NewNewNewBloomSHA3 () (*BloomSHA3, error){
    return NewBloomSHA3(20, 8)
}
// Clear the filter, unsynchronized 
func (b3 *BloomSHA3) doClear() {
    for i := 0; i < b3.filterWords; i++ {
        b3.Filter[i] = 0
    }
}
// Synchronized version */
func (b3 *BloomSHA3) Clear() {
	b3.mu.Lock()
    b3.doClear()
    b3.count = 0;          // jdd added 2005-02-19
    b3.mu.Unlock()
}
// Returns the number of keys which have been inserted.  This
// class (BloomSHA3) does not guarantee uniqueness in any sense; if the
// same key is added N times, the number of set members reported
// will increase by N.
func (b3 *BloomSHA3) Size() uint {
	b3.mu.Lock()
	defer b3.mu.Unlock()
    return b3.count
}

// Capacity returns the number of bits in the filter.
 
func (b3 *BloomSHA3)  Capacity() int  {
    return b3.filterBits
}

// Add a key to the set represented by the filter.
//
// XXX This version does not maintain 4-bit counters, it is not
// a counting Bloom filter.
func (b3 *BloomSHA3) Insert (b []byte) {
	b3.mu.Lock()
	defer b3.mu.Unlock()
    
	b3.ks.getOffsets(b)
    for i := uint(0); i < b3.k; i++ {
        b3.Filter[b3.wordOffset[i]] |=  1 << b3.bitOffset[i]
    }
    b3.count++
}

//
// Whether a key is in the filter.  Sets up the bit and word offset
// arrays.
//
// @param b byte array representing a key (SHA3 digest)
// @return true if b is in the filter
func (b3 *BloomSHA3) isMember(b []byte) bool {
    b3.ks.getOffsets(b)
    for i := uint(0); i < b3.k; i++ {
        if ! ((b3.Filter[b3.wordOffset[i]] & (1 << b3.bitOffset[i])) != 0)  {
            return false
        }
    }
    return true
}
// Whether a key is in the filter.  External interface, internally
// synchronized.
//
// @param b byte array representing a key (SHA3 digest)
// @return true if b is in the filter
func (b3 *BloomSHA3) Member(b []byte) bool {
	b3.mu.Lock()
	defer b3.mu.Unlock()
    
    return b3.isMember(b)
}

// For n the number of set members, return approximate false positive rate.
// XXX why two functions?? 
func (b3 *BloomSHA3) falsePositives(n uint) float64 {
    // (1 - e(-kN/M))^k

	fK := float64( b3.k )
	fN := float64(n)
	fB := float64(b3.filterBits)
    return math.Pow ( (1.0 - math.Exp( -fK * fN / fB)), fK)
}

func (b3 *BloomSHA3) FalsePositives() float64 {
    return b3.falsePositives(b3.count)
}
// DEBUG METHODS 
func  KeyToString(key []byte) string {
	return hex.EncodeToString(key)
}
// convert 64-bit integer to hex String */
func ltoh (i uint64) string {
	return fmt.Sprintf("#%x", i)
}

// convert 32-bit integer to String */
func itoh (i uint32) string {
	return fmt.Sprintf("#%x", i)
}

// convert single byte to String */
func btoh (b byte) string {
	return fmt.Sprintf("#%x", b)
}
