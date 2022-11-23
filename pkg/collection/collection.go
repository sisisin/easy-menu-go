package collection

func Map[T, V any](elms []T, fn func(T, int) V) []V {
	outputs := make([]V, len(elms), cap(elms))
	for i, elm := range elms {
		outputs[i] = fn(elm, i)
	}
	return outputs
}

func Filter[V any](elms []V, fn func(V, int) bool) []V {
	outputs := []V{}
	for i, elm := range elms {
		if fn(elm, i) {
			outputs = append(outputs, elm)
		}
	}
	return outputs
}

func FindIndex[V any](elms []V, fn func(V, int) bool) int {
	for i, v := range elms {
		if fn(v, i) {
			return i
		}
	}
	return -1
}
