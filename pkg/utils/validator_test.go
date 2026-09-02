package utils

import (
	"strings"
	"testing"
)

type TestArticlePayload struct {
	Title    string `json:"title" validate:"required,min=20,max=200"`
	Content  string `json:"content" validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3,max=100"`
	Status   string `json:"status" validate:"required,oneof=publish draft thrash Publish Draft Thrash"`
}

func TestValidationRules(t *testing.T) {
	validContent := strings.Repeat("This is valid content for testing article post. ", 10)

	tests := []struct {
		name        string
		payload     TestArticlePayload
		expectError bool
		errorSubstr string
	}{
		{
			name: "Valid payload - lowercase status",
			payload: TestArticlePayload{
				Title:    "Valid Article Title with More than 20 chars",
				Content:  validContent,
				Category: "Technology",
				Status:   "publish",
			},
			expectError: false,
		},
		{
			name: "Valid payload - uppercase status",
			payload: TestArticlePayload{
				Title:    "Valid Article Title with More than 20 chars",
				Content:  validContent,
				Category: "Technology",
				Status:   "Draft",
			},
			expectError: false,
		},
		{
			name: "Invalid title - too short (<20 chars)",
			payload: TestArticlePayload{
				Title:    "Short Title",
				Content:  validContent,
				Category: "Technology",
				Status:   "publish",
			},
			expectError: true,
			errorSubstr: "title must be at least 20 characters",
		},
		{
			name: "Invalid content - too short (<200 chars)",
			payload: TestArticlePayload{
				Title:    "Valid Article Title with More than 20 chars",
				Content:  "Short content under 200 chars.",
				Category: "Technology",
				Status:   "publish",
			},
			expectError: true,
			errorSubstr: "content must be at least 200 characters",
		},
		{
			name: "Invalid category - too short (<3 chars)",
			payload: TestArticlePayload{
				Title:    "Valid Article Title with More than 20 chars",
				Content:  validContent,
				Category: "IT",
				Status:   "publish",
			},
			expectError: true,
			errorSubstr: "category must be at least 3 characters",
		},
		{
			name: "Invalid status - not in oneof",
			payload: TestArticlePayload{
				Title:    "Valid Article Title with More than 20 chars",
				Content:  validContent,
				Category: "Technology",
				Status:   "archived",
			},
			expectError: true,
			errorSubstr: "status must be one of",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(&tc.payload)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error containing '%s', got nil", tc.errorSubstr)
				}
				if !strings.Contains(err.Error(), tc.errorSubstr) {
					t.Errorf("expected error to contain '%s', got '%s'", tc.errorSubstr, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
			}
		})
	}
}
