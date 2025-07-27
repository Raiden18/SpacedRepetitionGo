package utils

func Map[T any, R any](input []T, f func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}

func Filter[T any](input []T, f func(T) bool) []T {
	result := []T{}
	for _, v := range input {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func ForEach[T any](input []T, f func(T)) {
	for _, v := range input {
		f(v)
	}
}
