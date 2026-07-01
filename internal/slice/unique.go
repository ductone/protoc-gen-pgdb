package slice

func Unique[T comparable](slice []T) []T {
	ret := make([]T, 0, len(slice))
	dupeTrack := make(map[T]struct{})

	for _, o := range slice {
		if _, ok := dupeTrack[o]; ok {
			continue
		}

		dupeTrack[o] = struct{}{}
		ret = append(ret, o)
	}

	return ret
}
