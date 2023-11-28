package service

import (
	"github.com/stretchr/testify/assert"
	"noda/data/types"
	"testing"
)

func TestHelpers_doTrim(t *testing.T) {
	var trimmed = 0

	t.Run("success", func(t *testing.T) {
		var str1, str2 = "  \a\b\f\n\t\va  \a\b\f\n\t\v", "  \a\b\f\n\t\vb  \a\b\f\n\t\v"
		trimmed = doTrim(nil, nil, &str1, nil, nil, nil, &str2)
		assert.Equal(t, "a", str1)
		assert.Equal(t, "b", str2)
		assert.Equal(t, 2, trimmed)
	})

	t.Run("does nothing", func(t *testing.T) {
		trimmed = doTrim()
		assert.Equal(t, 0, trimmed)
	})

	t.Run("does nothing for nil parameters", func(t *testing.T) {
		trimmed = doTrim(nil, nil, nil, nil)
		assert.Equal(t, 0, trimmed)
	})

}

func TestHelpers_doDefaultPagination(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var pagination = &types.Pagination{Page: -1, RPP: 0}
		doDefaultPagination(pagination)
		assert.Equal(t, pagination.Page, int64(1))
		assert.Equal(t, pagination.RPP, int64(10))
		pagination.Page = 0
		pagination.RPP = -6
		doDefaultPagination(pagination)
		assert.Equal(t, pagination.Page, int64(1))
		assert.Equal(t, pagination.RPP, int64(10))
	})

	t.Run("does nothing for nil parameter", func(t *testing.T) {
		doDefaultPagination(nil)
	})
}
