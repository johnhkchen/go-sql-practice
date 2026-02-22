/**
 * PresenterController - Controls presentation navigation and progress
 * Extracted from components/PresenterController.astro
 */

import {
  stepToProgress,
  progressToStep,
  getNavigationState,
  getNextStep,
  getPreviousStep,
  formatStepDisplay,
  type StepNavigationState
} from '../utils/stepConversion';

interface PresentationData {
  id: string;
  name: string;
  step_count: number;
  step_labels: string[];
  current_step: number;
  progress: number;
}

export class PresenterController {
  private sessionId: string;
  private adminToken: string;
  private presentationData: PresentationData;
  private apiBase: string;
  private currentStep: number;
  private currentProgress: number;
  private isUpdating: boolean = false;
  private throttleTimeout: number | null = null;

  // DOM elements
  private prevButton: HTMLButtonElement | null = null;
  private nextButton: HTMLButtonElement | null = null;
  private jumpButtons: NodeListOf<HTMLButtonElement> | null = null;
  private progressSlider: HTMLInputElement | null = null;
  private progressDisplay: HTMLSpanElement | null = null;
  private stepDisplay: HTMLElement | null = null;
  private copyButton: HTMLButtonElement | null = null;
  private viewerUrlInput: HTMLInputElement | null = null;
  private stopButton: HTMLButtonElement | null = null;
  private statusArea: HTMLElement | null = null;
  private loadingOverlay: HTMLElement | null = null;

  constructor(sessionId: string, adminToken: string, presentationData: PresentationData, apiBase: string) {
    this.sessionId = sessionId;
    this.adminToken = adminToken;
    this.presentationData = presentationData;
    this.apiBase = apiBase;
    this.currentStep = presentationData.current_step;
    this.currentProgress = presentationData.progress;

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
    // Get DOM elements
    this.prevButton = document.querySelector('.prev-button');
    this.nextButton = document.querySelector('.next-button');
    this.jumpButtons = document.querySelectorAll('.jump-button');
    this.progressSlider = document.getElementById('progress-slider') as HTMLInputElement;
    this.progressDisplay = document.getElementById('progress-display') as HTMLSpanElement;
    this.stepDisplay = document.getElementById('step-display');
    this.copyButton = document.getElementById('copy-button') as HTMLButtonElement;
    this.viewerUrlInput = document.getElementById('viewer-url') as HTMLInputElement;
    this.stopButton = document.querySelector('.stop-button');
    this.statusArea = document.getElementById('status-area');
    this.loadingOverlay = document.getElementById('loading-overlay');

    if (!this.validateDOMElements()) {
      console.warn('PresenterController: Some DOM elements not found');
      return;
    }

    this.setupEventListeners();
  }

  private validateDOMElements(): boolean {
    return !!(
      this.prevButton && this.nextButton && this.jumpButtons &&
      this.progressSlider && this.progressDisplay && this.stepDisplay &&
      this.copyButton && this.viewerUrlInput && this.stopButton &&
      this.statusArea && this.loadingOverlay
    );
  }

  private setupEventListeners(): void {
    // Previous/Next buttons
    this.prevButton?.addEventListener('click', () => this.handlePrevious());
    this.nextButton?.addEventListener('click', () => this.handleNext());

    // Jump buttons
    this.jumpButtons?.forEach((button) => {
      button.addEventListener('click', () => {
        const stepIndex = parseInt(button.dataset.stepIndex || '0', 10);
        this.handleStepJump(stepIndex);
      });
    });

    // Progress slider
    this.progressSlider?.addEventListener('input', (e) => {
      const target = e.target as HTMLInputElement;
      this.handleSliderChange(parseFloat(target.value));
    });

    // Enhanced keyboard support for slider
    this.progressSlider?.addEventListener('keydown', (e) => {
      this.handleSliderKeyboard(e);
    });

    // Copy button
    this.copyButton?.addEventListener('click', () => this.copyToClipboard());

    // Stop presenting button
    this.stopButton?.addEventListener('click', () => this.handleStopPresenting());

    // URL input focus
    this.viewerUrlInput?.addEventListener('focus', () => {
      this.viewerUrlInput?.select();
    });
  }

  private async handlePrevious(): Promise<void> {
    if (!this.canGoPrevious()) return;

    const prevStep = getPreviousStep(this.currentStep, this.presentationData.step_count);
    await this.updateStep(prevStep);
  }

