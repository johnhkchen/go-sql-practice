package routes

import "strings"

// extractPathParam extracts a parameter value from a URL path by finding the segment
// that follows the specified target segment.
//
// Example:
//   extractPathParam("/api/presentations/abc123/status", "presentations") returns "abc123"
//   extractPathParam("/api/sync/def456/progress", "sync") returns "def456"
//   extractPathParam("/api/links/ghi789/view", "links") returns "ghi789"
//
// Returns empty string if:
//   - path is empty
//   - target segment is not found
//   - no segment follows the target segment
//   - the following segment is empty
func extractPathParam(path, targetSegment string) string {
	if path == "" || targetSegment == "" {
		return ""
	}

	// Split path by "/" and filter out empty segments
	parts := strings.Split(path, "/")

	// Find the target segment and return the next non-empty segment
	for i, part := range parts {
		if part == targetSegment && i+1 < len(parts) {
			param := parts[i+1]
			if param != "" {
				return param
			}
		}
	}

	return ""
}

// extractLastPathParam extracts the parameter that appears before a specified ending segment.
// This is useful for URLs where the parameter is the second-to-last segment.
//
// Example:
//   extractLastPathParam("/api/sync/abc123/progress", "progress") returns "abc123"
//   extractLastPathParam("/api/presentations/def456/status", "status") returns "def456"
//   extractLastPathParam("/api/links/ghi789/view", "view") returns "ghi789"
//
// Returns empty string if:
//   - path is empty
//   - ending segment is not found
//   - no segment precedes the ending segment
//   - the preceding segment is empty
func extractLastPathParam(path, endingSegment string) string {
	if path == "" || endingSegment == "" {
		return ""
	}

	// Split path by "/" and filter out empty segments
	parts := strings.Split(path, "/")

	// Find the ending segment and return the previous non-empty segment
	for i, part := range parts {
		if part == endingSegment && i > 0 {
			param := parts[i-1]
			if param != "" {
				return param
			}
		}
	}

	return ""
}