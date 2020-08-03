package helpers

func CoalesceString(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}

	return ""
}

func DefaultString(s *string, def string) string {
	if s == nil {
		return def
	}

	return *s
}

func StringSliceContains(ss []string, s string) bool {
	for _, i := range ss {
		if i == s {
			return true
		}
	}

	return false
}
