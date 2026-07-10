package pkg

const (
	DefaultPage     = 1		// 第几页
	DefaultPageSize = 20	// 每页几条
	MaxPageSize     = 100	// 最多几条
)

// Page 规范分页参数
func Page(page, pageSize int) (int, int) {

	if page <= 0 {
		page = DefaultPage
	}

	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return page, pageSize
}

// Offset 返回分页偏移量
func Offset(page, pageSize int) int {

	page, pageSize = Page(page, pageSize)

	return (page - 1) * pageSize
}