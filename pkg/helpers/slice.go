package helpers

func SliceAdd(ss []string, s string) []string {
	for _, x := range ss {
		if x == s {
			return ss
		}
	}

	ss = append(ss, s)

	return ss
}

func SliceContains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}

	return false
}

func SliceRemove(ss []string, s string) []string {
	for i, x := range ss {
		if x == s {
			ss = append(ss[0:i], ss[i+1:]...)
		}
	}

	return ss
}
