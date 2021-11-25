package pagination

import (
	"math"

	"github.com/rayyone/go-core/helpers/method"
)

type Config struct {
	Page   int
	Offset int
	Limit  int
}

type Paginator struct {
	TotalItems   int  `json:"total_items"`
	TotalPages   int  `json:"total_pages"`
	ItemFrom     int  `json:"item_from"`
	ItemTo       int  `json:"item_to"`
	CurrentPage  int  `json:"current_page"`
	Limit        int  `json:"limit"`
	NextPage     *int `json:"next_page"`
	PreviousPage *int `json:"previous_page"`
}

func GetPaginationConfig(page int, limit int) Config {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 30
	}

	offset := (page - 1) * limit

	return Config{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func BuildPaginator(total int, limit int, offset int) *Paginator {
	var paginator Paginator

	currentPage := offset/limit + 1
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	itemFrom := validateMaxNumber((currentPage-1)*limit+1, total)
	itemTo := validateMaxNumber(currentPage*limit, total)

	paginator.TotalItems = total
	paginator.TotalPages = totalPages
	paginator.ItemFrom = itemFrom
	paginator.ItemTo = itemTo
	paginator.CurrentPage = currentPage
	paginator.Limit = limit
	if currentPage < totalPages {
		paginator.NextPage = method.NewInt(currentPage + 1)
	}
	if currentPage > 1 {
		paginator.PreviousPage = method.NewInt(currentPage - 1)
	}

	return &paginator
}

func validateMaxNumber(number int, max int) int {
	if number < 0 {
		return 0
	}
	if number < max {
		return number
	}
	return max
}
