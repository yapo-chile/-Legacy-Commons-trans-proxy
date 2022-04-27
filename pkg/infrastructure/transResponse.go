package infrastructure

import (
	"bytes"
	"fmt"
	"strconv"
)

// TransResponse a Trans response in bytes.
type TransResponse []byte

// Map returns a new map from a response.
func (r TransResponse) Map() (map[string]string, error) {
	m := make(map[string]string)
	err := r.apply(func(key, value string) {
		m[key] = value
	})
	return m, err
}

// apply applies the given function on all key-value pairs of the response.
func (r TransResponse) apply(f func(key, value string)) error {
	n := 0
	for n < len(r) {
		blobLen := 0
		// Check if the value is a blob.
		if len(r) > n+5 && bytes.Equal(r[n:n+5], []byte("blob:")) {
			i := bytes.IndexByte(r[n+5:], ':')
			if i == -1 {
				return fmt.Errorf("trans: invalid blob %q", r[n:])
			}
			n += 5
			var err error
			blobLen, err = strconv.Atoi(string(r[n : n+i]))
			if err != nil {
				return fmt.Errorf("trans: cannot parse blob length: %v", err)
			}
			n += i + 1
		}

		// if current field is blob field - key terminator is newline, not ':'
		var i int
		if blobLen > 0 {
			i = bytes.IndexByte(r[n:], '\n')
		} else {
			i = bytes.IndexByte(r[n:], ':')
		}

		if i == -1 {
			return fmt.Errorf("trans: invalid key-value format: %q", r[n:])
		}

		key := string(r[n : n+i])
		n += i + 1

		vl := n + blobLen
		// if current field is not blob field - read until newline, if there is a glob field,
		// we already have value length in blobLen variable
		if blobLen <= 0 {
			i = bytes.IndexByte(r[n:], '\n')
			if i == -1 {
				return fmt.Errorf("trans: newline is missing: %q", r[n:])
			}
			vl += i
		}

		f(key, string(r[n:vl]))
		n = vl + 1
	}
	return nil
}
