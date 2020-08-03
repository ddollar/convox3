package helpers

func Truncate(s string, max int, ellipsis bool) string {
	if len(s) > max {
		r := s[0:max]
		if ellipsis {
			r += "..."
		}
		return r
	}
	return s
}
