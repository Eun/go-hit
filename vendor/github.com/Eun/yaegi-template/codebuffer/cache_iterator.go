package codebuffer

type cacheIterator struct {
	parts        []*Part
	size         int
	currentIndex int
}

func newCacheIterator(parts []*Part) (Iterator, error) {
	return &cacheIterator{
		parts:        parts,
		size:         len(parts),
		currentIndex: -1,
	}, nil
}

func (i *cacheIterator) Next() bool {
	i.currentIndex++
	return i.currentIndex < i.size
}

func (i *cacheIterator) Value() *Part {
	if i.currentIndex < 0 || i.currentIndex >= i.size {
		return nil
	}
	return i.parts[i.currentIndex]
}

func (i *cacheIterator) Error() error {
	return nil
}
