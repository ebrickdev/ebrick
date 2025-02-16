package middleware

import (
	"strconv"

	"github.com/ebrickdev/ebrick/transport/http"
)

const PagingParamsKey = "paging"

// PagingParams holds paging details.
type PagingParams struct {
	Page   int
	Limit  int
	Offset int
}

// PagingMiddleware extracts paging query params and stores them in the context.
func PagingMiddleware() http.HandlerFunc {
	return func(c *http.Context) {
		// Default values
		page := 1
		limit := 15

		if pageStr := c.Query("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l >= 0 {
				limit = l
			}
		}
		offset := (page - 1) * limit

		// Set paging parameters in the context
		c.Set(PagingParamsKey, &PagingParams{
			Page:   page,
			Limit:  limit,
			Offset: offset,
		})

		c.Next()
	}
}