  private async handleNext(): Promise<void> {
    if (!this.canGoNext()) return;

    const nextStep = getNextStep(this.currentStep, this.presentationData.step_count);
    await this.updateStep(nextStep);
  }

  private async handleStepJump(stepIndex: number): Promise<void> {
    if (stepIndex < 0 || stepIndex >= this.presentationData.step_count) {
      console.warn('Invalid step index:', stepIndex);
      return;
    }

    await this.updateStep(stepIndex);
  }

  private handleSliderChange(progress: number): void {
    this.updateProgressDisplay(progress);
    this.throttledProgressUpdate(progress);
  }

  private handleSliderKeyboard(e: KeyboardEvent): void {
    let increment = 0;
    const stepSize = 1 / (this.presentationData.step_count - 1);

    switch (e.key) {
      case 'ArrowRight':
      case 'ArrowUp':
        increment = e.shiftKey ? stepSize : e.ctrlKey ? 0.1 : 0.001;
        break;
      case 'ArrowLeft':
      case 'ArrowDown':
        increment = e.shiftKey ? -stepSize : e.ctrlKey ? -0.1 : -0.001;
        break;
      case 'Home':
        e.preventDefault();
        this.handleSliderChange(0);
        return;
      case 'End':
        e.preventDefault();
        this.handleSliderChange(1);
        return;
      case 'PageUp':
        increment = stepSize;
        break;
      case 'PageDown':
        increment = -stepSize;
        break;
      default:
        return;
    }

    if (increment !== 0) {
      e.preventDefault();
      const newValue = Math.max(0, Math.min(1, this.currentProgress + increment));
      if (this.progressSlider) {
        this.progressSlider.value = newValue.toString();
      }
      this.handleSliderChange(newValue);
    }
  }

  private async updateStep(stepIndex: number): Promise<void> {
    const progress = stepToProgress(stepIndex, this.presentationData.step_count);
    await this.sendProgressUpdate(progress);
  }

  private updateProgressDisplay(progress: number): void {
    this.currentProgress = progress;

    if (this.progressDisplay) {
      this.progressDisplay.textContent = progress.toFixed(3);
    }

    // Update slider ARIA attributes
    if (this.progressSlider) {
      this.progressSlider.setAttribute('aria-valuenow', progress.toString());
      this.progressSlider.setAttribute('aria-valuetext', `${progress.toFixed(3)} out of 1`);
    }

    // Calculate current step from progress
    const newStep = progressToStep(progress, this.presentationData.step_count);
    this.updateStepDisplay(newStep);
  }

  private updateStepDisplay(stepIndex: number): void {
    this.currentStep = stepIndex;

    // Update step indicator
    if (this.stepDisplay) {
      const displayText = formatStepDisplay(
        stepIndex,
        this.presentationData.step_count,
        this.presentationData.step_labels
      );
      this.stepDisplay.textContent = displayText;
    }

    // Update navigation button states
    const navState = getNavigationState(stepIndex, this.presentationData.step_count);

    if (this.prevButton) {
      this.prevButton.disabled = !navState.canGoPrevious;
    }

    if (this.nextButton) {
      this.nextButton.disabled = !navState.canGoNext;
    }

    // Update jump button states
    this.jumpButtons?.forEach((button, index) => {
      const isActive = index === stepIndex;
      button.classList.toggle('active', isActive);
      button.setAttribute('aria-pressed', isActive.toString());
    });
  }

  private throttledProgressUpdate(progress: number): void {
    if (this.throttleTimeout) {
      clearTimeout(this.throttleTimeout);
    }

    this.throttleTimeout = setTimeout(() => {
      this.sendProgressUpdate(progress);
    }, 33); // ~30 updates/sec max
  }

