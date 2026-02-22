/**
 * SyncController - Controls sync session progress
 * Extracted from sync/[id]/control.astro
 */

export class SyncController {
  private sessionId: string;
  private adminToken: string;
  private currentProgress: number;
  private apiBase: string;

  // DOM elements
  private progressSlider: HTMLInputElement | null = null;
  private progressValue: HTMLElement | null = null;
  private copyButton: HTMLButtonElement | null = null;
  private viewerUrlInput: HTMLInputElement | null = null;

  // State management
  private isUpdating: boolean = false;
  private throttleTimeout: number | null = null;
  private lastAnnouncedValue?: number;

  constructor(sessionId: string, adminToken: string, initialProgress: number, apiBase: string) {
    this.sessionId = sessionId;
    this.adminToken = adminToken;
    this.currentProgress = initialProgress;
    this.apiBase = apiBase;

    this.init();
  }

  private init(): void {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => this.setupDOM());
    } else {
      this.setupDOM();
    }
  }

  private setupDOM(): void {
    this.progressSlider = document.getElementById('progress-slider') as HTMLInputElement;
    this.progressValue = document.getElementById('progress-value');
    this.copyButton = document.getElementById('copy-button') as HTMLButtonElement;
    this.viewerUrlInput = document.getElementById('viewer-url') as HTMLInputElement;

    if (!this.progressSlider || !this.progressValue || !this.copyButton || !this.viewerUrlInput) {
      console.warn('SyncController: Some DOM elements not found. Page may not be fully loaded.');
      return;
    }

    this.setupEventListeners();
  }

  private setupEventListeners(): void {
    this.progressSlider?.addEventListener('input', (e) => {
      this.updateProgress(parseFloat((e.target as HTMLInputElement).value));
    });

    this.progressSlider?.addEventListener('keydown', (e) => {
      this.handleSliderKeyboard(e);
    });

    this.copyButton?.addEventListener('click', () => {
      this.copyToClipboard();
    });

    this.copyButton?.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        this.copyToClipboard();
      }
    });

    this.viewerUrlInput?.addEventListener('focus', () => {
      this.viewerUrlInput?.select();
    });
  }

  private handleSliderKeyboard(e: KeyboardEvent): void {
    let increment = 0;

    switch (e.key) {
      case 'ArrowRight':
      case 'ArrowUp':
        increment = e.shiftKey ? 0.01 : e.ctrlKey ? 0.1 : 0.001;
        break;
      case 'ArrowLeft':
      case 'ArrowDown':
        increment = e.shiftKey ? -0.01 : e.ctrlKey ? -0.1 : -0.001;
        break;
      case 'Home':
        e.preventDefault();
        this.updateProgress(0);
        return;
      case 'End':
        e.preventDefault();
        this.updateProgress(1);
        return;
      case 'PageUp':
        increment = 0.1;
        break;
      case 'PageDown':
        increment = -0.1;
        break;
      default:
        return;
    }

    if (increment !== 0) {
      e.preventDefault();
      const newValue = parseFloat(this.progressSlider!.value) + increment;
      this.progressSlider!.value = Math.max(0, Math.min(1, newValue)).toString();
      this.updateProgress(parseFloat(this.progressSlider!.value));
    }
  }

  private updateProgress(value: number): void {
    const validatedValue = this.validateProgress(value);
    this.currentProgress = validatedValue;

    this.updateProgressDisplay(validatedValue);
    this.throttledUpdate(validatedValue);
  }

  private updateProgressDisplay(value: number): void {
    const span = this.progressValue?.querySelector('span');
    if (span) {
      span.textContent = value.toFixed(3);
    }

    this.progressSlider?.setAttribute('aria-valuenow', value.toString());
    this.progressSlider?.setAttribute('aria-valuetext', `${value.toFixed(3)} out of 1`);

    if (this.progressValue) {
      this.progressValue.style.color = 'var(--color-primary)';
      this.progressValue.style.fontWeight = '700';

      setTimeout(() => {
        if (this.progressValue) {
          this.progressValue.style.color = '';
          this.progressValue.style.fontWeight = '';
        }
      }, 150);
    }

    if (Math.abs(value - (this.lastAnnouncedValue || 0)) >= 0.1 || !this.lastAnnouncedValue) {
      this.announceProgressChange(value);
      this.lastAnnouncedValue = value;
    }
  }

  private announceProgressChange(value: number): void {
    const announcement = document.createElement('div');
    announcement.setAttribute('aria-live', 'polite');
    announcement.setAttribute('aria-atomic', 'true');
    announcement.className = 'sr-only';

    const percentage = Math.round(value * 100);
    announcement.textContent = `Progress updated to ${percentage}%`;

    document.body.appendChild(announcement);

    setTimeout(() => {
      if (document.body.contains(announcement)) {
        document.body.removeChild(announcement);
      }
    }, 1000);
  }

  private validateProgress(value: number): number {
    if (typeof value !== 'number' || isNaN(value)) {
      return this.currentProgress;
    }

    return Math.max(0, Math.min(1, value));
  }

  private throttledUpdate(value: number): void {
    if (this.throttleTimeout) {
      clearTimeout(this.throttleTimeout);
    }

    this.throttleTimeout = setTimeout(() => {
      this.sendProgressUpdate(value);
    }, 33);
  }

  private async sendProgressUpdate(value: number): Promise<void> {
    if (this.isUpdating) {
      return;
    }

    this.isUpdating = true;

    try {
      this.addUpdateFeedback('updating');

      const response = await fetch(`${this.apiBase}/api/sync/${this.sessionId}/progress`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify({
          progress: value,
          token: this.adminToken
        })
      });

      if (response.ok) {
        const result = await response.json();
        this.addUpdateFeedback('success');

        if (result.progress !== undefined) {
          this.currentProgress = result.progress;
          this.updateProgressDisplay(result.progress);
        }
      } else {
        let errorMessage = 'Failed to update progress';

        if (response.status === 403) {
          errorMessage = 'Invalid admin token';
          this.handleError('invalid_token');
        } else if (response.status === 404) {
          errorMessage = 'Session not found';
          this.handleError('session_not_found');
        } else if (response.status === 400) {
          const errorData = await response.json().catch(() => ({}));
          errorMessage = errorData.message || 'Invalid request';
          this.handleError('validation_error', errorMessage);
        } else {
          errorMessage = `Server error (${response.status})`;
          this.handleError('server_error');
        }

        console.error('API Error:', errorMessage);
        this.addUpdateFeedback('error', errorMessage);
      }
    } catch (error) {
      console.error('Network error:', error);

      let errorMessage = 'Network error';
      if (error instanceof Error && error.message.includes('fetch')) {
        errorMessage = 'Unable to connect to server';
      }

      this.handleError('network_error', errorMessage);
      this.addUpdateFeedback('error', errorMessage);
    } finally {
      this.isUpdating = false;
    }
  }

  private addUpdateFeedback(type: 'updating' | 'success' | 'error', message = ''): void {
    this.removeUpdateFeedback();

    const feedbackElement = document.createElement('div');
    feedbackElement.className = `update-feedback update-feedback-${type}`;
    feedbackElement.id = 'update-feedback';

    let content = '';
    switch (type) {
      case 'updating':
        content = '<span>Updating...</span>';
        break;
      case 'success':
        content = '<span>✓ Updated</span>';
        break;
      case 'error':
        content = `<span>⚠ ${message || 'Update failed'}</span>`;
        break;
    }

    feedbackElement.innerHTML = content;

    const progressControls = document.querySelector('.progress-controls');
    if (progressControls) {
      progressControls.appendChild(feedbackElement);
    }

    if (type === 'success' || type === 'error') {
      setTimeout(() => {
        this.removeUpdateFeedback();
      }, type === 'success' ? 2000 : 5000);
    }
  }

  private removeUpdateFeedback(): void {
    const existing = document.getElementById('update-feedback');
    if (existing) {
      existing.remove();
    }
  }

  private async copyToClipboard(): Promise<void> {
    const viewerUrl = this.viewerUrlInput!.value;

    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(viewerUrl);
        this.showCopyFeedback('success');
        this.announceCopySuccess();
      } else {
        this.copyWithLegacyMethod(viewerUrl);
      }
    } catch (error) {
      console.error('Copy failed:', error);

      try {
        this.copyWithLegacyMethod(viewerUrl);
      } catch (fallbackError) {
        console.error('Fallback copy also failed:', fallbackError);
        this.showCopyFeedback('error', 'Copy failed. Please copy the URL manually.');
      }
    }
  }

  private copyWithLegacyMethod(text: string): void {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    textarea.style.left = '-9999px';
    textarea.setAttribute('aria-hidden', 'true');

    document.body.appendChild(textarea);

    try {
      textarea.select();
      textarea.setSelectionRange(0, 99999);

      const successful = document.execCommand('copy');

      if (successful) {
        this.showCopyFeedback('success');
        this.announceCopySuccess();
      } else {
        throw new Error('execCommand copy failed');
      }
    } finally {
      document.body.removeChild(textarea);
    }
  }

  private showCopyFeedback(type: 'success' | 'error', message = ''): void {
    const existing = document.querySelector('.copy-feedback');
    if (existing) {
      existing.remove();
    }

    const originalText = this.copyButton!.textContent;

    if (type === 'success') {
      this.copyButton!.textContent = '✓ Copied!';
      this.copyButton!.classList.add('copy-success');

      setTimeout(() => {
        if (this.copyButton) {
          this.copyButton.textContent = originalText;
          this.copyButton.classList.remove('copy-success');
        }
      }, 2000);

    } else if (type === 'error') {
      this.copyButton!.textContent = '⚠ Error';
      this.copyButton!.classList.add('copy-error');

      const feedbackElement = document.createElement('div');
      feedbackElement.className = 'copy-feedback copy-feedback-error';
      feedbackElement.textContent = message || 'Copy failed';

      const urlContainer = document.querySelector('.url-container');
      if (urlContainer) {
        urlContainer.appendChild(feedbackElement);

        setTimeout(() => {
          feedbackElement.remove();
        }, 5000);
      }

      setTimeout(() => {
        if (this.copyButton) {
          this.copyButton.textContent = originalText;
          this.copyButton.classList.remove('copy-error');
        }
      }, 3000);
    }
  }

  private announceCopySuccess(): void {
    const announcement = document.createElement('div');
    announcement.setAttribute('aria-live', 'polite');
    announcement.setAttribute('aria-atomic', 'true');
    announcement.style.position = 'absolute';
    announcement.style.left = '-10000px';
    announcement.style.width = '1px';
    announcement.style.height = '1px';
    announcement.style.overflow = 'hidden';

    announcement.textContent = 'Viewer URL copied to clipboard';

    document.body.appendChild(announcement);

    setTimeout(() => {
      document.body.removeChild(announcement);
    }, 1000);
  }

  private handleError(errorType: string, message = ''): void {
    console.error('SyncController error:', errorType, message);

    switch (errorType) {
      case 'invalid_token':
        this.showErrorRecovery('invalid_token', 'Invalid admin token. Please check your URL and refresh the page.');
        break;

      case 'session_not_found':
        this.showErrorRecovery('session_not_found', 'Session not found. It may have been deleted.');
        break;

      case 'network_error':
        this.showNetworkError();
        break;

      case 'server_error':
        this.showErrorRecovery('server_error', 'Server error. Please try again later.');
        break;

      case 'validation_error':
        this.showErrorRecovery('validation_error', message || 'Invalid request.');
        break;

      default:
        this.showErrorRecovery('unknown_error', message || 'An unexpected error occurred.');
    }
  }

  private showErrorRecovery(errorType: string, message: string): void {
    const errorDialog = this.createErrorDialog(errorType, message);
    document.body.appendChild(errorDialog);

    const firstButton = errorDialog.querySelector('button');
    if (firstButton) {
      firstButton.focus();
    }
  }

  private showNetworkError(): void {
    const errorDialog = this.createNetworkErrorDialog();
    document.body.appendChild(errorDialog);

    const retryButton = errorDialog.querySelector('.retry-button');
    if (retryButton) {
      retryButton.focus();
    }
  }

  private createErrorDialog(errorType: string, message: string): HTMLElement {
    const dialog = document.createElement('div');
    dialog.className = 'error-dialog-overlay';
    dialog.setAttribute('role', 'dialog');
    dialog.setAttribute('aria-modal', 'true');
    dialog.setAttribute('aria-labelledby', 'error-title');

    dialog.innerHTML = `
      <div class="error-dialog">
        <div class="error-dialog-header">
          <h3 id="error-title">⚠ Error</h3>
        </div>
        <div class="error-dialog-body">
          <p>${message}</p>
        </div>
        <div class="error-dialog-actions">
          <button type="button" class="error-dialog-button primary" onclick="location.reload()">
            Reload Page
          </button>
          <button type="button" class="error-dialog-button secondary" onclick="this.closest('.error-dialog-overlay').remove()">
            Dismiss
          </button>
        </div>
      </div>
    `;

    dialog.addEventListener('keydown', (e) => {
      if (e.key === 'Escape') {
        dialog.remove();
      }
    });

    dialog.addEventListener('click', (e) => {
      if (e.target === dialog) {
        dialog.remove();
      }
    });

    return dialog;
  }

  private createNetworkErrorDialog(): HTMLElement {
    const dialog = document.createElement('div');
    dialog.className = 'error-dialog-overlay';
    dialog.setAttribute('role', 'dialog');
    dialog.setAttribute('aria-modal', 'true');
    dialog.setAttribute('aria-labelledby', 'error-title');

    dialog.innerHTML = `
      <div class="error-dialog">
        <div class="error-dialog-header">
          <h3 id="error-title">🌐 Connection Problem</h3>
        </div>
        <div class="error-dialog-body">
          <p>Unable to connect to the server. Please check your connection and try again.</p>
          <div class="connection-status" id="connection-status">
            <span class="status-indicator offline"></span>
            <span class="status-text">Offline</span>
          </div>
        </div>
        <div class="error-dialog-actions">
          <button type="button" class="error-dialog-button primary retry-button" onclick="this.retryConnection(this)">
            Retry Connection
          </button>
          <button type="button" class="error-dialog-button secondary" onclick="this.closest('.error-dialog-overlay').remove()">
            Continue Offline
          </button>
        </div>
      </div>
    `;

    const retryButton = dialog.querySelector('.retry-button');
    retryButton?.addEventListener('click', () => {
      this.testConnection(dialog);
    });

    this.setupConnectionMonitoring(dialog);

    return dialog;
  }

  private async testConnection(dialog: HTMLElement): Promise<void> {
    const statusIndicator = dialog.querySelector('.status-indicator');
    const statusText = dialog.querySelector('.status-text');
    const retryButton = dialog.querySelector('.retry-button') as HTMLButtonElement;

    if (statusIndicator && statusText && retryButton) {
      statusIndicator.className = 'status-indicator testing';
      statusText.textContent = 'Testing...';
      retryButton.disabled = true;

      try {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 5000);

        const response = await fetch(`${this.apiBase}/api/health`, {
          method: 'GET',
          signal: controller.signal
        });

        clearTimeout(timeoutId);

        if (response.ok) {
          statusIndicator.className = 'status-indicator online';
          statusText.textContent = 'Connected';
          retryButton.textContent = 'Close';
          retryButton.disabled = false;

          setTimeout(() => {
            dialog.remove();
          }, 1000);
        } else {
          throw new Error('Server responded with error');
        }
      } catch (error) {
        statusIndicator.className = 'status-indicator offline';
        statusText.textContent = 'Still offline';
        retryButton.disabled = false;
      }
    }
  }

  private setupConnectionMonitoring(dialog: HTMLElement): void {
    const statusIndicator = dialog.querySelector('.status-indicator');
    const statusText = dialog.querySelector('.status-text');

    const updateStatus = () => {
      if (navigator.onLine) {
        if (statusIndicator) statusIndicator.className = 'status-indicator online';
        if (statusText) statusText.textContent = 'Back online';
      } else {
        if (statusIndicator) statusIndicator.className = 'status-indicator offline';
        if (statusText) statusText.textContent = 'Offline';
      }
    };

    window.addEventListener('online', updateStatus);
    window.addEventListener('offline', updateStatus);

    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        mutation.removedNodes.forEach((node) => {
          if (node === dialog) {
            window.removeEventListener('online', updateStatus);
            window.removeEventListener('offline', updateStatus);
            observer.disconnect();
          }
        });
      });
    });

    observer.observe(document.body, { childList: true });
    updateStatus();
  }

  public checkJavaScriptSupport(): void {
    document.documentElement.classList.add('js-enabled');
  }
}