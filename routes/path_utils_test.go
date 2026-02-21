package routes

import "testing"

func TestExtractPathParam_ValidPaths(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		targetSegment  string
		expected       string
	}{
		{
			name:          "presentations status endpoint",
			path:          "/api/presentations/abc123/status",
			targetSegment: "presentations",
			expected:      "abc123",
		},
		{
			name:          "presentations live endpoint",
			path:          "/api/presentations/def456/live",
			targetSegment: "presentations",
			expected:      "def456",
		},
		{
			name:          "presentations stop endpoint",
			path:          "/api/presentations/ghi789/stop",
			targetSegment: "presentations",
			expected:      "ghi789",
		},
		{
			name:          "sync progress endpoint",
			path:          "/api/sync/session123/progress",
			targetSegment: "sync",
			expected:      "session123",
		},
		{
			name:          "links view endpoint",
			path:          "/api/links/link456/view",
			targetSegment: "links",
			expected:      "link456",
		},
		{
			name:          "uuid format parameter",
			path:          "/api/presentations/550e8400-e29b-41d4-a716-446655440000/status",
			targetSegment: "presentations",
			expected:      "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPathParam(tt.path, tt.targetSegment)
			if result != tt.expected {
				t.Errorf("extractPathParam(%q, %q) = %q, want %q",
					tt.path, tt.targetSegment, result, tt.expected)
			}
		})
	}
}

func TestExtractPathParam_InvalidPaths(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		targetSegment  string
		expected       string
	}{
		{
			name:          "empty path",
			path:          "",
			targetSegment: "presentations",
			expected:      "",
		},
		{
			name:          "empty target segment",
			path:          "/api/presentations/abc123/status",
			targetSegment: "",
			expected:      "",
		},
		{
			name:          "target segment not found",
			path:          "/api/presentations/abc123/status",
			targetSegment: "nonexistent",
			expected:      "",
		},
		{
			name:          "target segment at end of path",
			path:          "/api/presentations",
			targetSegment: "presentations",
			expected:      "",
		},
		{
			name:          "empty parameter after target",
			path:          "/api/presentations//status",
			targetSegment: "presentations",
			expected:      "",
		},
		{
			name:          "parameter is just slash",
			path:          "/api/presentations/status",
			targetSegment: "presentations",
			expected:      "status",  // This is actually valid - status becomes the parameter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPathParam(tt.path, tt.targetSegment)
			if result != tt.expected {
				t.Errorf("extractPathParam(%q, %q) = %q, want %q",
					tt.path, tt.targetSegment, result, tt.expected)
			}
		})
	}
}

func TestExtractPathParam_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		targetSegment  string
		expected       string
	}{
		{
			name:          "multiple slashes in path",
			path:          "/api//presentations///abc123//status",
			targetSegment: "presentations",
			expected:      "abc123",
		},
		{
			name:          "path without leading slash",
			path:          "api/presentations/abc123/status",
			targetSegment: "presentations",
			expected:      "abc123",
		},
		{
			name:          "path with trailing slash",
			path:          "/api/presentations/abc123/status/",
			targetSegment: "presentations",
			expected:      "abc123",
		},
		{
			name:          "repeated target segment",
			path:          "/api/presentations/presentations/abc123",
			targetSegment: "presentations",
			expected:      "presentations",  // Returns first match
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPathParam(tt.path, tt.targetSegment)
			if result != tt.expected {
				t.Errorf("extractPathParam(%q, %q) = %q, want %q",
					tt.path, tt.targetSegment, result, tt.expected)
			}
		})
	}
}

func TestExtractLastPathParam_ValidPaths(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		endingSegment  string
		expected       string
	}{
		{
			name:          "presentations status endpoint",
			path:          "/api/presentations/abc123/status",
			endingSegment: "status",
			expected:      "abc123",
		},
		{
			name:          "sync progress endpoint",
			path:          "/api/sync/session456/progress",
			endingSegment: "progress",
			expected:      "session456",
		},
		{
			name:          "links view endpoint",
			path:          "/api/links/link789/view",
			endingSegment: "view",
			expected:      "link789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractLastPathParam(tt.path, tt.endingSegment)
			if result != tt.expected {
				t.Errorf("extractLastPathParam(%q, %q) = %q, want %q",
					tt.path, tt.endingSegment, result, tt.expected)
			}
		})
	}
}

func TestExtractLastPathParam_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		endingSegment  string
		expected       string
	}{
		{
			name:          "empty path",
			path:          "",
			endingSegment: "status",
			expected:      "",
		},
		{
			name:          "empty ending segment",
			path:          "/api/presentations/abc123/status",
			endingSegment: "",
			expected:      "",
		},
		{
			name:          "ending segment not found",
			path:          "/api/presentations/abc123/status",
			endingSegment: "nonexistent",
			expected:      "",
		},
		{
			name:          "ending segment at beginning",
			path:          "/status/abc123",
			endingSegment: "status",
			expected:      "",
		},
		{
			name:          "empty parameter before ending",
			path:          "/api/presentations//status",
			endingSegment: "status",
			expected:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractLastPathParam(tt.path, tt.endingSegment)
			if result != tt.expected {
				t.Errorf("extractLastPathParam(%q, %q) = %q, want %q",
					tt.path, tt.endingSegment, result, tt.expected)
			}
		})
	}
}