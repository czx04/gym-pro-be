package helper

func SetIfNotNil[T comparable](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}
