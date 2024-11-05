package search

type Base[CType any] struct {
	Elastic[CType]
	Open[CType]
}
