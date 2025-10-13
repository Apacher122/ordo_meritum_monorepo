package formatters

func PtrString(s *string, fallback string) string {
	if s == nil || *s == "" {
		return fallback
	}
	return *s
}

func PtrInt(i *int, fallback int) int {
	if i == nil {
		return fallback
	}
	return *i
}
