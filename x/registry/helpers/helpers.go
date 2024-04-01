package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"sort"

	wnsUtils "git.vdb.to/cerc-io/laconicd/utils"
	set "github.com/deckarep/golang-set"
)

func StringToBytes(val string) []byte {
	return []byte(val)
}

func BytesToString(val []byte) string {
	return string(val)
}

func StrArrToBytesArr(val []string) ([]byte, error) {
	buffer := &bytes.Buffer{}

	err := gob.NewEncoder(buffer).Encode(val)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func BytesArrToStringArr(val []byte) ([]string, error) {
	buffer := bytes.NewReader(val)
	var v []string
	err := gob.NewDecoder(buffer).Decode(&v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func Int64ToBytes(num int64) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, num)
	return buf.Bytes()
}

func MustMarshalJSON[T any](val T) (bytes []byte) {
	bytes, err := json.Marshal(val)
	if err != nil {
		panic("JSON marshal error:" + err.Error())
	}
	return
}

func MustUnmarshalJSON[T any](bytes []byte) T {
	var val T
	err := json.Unmarshal(bytes, &val)
	if err != nil {
		panic("JSON unmarshal error:" + err.Error())
	}
	return val
}

// GetCid gets the content ID.
func GetCid(content []byte) (string, error) {
	return wnsUtils.CIDFromJSONBytes(content)
}

// BytesToBase64 encodes a byte array as a base64 string.
func BytesToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// BytesFromBase64 decodes a byte array from a base64 string.
func BytesFromBase64(str string) []byte {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic("Error decoding string to bytes.")
	}

	return bytes
}

// BytesToHex encodes a byte array as a hex string.
func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// BytesFromHex decodes a byte array from a hex string.
func BytesFromHex(str string) []byte {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic("Error decoding hex to bytes.")
	}

	return bytes
}

func SetToSlice(set set.Set) []string {
	names := []string{}

	for name := range set.Iter() {
		if name, ok := name.(string); ok && name != "" {
			names = append(names, name)
		}
	}

	sort.SliceStable(names, func(i, j int) bool { return names[i] < names[j] })

	return names
}

func SliceToSet(names []string) set.Set {
	set := set.NewThreadUnsafeSet()

	for _, name := range names {
		if name != "" {
			set.Add(name)
		}
	}

	return set
}

func AppendUnique(list []string, element string) []string {
	set := SliceToSet(list)
	set.Add(element)

	return SetToSlice(set)
}
