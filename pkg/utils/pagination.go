package utils

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

func ValidatePagination(page, pageSize int) (int, int, error) {
	if page < 1 {
		page = DefaultPage
	}

	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return page, pageSize, nil
}
