package store

type Json interface {
	FromJson(rawJson []byte) error
	ToJson() ([]byte, error)
}
