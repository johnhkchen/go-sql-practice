// API Utilities and Configuration
// Shared utilities for making API calls across the frontend

// API Configuration Constants
export const API_BASE = import.meta.env.PUBLIC_API_URL || 'http://localhost:8090';
export const FETCH_TIMEOUT = 5000; // Default timeout for most operations
export const FETCH_TIMEOUT_LONG = 10000; // Longer timeout for complex operations

// Custom error class for API errors
export class ApiError extends Error {
  constructor(
    message: string,
    public code: 'timeout' | 'network' | 'server' | 'notfound',
    public status?: number
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

// Generic fetch helper with automatic timeout and error handling
export async function apiFetch<T>(
  url: string,
  options?: RequestInit & { timeout?: number }
): Promise<T> {
  const timeout = options?.timeout || FETCH_TIMEOUT;

  // Set up AbortController for timeout
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeout);

  try {
    const response = await fetch(url, {
      ...options,
      signal: controller.signal,
      headers: {
        'Accept': 'application/json',
        ...options?.headers,
      },
    });

    clearTimeout(timeoutId);

    // Handle HTTP errors
    if (!response.ok) {
      if (response.status === 404) {
        throw new ApiError(
          'Resource not found',
          'notfound',
          response.status
        );
      }
      throw new ApiError(
        `HTTP ${response.status}: ${response.statusText}`,
        'server',
        response.status
      );
    }

    // Parse and return JSON
    const data = await response.json();
    return data as T;

  } catch (error) {
    clearTimeout(timeoutId);

    // Re-throw ApiError instances
    if (error instanceof ApiError) {
      throw error;
    }

    // Handle abort/timeout
    if (error instanceof Error && error.name === 'AbortError') {
      throw new ApiError('Request timed out', 'timeout');
    }

    // Handle network errors
    throw new ApiError(
      error instanceof Error ? error.message : 'Network error occurred',
      'network'
    );
  }
}