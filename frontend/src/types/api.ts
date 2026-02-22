// API Type Definitions
// Single source of truth for all API-related TypeScript interfaces

// Link-related types
export interface LinkItem {
  id: string;
  title: string;
  url: string;
  description: string;
  tags: string[];
  created_at: string;  // Standardized field name
  view_count: number;
  updated?: string;    // Optional for backwards compatibility
}

export interface SearchResponse {
  links: LinkItem[];
  page: number;
  per_page: number;
  total_count: number;
  total_pages: number;
}

// Generic PocketBase response wrapper
export interface PocketBaseResponse<T = LinkItem> {
  page: number;
  perPage: number;
  totalItems: number;
  totalPages: number;
  items: T[];
}

// Stats types
export interface StatsData {
  total_links: number;
  total_tags: number;
  total_views: number;
  top_tags: Array<{
    name: string;
    slug: string;
    link_count: number;
  }>;
  most_viewed: Array<{
    id: string;
    title: string;
    url: string;
    view_count: number;
  }>;
}

export interface StatsState {
  loading: boolean;
  error: string | null;
  data: StatsData | null;
}

// Presentation types
export interface Presentation {
  id: string;
  name: string;
  step_count: number;
  step_labels: string[] | null;
  active_session: string | null;
  created_by: string | null;
  created: string;
  updated: string;
}

export interface PresentationStatus {
  id: string;
  name: string;
  step_count: number;
  step_labels: string[];
  is_live: boolean;
  progress: number | null;
  current_step: number | null;
}

// Search state for client-side search
export interface SearchState {
  query: string;
  isLoading: boolean;
  results: LinkItem[];
  totalCount: number;
  error: string | null;
}