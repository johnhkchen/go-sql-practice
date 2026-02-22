/**
 * StatsController - Manages statistics display and interactions
 * Extracted from components/StatsSummary.astro
 */

import type { StatsData, StatsState } from '../types/api';

// API URL resolution helper
function getApiBase(): string {
  if (typeof import.meta !== 'undefined' && import.meta.env?.PUBLIC_API_URL) {
    return import.meta.env.PUBLIC_API_URL;
  }
  // Check for data-api-base attribute on container
  const container = document.querySelector('[data-api-base]');
  if (container instanceof HTMLElement && container.dataset.apiBase) {
    return container.dataset.apiBase;
  }
  return 'http://localhost:8090';
}

export class StatsController {
  private state: StatsState;
  private container: HTMLElement;

  constructor() {
    this.state = {
      loading: true,
      error: null,
      data: null
    };

    const container = document.getElementById('stats-container');
    if (!container) {
      console.error('Stats container not found');
      throw new Error('Stats container not found');
    }

    this.container = container;
    this.init();
  }

  private async init(): Promise<void> {
    await this.fetchStats();
    this.setupEventListeners();
  }

  private setupEventListeners(): void {
    const refreshBtn = document.getElementById('refresh-stats');
    if (refreshBtn) {
      refreshBtn.addEventListener('click', () => {
        this.refreshStats();
      });
    }
  }

