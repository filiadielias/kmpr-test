// Package general contains commonly use methods and variables
// such as config, server module, etc.
package general

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"io/ioutil"
)

// JSONUnmarshal : unmarshal json directly from Reader
func JSONUnmarshal(src io.Reader, dest interface{}) error {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &dest); err != nil {
		return err
	}

	return nil
}

// GobEncode : encode to bytes, used for convert struct into bytes (for NSQ message)
func GobEncode(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(data); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

// GobDecode : Decode bytes to struct / other data types, used for decoding NSQ message
func GobDecode(b []byte, data interface{}) error {

	d := gob.NewDecoder(bytes.NewReader(b))
	if err := d.Decode(data); err != nil {
		return err
	}

	return nil
}
