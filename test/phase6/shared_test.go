package phase6_test

import (
	"akadia/internal/shared"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPageableNormalizeAndOffset(t *testing.T) {
	tests := []struct {
		name           string
		pageable       shared.Pageable
		expectedPage   int
		expectedSize   int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "invalid page and size fall back to defaults",
			pageable:       shared.Pageable{Page: 0, Size: 0},
			expectedPage:   1,
			expectedSize:   10,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "size is capped at maximum",
			pageable:       shared.Pageable{Page: 3, Size: 999},
			expectedPage:   3,
			expectedSize:   100,
			expectedLimit:  100,
			expectedOffset: 200,
		},
		{
			name:           "valid values are preserved",
			pageable:       shared.Pageable{Page: 2, Size: 15},
			expectedPage:   2,
			expectedSize:   15,
			expectedLimit:  15,
			expectedOffset: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageable := tt.pageable
			pageable.Normalize()

			if pageable.Page != tt.expectedPage {
				t.Fatalf("expected page %d, got %d", tt.expectedPage, pageable.Page)
			}
			if pageable.Size != tt.expectedSize {
				t.Fatalf("expected size %d, got %d", tt.expectedSize, pageable.Size)
			}
			if pageable.Limit() != tt.expectedLimit {
				t.Fatalf("expected limit %d, got %d", tt.expectedLimit, pageable.Limit())
			}
			if pageable.Offset() != tt.expectedOffset {
				t.Fatalf("expected offset %d, got %d", tt.expectedOffset, pageable.Offset())
			}
		})
	}
}

func TestGetPageableBindsAndNormalizesQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("GET", "/payment-policies?page=-5&size=101", nil)

	pageable, err := shared.GetPageable(c)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pageable.Page != 1 {
		t.Fatalf("expected normalized page 1, got %d", pageable.Page)
	}
	if pageable.Size != 100 {
		t.Fatalf("expected normalized size 100, got %d", pageable.Size)
	}
}

func TestNewPageAndNewPageEmptyBuildPaginationMetadata(t *testing.T) {
	pageable := &shared.Pageable{Page: 2, Size: 10}
	page := shared.NewPage(pageable, 25, []string{"a", "b"})

	if len(page.Data) != 2 {
		t.Fatalf("expected two data rows, got %d", len(page.Data))
	}
	if page.Pagination.TotalCount != 25 {
		t.Fatalf("expected total count 25, got %d", page.Pagination.TotalCount)
	}
	if page.Pagination.TotalPage != 3 {
		t.Fatalf("expected total pages 3, got %d", page.Pagination.TotalPage)
	}
	if page.Pagination.PreviousPage == nil || *page.Pagination.PreviousPage != 1 {
		t.Fatalf("expected previous page 1")
	}
	if page.Pagination.NextPage == nil || *page.Pagination.NextPage != 3 {
		t.Fatalf("expected next page 3")
	}

	empty := shared.NewPageEmpty[string](pageable)
	if len(empty.Data) != 0 {
		t.Fatalf("expected empty page data")
	}
	if empty.Pagination.TotalCount != 0 || empty.Pagination.TotalPage != 0 {
		t.Fatalf("expected empty pagination totals to be zero")
	}
	if empty.Pagination.PreviousPage != nil || empty.Pagination.NextPage != nil {
		t.Fatalf("expected empty page to omit previous and next page")
	}
}