  private async fetchStats(): Promise<void> {
    this.updateRefreshButton(true);

    try {
      this.setState({ loading: true, error: null });

      const response = await fetch(`${getApiBase()}/api/stats`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();

      // Validate data structure
      if (!this.isValidStatsData(data)) {
        throw new Error('Invalid data format received from API');
      }

      this.setState({
        loading: false,
        error: null,
        data: data as StatsData
      });

    } catch (error) {
      console.error('Error fetching stats:', error);
      this.setState({
        loading: false,
        error: error instanceof Error ? error.message : 'Failed to load statistics',
        data: null
      });
    } finally {
      this.updateRefreshButton(false);
    }
  }

  private updateRefreshButton(loading: boolean): void {
    const refreshBtn = document.getElementById('refresh-stats') as HTMLButtonElement;
    if (refreshBtn) {
      refreshBtn.disabled = loading;
      if (loading) {
        refreshBtn.classList.add('loading');
        refreshBtn.textContent = 'Refreshing';
      } else {
        refreshBtn.classList.remove('loading');
        refreshBtn.textContent = 'Refresh Statistics';
      }
    }
  }

  private isValidStatsData(data: any): data is StatsData {
    return (
      typeof data === 'object' &&
      data !== null &&
      typeof data.total_links === 'number' &&
      typeof data.total_tags === 'number' &&
      typeof data.total_views === 'number' &&
      Array.isArray(data.top_tags) &&
      Array.isArray(data.most_viewed)
    );
  }

  private setState(newState: Partial<StatsState>): void {
    this.state = { ...this.state, ...newState };
    this.updateDOM();
  }

  private updateDOM(): void {
    if (this.state.loading) {
      this.renderLoadingState();
    } else if (this.state.error) {
      this.renderErrorState(this.state.error);
    } else if (this.state.data) {
      this.renderData(this.state.data);
    }
  }

  private renderLoadingState(): void {
    // Loading state is handled by CSS - just ensure loading classes are present
    const cards = this.container.querySelectorAll('.stats-card');
    const lists = this.container.querySelectorAll('.stats-list');

    cards.forEach(card => {
      card.classList.add('loading');
      const numberEl = card.querySelector('.stats-number');
      if (numberEl) {
        numberEl.classList.add('skeleton');
        numberEl.textContent = '';
      }
    });

    lists.forEach(list => {
      list.classList.add('loading');
      list.innerHTML = `
        <div class="skeleton-item"></div>
        <div class="skeleton-item"></div>
        <div class="skeleton-item"></div>
      `;
    });
  }

  private renderErrorState(error: string): void {
    const cards = this.container.querySelectorAll('.stats-card');
    const lists = this.container.querySelectorAll('.stats-list');

    cards.forEach(card => {
      card.classList.remove('loading');
      const numberEl = card.querySelector('.stats-number');
      if (numberEl) {
        numberEl.classList.remove('skeleton');
        numberEl.textContent = '—';
      }
    });

    lists.forEach(list => {
      list.classList.remove('loading');
      list.innerHTML = `
        <div class="error-message">
          <p>${error}</p>
          <button class="retry-btn" id="retry-stats-${Math.random().toString(36).substr(2, 9)}">Try Again</button>
        </div>
      `;
      const retryBtn = list.querySelector('.retry-btn');
      if (retryBtn) {
        retryBtn.addEventListener('click', () => this.refreshStats());
      }
    });
  }

  private renderData(data: StatsData): void {
    this.renderSummaryCards(data);
    this.renderTopTags(data.top_tags);
    this.renderMostViewed(data.most_viewed);
    this.announceToScreenReader('Statistics loaded successfully');
  }

  private renderSummaryCards(data: StatsData): void {
    const cards = this.container.querySelectorAll('.stats-card');
    const values = [data.total_links, data.total_tags, data.total_views];

    cards.forEach((card, index) => {
      card.classList.remove('loading');
      const numberEl = card.querySelector('.stats-number');
      if (numberEl) {
        numberEl.classList.remove('skeleton');
        numberEl.textContent = this.formatNumber(values[index]);
      }
    });
  }

  private renderTopTags(tags: StatsData['top_tags']): void {
    const topTagsSection = this.container.querySelector('.stats-section h2');
    if (!topTagsSection || topTagsSection.textContent !== 'Top Tags') return;

    const listContainer = topTagsSection.parentElement?.querySelector('.stats-list');
    if (!listContainer) return;

    listContainer.classList.remove('loading');

    if (tags.length === 0) {
      listContainer.innerHTML = '<p class="empty-message">No tags found</p>';
      return;
    }

    const listHTML = tags.map((tag, index) => `
      <div class="ranked-item">
        <span class="rank">${index + 1}</span>
        <span class="name">${this.escapeHtml(tag.name)}</span>
        <span class="count">${this.formatNumber(tag.link_count)} ${tag.link_count === 1 ? 'link' : 'links'}</span>
      </div>
    `).join('');

    listContainer.innerHTML = `<div class="ranked-list">${listHTML}</div>`;
  }

  private renderMostViewed(links: StatsData['most_viewed']): void {
    const mostViewedSection = Array.from(this.container.querySelectorAll('.stats-section h2')).find(h2 => h2.textContent === 'Most Viewed');
    if (!mostViewedSection) return;

    const listContainer = mostViewedSection.parentElement?.querySelector('.stats-list');
    if (!listContainer) return;

    listContainer.classList.remove('loading');

    if (links.length === 0) {
      listContainer.innerHTML = '<p class="empty-message">No links found</p>';
      return;
    }

    const listHTML = links.map((link, index) => `
      <div class="ranked-item">
        <span class="rank">${index + 1}</span>
        <span class="name">
          <a href="${this.escapeHtml(link.url)}" target="_blank" rel="noopener noreferrer">
            ${this.escapeHtml(link.title)}
          </a>
        </span>
        <span class="count">${this.formatNumber(link.view_count)} ${link.view_count === 1 ? 'view' : 'views'}</span>
      </div>
    `).join('');

    listContainer.innerHTML = `<div class="ranked-list">${listHTML}</div>`;
  }

  private formatNumber(num: number): string {
    return new Intl.NumberFormat().format(num);
  }

  private escapeHtml(text: string): string {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  private announceToScreenReader(message: string): void {
    const statusElement = document.getElementById('stats-status');
    if (statusElement) {
      statusElement.textContent = message;
    }
  }

  public async refreshStats(): Promise<void> {
    await this.fetchStats();
  }
}