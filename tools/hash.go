package tools

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/fnv"
)

// SHA1 возвращает content-id вложения.
// sha1.Sum, а не общий hash.Hash: письма собираются параллельно,
// общее состояние хэшера давало гонку и неверные cid.
func SHA1(v []byte) string {
	if v == nil {
		return ""
	}
	sum := sha1.Sum(v)
	return hex.EncodeToString(sum[:])
}

// RandomID — случайная hex-строка (8 байт)
func RandomID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// HashLink — короткий хэш ссылки (fnv32), для имён скачанных картинок
func HashLink(link string) string {
	p := fnv.New32()
	p.Write([]byte(link))
	return fmt.Sprintf("%x", p.Sum(nil))
}