  private async sendProgressUpdate(progress: number): Promise<void> {
    if (this.isUpdating) {
      return;
    }

    this.isUpdating = true;
    this.showLoading(true);

    try {
      const response = await fetch(`${this.apiBase}/api/sync/${this.sessionId}/progress`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify({
          progress: progress,
          token: this.adminToken
        })
      });

      if (response.ok) {
        const result = await response.json();

        // Update local state with server response
        if (result.progress !== undefined) {
          this.updateProgressDisplay(result.progress);
          if (this.progressSlider) {
            this.progressSlider.value = result.progress.toString();
          }
        }

        this.showStatus('success', 'Progress updated successfully');

      } else {
        let errorMessage = 'Failed to update progress';

        if (response.status === 403) {
          errorMessage = 'Invalid admin token';
        } else if (response.status === 404) {
          errorMessage = 'Session not found';
        } else if (response.status === 400) {
          const errorData = await response.json().catch(() => ({}));
          errorMessage = errorData.message || 'Invalid request';
        }

        this.showStatus('error', errorMessage);
        console.error('API Error:', errorMessage);
      }

    } catch (error) {
      console.error('Network error:', error);
      const errorMessage = error instanceof Error && error.message.includes('fetch')
        ? 'Unable to connect to server'
        : 'Network error';
      this.showStatus('error', errorMessage);

    } finally {
      this.isUpdating = false;
      this.showLoading(false);
    }
  }

  private async copyToClipboard(): Promise<void> {
    if (!this.viewerUrlInput) return;

    const viewerUrl = this.viewerUrlInput.value;

    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(viewerUrl);
        this.showCopyFeedback('success');
        this.announceToScreenReader('Viewer URL copied to clipboard');
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
        this.announceToScreenReader('Viewer URL copied to clipboard');
      } else {
        throw new Error('execCommand copy failed');
      }
    } finally {
      document.body.removeChild(textarea);
    }
  }

  private showCopyFeedback(type: 'success' | 'error', message?: string): void {
    if (!this.copyButton) return;

    const originalText = this.copyButton.textContent;

    if (type === 'success') {
      this.copyButton.textContent = '✓ Copied!';
      this.copyButton.classList.add('copy-success');

      setTimeout(() => {
        if (this.copyButton) {
          this.copyButton.textContent = originalText;
          this.copyButton.classList.remove('copy-success');
        }
      }, 2000);

    } else {
      this.copyButton.textContent = '⚠ Error';
      this.copyButton.classList.add('copy-error');

      if (message) {
        this.showStatus('error', message);
      }

      setTimeout(() => {
        if (this.copyButton) {
          this.copyButton.textContent = originalText;
          this.copyButton.classList.remove('copy-error');
        }
      }, 3000);
    }
  }

  private async handleStopPresenting(): Promise<void> {
    const confirmed = confirm('Are you sure you want to stop presenting? This will end the live session.');

    if (!confirmed) return;

    this.showLoading(true);

    try {
      const response = await fetch(`${this.apiBase}/api/presentations/${this.presentationData.id}/stop`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify({
          token: this.adminToken
        })
      });

      if (response.ok) {
        this.showStatus('success', 'Presentation stopped. Redirecting to dashboard...');

        // Redirect to dashboard after brief delay
        setTimeout(() => {
          window.location.href = '/present';
        }, 1500);

      } else {
        const errorData = await response.json().catch(() => ({}));
        const errorMessage = errorData.message || `Failed to stop presentation (HTTP ${response.status})`;
        this.showStatus('error', errorMessage);
        console.error('Stop presentation error:', errorMessage);
      }

    } catch (error) {
      console.error('Network error stopping presentation:', error);
      this.showStatus('error', 'Unable to stop presentation. Please try again.');
    } finally {
      this.showLoading(false);
    }
  }

  private showStatus(type: 'success' | 'error' | 'info', message: string): void {
    if (!this.statusArea) return;

    // Clear existing status
    this.statusArea.innerHTML = '';

    const statusElement = document.createElement('div');
    statusElement.className = `status-message ${type}`;

    const icon = type === 'success' ? '✓' : type === 'error' ? '⚠' : 'ℹ';
    statusElement.innerHTML = `<span aria-hidden="true">${icon}</span> ${message}`;

    this.statusArea.appendChild(statusElement);

    // Auto-remove after delay
    setTimeout(() => {
      if (statusElement.parentNode) {
        statusElement.remove();
      }
    }, type === 'success' ? 3000 : 5000);

    // Announce to screen readers
    this.announceToScreenReader(message);
  }

  private showLoading(show: boolean): void {
    if (!this.loadingOverlay) return;

    this.loadingOverlay.setAttribute('aria-hidden', show ? 'false' : 'true');
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
      if (document.body.contains(announcement)) {
        document.body.removeChild(announcement);
      }
    }, 1000);
  }

  private canGoPrevious(): boolean {
    return this.currentStep > 0;
  }

  private canGoNext(): boolean {
    return this.currentStep < this.presentationData.step_count - 1;
  }
}