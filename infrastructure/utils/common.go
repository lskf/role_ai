package utils

func DefaultOmitsColumn() []string {
	return []string{"id", "created_at", "creator", "uid", "account",
		"order_id", "status", "creator_name"}
}

func CustomDefaultOmitsColumn(customs ...string) []string {
	return append(DefaultOmitsColumn(), customs...)
}

func SliceContainer(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
