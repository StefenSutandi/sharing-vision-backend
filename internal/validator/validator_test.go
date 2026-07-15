package validator_test

import (
	"github.com/StefenSutandi/sharing-vision-backend/internal/dto"
	"github.com/StefenSutandi/sharing-vision-backend/internal/validator"
	"strings"
	"testing"
)

func TestArticleValidator(t *testing.T) {
	validContent := strings.Repeat("a", 200)
	tooLongTitle := strings.Repeat("a", 201)

	tests := []struct {
		name    string
		payload dto.CreateArticleReq
		valid   bool
	}{
		{
			name: "valid publish",
			payload: dto.CreateArticleReq{
				Title:    "Valid title minimum twenty",
				Content:  validContent,
				Category: "Tech",
				Status:   "publish",
			},
			valid: true,
		},
		{
			name: "valid draft",
			payload: dto.CreateArticleReq{
				Title:    "Valid title minimum twenty",
				Content:  validContent,
				Category: "Tech",
				Status:   "draft",
			},
			valid: true,
		},
		{
			name: "valid thrash",
			payload: dto.CreateArticleReq{
				Title:    "Valid title minimum twenty",
				Content:  validContent,
				Category: "Tech",
				Status:   "thrash",
			},
			valid: true,
		},
		{
			name: "empty title",
			payload: dto.CreateArticleReq{
				Title:    "",
				Content:  validContent,
				Category: "Tech",
				Status:   "publish",
			},
			valid: false,
		},
		{
			name: "title below 20 chars",
			payload: dto.CreateArticleReq{
				Title:    "too short",
				Content:  validContent,
				Category: "Tech",
				Status:   "publish",
			},
			valid: false,
		},
		{
			name: "title above 200 chars",
			payload: dto.CreateArticleReq{
				Title:    tooLongTitle,
				Content:  validContent,
				Category: "Tech",
				Status:   "publish",
			},
			valid: false,
		},
		{
			name: "empty content",
			payload: dto.CreateArticleReq{
				Title:    "Valid title minimum twenty",
				Content:  "",
				Category: "Tech",
				Status:   "publish",
			},
			valid: false,
		},
		{
			name: "content below 200 chars",
			payload: dto.CreateArticleReq{
				Title:    "Valid title minimum twenty",
				Content:  "too short",
				Category: "Tech",
				Status:   "publish",
			},
			valid: false,
		},
		{
			name: "empty category",
			payload: dto.CreateArticleReq{
				Title:    "Valid title minimum twenty",
				Content:  validContent,
				Category: "",
				Status:   "publish",
			},
			valid: false,
		},
		{
			name: "category below 3 chars",
			payload: dto.CreateArticleReq{
				Title:    "Valid title minimum twenty",
				Content:  validContent,
				Category: "Te",
				Status:   "publish",
			},
			valid: false,
		},
		{
			name: "invalid status",
			payload: dto.CreateArticleReq{
				Title:    "Valid title minimum twenty",
				Content:  validContent,
				Category: "Tech",
				Status:   "invalid",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validator.ValidateStruct(tt.payload)
			if tt.valid && len(errs) > 0 {
				t.Errorf("expected valid, got error: %v", errs)
			}
			if !tt.valid && len(errs) == 0 {
				t.Errorf("expected error, got valid")
			}
		})
	}
}
