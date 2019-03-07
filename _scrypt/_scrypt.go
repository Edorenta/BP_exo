package _scrypt

import(
	// "fmt" //only for main test
	"time"
	"math/rand"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
)

//HashGet related:
const(
	// salt = sha256("NaCl")
	// salt = "DBCD6D34E827BCBDFCC06F0D7C6B54880D8F892701F81880AD319883EC6D6510"
	N = 32768	//CPU exp. difficulty
	r = 8		//memory exp. difficulty
	p = 1		//parallelization exp. difficulty
	l = 64		//output len 
)

//SaltGenerate generator related:
const(
	alpha_bytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    byte_bits = 6                    // 6 bits to represent a letter index
    byte_mask = 1 << byte_bits - 1 // All 1-bits, as many as byte_bits
    max_alpha_bit  = 63 / byte_bits   // # of letter indices fitting in 63 bits
)

func HashGet(s string, salt string) string {
	// salt := SaltGenerate(32)
	byte_hash, err := scrypt.Key([]byte(s), []byte(salt), N, r, p, l)
	if err != nil { return "" }
	return hex.EncodeToString(byte_hash[:l])//string(byte_hash[:l])
}

func HashMatch(s string, salt string, h string) bool {
	h2 := HashGet(s, salt)
	// fmt.Println("h1:", h, "h2:", h2)
    return h == h2
}

func PwdHash(pwd string) (string, string, error) {
	salt := SaltGenerate(32)
	hash, err := scrypt.Key([]byte(pwd), salt, N, r, p, l)
	if err != nil { return "", "", err }
	return hex.EncodeToString(hash), string(salt), nil
}

//unique rand seed for all salt generation
var src = rand.NewSource(time.Now().UnixNano())

//not goroutine safe, needs to instanciate an rand.Int63() per task for concurrent access
func SaltGenerate(n int) []byte {
    salt := make([]byte, n)
    // A src.Int63() generates 63 random bits, enough for max_alpha_bit characters!
    for i, cache, remain := n-1, src.Int63(), max_alpha_bit; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), max_alpha_bit
        }
        if j := int(cache & byte_mask); j < len(alpha_bytes) {
            salt[i] = alpha_bytes[j]
            i--
        }
        cache >>= byte_bits
        remain--
    }
    return salt
}
