/**
 * SearchEnhancer - Enhances search interface with client-side functionality
 * Extracted and simplified from components/SearchInterface.astro
 * Uses server-side navigation as requested instead of complex client-side rendering
 */

import type { SearchState } from '../types/api';

export class SearchEnhancer {
  private API_BASE: string;
  private DEBOUNCE_DELAY: number = 300;
  private FETCH_TIMEOUT: number = 5000;

  // State management
  private state: SearchState;
  private debounceTimeout: number | null = null;
  private currentController: AbortController | null = null;

  // DOM elements
  private searchInterface: HTMLElement;
  private searchForm: HTMLFormElement;
  private searchInput: HTMLInputElement;
  private loadingIndicator: HTMLElement;
  private errorDisplay: HTMLElement;

  constructor() {
    // Find required elements
    this.searchInterface = document.querySelector('.search-interface') as HTMLElement;
    this.searchForm = document.querySelector('.search-form') as HTMLFormElement;
    this.searchInput = document.querySelector('.search-input') as HTMLInputElement;
    this.loadingIndicator = document.querySelector('.loading-indicator') as HTMLElement;
    this.errorDisplay = document.querySelector('.error-display') as HTMLElement;

    if (!this.searchInterface || !this.searchForm || !this.searchInput || !this.loadingIndicator || !this.errorDisplay) {
      return;
    }

    // API configuration - try data attribute first
    const container = document.querySelector('[data-api-base]');
    this.API_BASE = container instanceof HTMLElement && container.dataset.apiBase
      ? container.dataset.apiBase
      : 'http://localhost:8090';

    // Initialize state
    this.state = {
      query: this.searchInput.value || '',
      isLoading: false,
      results: [],
      totalCount: 0,
      error: null
    };

    this.init();
  }

  private init(): void {
    this.setupEventListeners();
    this.setupPopstateHandler();
  }

  private setupEventListeners(): void {
    // Prevent default form submission and handle with JavaScript
    this.searchForm.addEventListener('submit', (e) => {
      e.preventDefault();
      const formData = new FormData(this.searchForm);
      const query = formData.get('q') as string || '';
      this.performSearch(query);
    });

    // Debounced input handling for live search
    this.searchInput.addEventListener('input', (e) => {
      const target = e.target as HTMLInputElement;
      const query = target.value.trim();

      // Clear existing debounce
      if (this.debounceTimeout) {
        clearTimeout(this.debounceTimeout);
      }

      // Set new debounce
      this.debounceTimeout = setTimeout(() => {
        if (query !== this.state.query) {
          this.performSearch(query);
        }
      }, this.DEBOUNCE_DELAY) as unknown as number;
    });

    // Error retry functionality
    const retryButton = this.errorDisplay?.querySelector('.error-retry');
    retryButton?.addEventListener('click', () => {
      this.performSearch(this.state.query);
    });
  }

  private setupPopstateHandler(): void {
    // Handle browser back/forward buttons
    window.addEventListener('popstate', (event) => {
      const url = new URL(window.location.href);
      const query = url.searchParams.get('q') || '';
      this.searchInput.value = query;

      if (event.state && event.state.query !== undefined) {
        // We have state, so this was a programmatic navigation
        if (!event.state.query) {
          window.location.reload(); // Go back to server-rendered home page
        }
      } else {
        // Browser navigation, reload to ensure consistency
        window.location.reload();
      }
    });
  }

  private async performSearch(query: string): Promise<void> {
    // Cancel any existing request
    if (this.currentController) {
      this.currentController.abort();
    }

    // Update state
    this.state.query = query;
    this.state.isLoading = true;
    this.state.error = null;

    // Update UI to show loading
    this.updateLoadingState(true);
    this.updateErrorState(false);

    // If empty query, use server-side navigation to home page
    if (!query.trim()) {
      window.history.pushState({}, '', '/');
      window.location.reload();
      return;
    }

    try {
      this.currentController = new AbortController();
      const timeoutId = setTimeout(() => this.currentController?.abort(), this.FETCH_TIMEOUT);

      // Use server-side navigation instead of complex client-side rendering
      // This ensures consistency with server-rendered components and simplifies the code
      const url = new URL(window.location.href);
      url.searchParams.set('q', query);

      // Update URL without page reload first
      window.history.pushState({ query }, '', url.toString());

      // Then navigate to the server-rendered search results page
      // This is simpler than trying to render components client-side
      window.location.href = url.toString();

    } catch (err) {
      this.state.isLoading = false;

      if (err && typeof err === 'object' && 'name' in err && err.name === 'AbortError') {
        // Request was cancelled, don't show error
        return;
      }

      this.state.error = err instanceof Error ? err.message : 'Search failed';
      this.updateLoadingState(false);
      this.updateErrorState(true, this.state.error);

      // Announce error to screen readers
      this.announceToScreenReader('Search failed. Please try again.');
    }
  }

  private updateLoadingState(isLoading: boolean): void {
    if (isLoading) {
      this.loadingIndicator?.setAttribute('aria-hidden', 'false');
      this.searchInterface?.setAttribute('data-loading', 'true');
    } else {
      this.loadingIndicator?.setAttribute('aria-hidden', 'true');
      this.searchInterface?.removeAttribute('data-loading');
    }
  }

  private updateErrorState(hasError: boolean, errorMessage?: string): void {
    if (hasError && errorMessage) {
      this.errorDisplay?.setAttribute('aria-hidden', 'false');
      const errorText = this.errorDisplay?.querySelector('.error-text');
      if (errorText) {
        errorText.textContent = errorMessage;
      }
    } else {
      this.errorDisplay?.setAttribute('aria-hidden', 'true');
    }
  }

  private announceToScreenReader(message: string): void {
    const announcement = document.createElement('div');
    announcement.setAttribute('aria-live', 'assertive');
    announcement.setAttribute('aria-atomic', 'true');
    announcement.style.position = 'absolute';
    announcement.style.left = '-10000px';
    announcement.style.width = '1px';
    announcement.style.height = '1px';
    announcement.style.overflow = 'hidden';
    announcement.textContent = message;

    document.body.appendChild(announcement);

    setTimeout(() => {
      document.body.removeChild(announcement);
    }, 1000);
  }
}