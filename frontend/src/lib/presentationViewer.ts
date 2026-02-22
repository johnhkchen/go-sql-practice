/**
 * PresentationAutoViewer - Handles live presentation viewing with SSE updates
 * Extracted from watch/[id].astro
 */

// Step conversion utilities (inlined for client script)
function progressToStep(progress: number, stepCount: number): number {
  if (stepCount <= 1) return 0;
  const clampedProgress = Math.max(0.0, Math.min(progress, 1.0));
  return Math.round(clampedProgress * (stepCount - 1));
}

function formatStepDisplay(stepIndex: number, stepCount: number, stepLabels?: string[]): string {
  const stepNumber = stepIndex + 1;
  let display = `Step ${stepNumber} of ${stepCount}`;
  if (stepLabels && stepLabels[stepIndex]) {
    display += ` — ${stepLabels[stepIndex]}`;
  }
  return display;
}

export class PresentationAutoViewer {
  private presentationId: string;
  private sessionId: string | null;
  private currentProgress: number;
  private stepCount: number;
  private stepLabels: string[];
  private apiBase: string;

  // State management
  private currentState: 'waiting' | 'starting' | 'live' | 'ended';
  private connectionStatus: 'connected' | 'connecting' | 'disconnected' = 'disconnected';
  private eventSource: EventSource | null = null;
  private reconnectAttempts: number = 0;
  private maxReconnectAttempts: number = 10;
  private currentStep: number = 0;

  // Debouncing and error handling
  private lastStateChangeTime: number = 0;
  private stateChangeDebounceMs: number = 500;

  // DOM elements
  private container: HTMLElement | null = null;
  private statusIndicator: HTMLElement | null = null;
  private statusText: HTMLElement | null = null;
  private lastUpdate: HTMLElement | null = null;
  private stepHeading: HTMLElement | null = null;
  private stepIndicators: HTMLElement | null = null;
  private progressBar: HTMLProgressElement | null = null;
  private progressLabel: HTMLElement | null = null;
  private transitionMessage: HTMLElement | null = null;
  private waitAgainBtn: HTMLElement | null = null;

  constructor(
    presentationId: string,
    sessionId: string | null,
    initialProgress: number,
    stepCount: number,
    stepLabels: string[],
    apiBase: string,
    initialState: 'waiting' | 'starting' | 'live' | 'ended'
  ) {
    this.presentationId = presentationId;
    this.sessionId = sessionId;
    this.currentProgress = initialProgress;
    this.stepCount = stepCount;
    this.stepLabels = stepLabels;
    this.apiBase = apiBase;
    this.currentState = initialState;

    // DOM elements
    this.container = document.getElementById('presentation-container');
    this.statusIndicator = document.getElementById('status-indicator');
    this.statusText = document.getElementById('status-text');
    this.lastUpdate = document.getElementById('last-update');
    this.stepHeading = document.getElementById('step-heading');
    this.stepIndicators = document.getElementById('step-indicators');
    this.progressBar = document.getElementById('progress-bar') as HTMLProgressElement;
    this.progressLabel = document.getElementById('progress-label');
    this.transitionMessage = document.getElementById('transition-message');
    this.waitAgainBtn = document.getElementById('wait-again-btn');

    // Initialize
    this.currentStep = progressToStep(this.currentProgress, this.stepCount);
    this.updateConnectionStatus('connecting');
    this.updateLastUpdate('Initializing...');

    // Update connection status based on initial state
    setTimeout(() => {
      if (this.currentState === 'waiting') {
        this.updateConnectionStatus('connected');
        this.updateLastUpdate('Waiting for presenter to start');
      }
    }, 1000);

    this.setupEventListeners();
    this.connect();
  }

  private setupEventListeners(): void {
    if (this.waitAgainBtn) {
      this.waitAgainBtn.addEventListener('click', () => {
        this.transitionToState('waiting');
      });
    }
  }

  private transitionToState(newState: 'waiting' | 'starting' | 'live' | 'ended', transitionMessage?: string | null): void {
    if (this.currentState === newState) return;

    const now = Date.now();
    if (now - this.lastStateChangeTime < this.stateChangeDebounceMs) {
      return;
    }
    this.lastStateChangeTime = now;

    if (this.container) {
      this.container.classList.add('transitioning');
    }

    if (transitionMessage && this.transitionMessage) {
      this.transitionMessage.textContent = transitionMessage;
    }

    setTimeout(() => {
      this.currentState = newState;
      if (this.container) {
        this.container.setAttribute('data-state', newState);
        this.container.classList.remove('transitioning');
      }
    }, 300);
  }

  private updateLastUpdate(message?: string): void {
    if (this.lastUpdate) {
      const now = new Date().toLocaleTimeString();
      this.lastUpdate.textContent = message || `Updated ${now}`;
    }
  }

  private updateConnectionStatus(status: 'connected' | 'connecting' | 'disconnected'): void {
    if (!this.statusIndicator || !this.statusText) return;

    this.connectionStatus = status;
    this.statusIndicator.className = `status-indicator ${status}`;

    switch (status) {
      case 'connected':
        if (this.currentState === 'waiting') {
          this.statusText.textContent = 'Waiting for presenter';
        } else {
          this.statusText.textContent = 'Connected';
        }
        break;
      case 'connecting':
        this.statusText.textContent = 'Connecting...';
        break;
      case 'disconnected':
        if (this.currentState === 'ended') {
          this.statusText.textContent = 'Presentation ended';
        } else {
          this.statusText.textContent = 'Disconnected';
        }
        break;
    }
  }

  private connect(): void {
    if (this.eventSource) {
      this.eventSource.close();
    }

    try {
      this.eventSource = new EventSource(`${this.apiBase}/api/realtime`);

      this.eventSource.onopen = () => this.handleConnect();
      this.eventSource.onmessage = (event) => this.handleMessage(event);
      this.eventSource.onerror = () => this.handleError();
    } catch (error) {
      console.error('Failed to create SSE connection:', error);
      this.handleError();
    }
  }

  private handleConnect(): void {
    this.updateConnectionStatus('connected');
    this.updateLastUpdate('Connected');
    this.reconnectAttempts = 0;

    this.syncWithServerState();
  }

  private async syncWithServerState(): Promise<void> {
    try {
      const response = await fetch(`${this.apiBase}/api/presentations/${this.presentationId}/status`);
      if (response.ok) {
        const serverPresentation = await response.json();

        const serverHasActiveSession = serverPresentation.active_session !== null;
        const clientInWaitingState = this.currentState === 'waiting';
        const clientInLiveState = this.currentState === 'live' || this.currentState === 'starting';

        if (serverHasActiveSession && clientInWaitingState) {
          this.sessionId = serverPresentation.active_session;
          this.currentProgress = serverPresentation.progress || 0;

          this.transitionToState('starting', 'Joining live presentation...');

          setTimeout(() => {
            this.transitionToState('live');
            this.updateProgress(serverPresentation.progress || 0);
            this.updateConnectionStatus('connected');
            this.updateLastUpdate('Joined live presentation');
          }, 1500);
        }
        else if (!serverHasActiveSession && clientInLiveState) {
          this.transitionToState('ended');
          this.updateConnectionStatus('disconnected');
          this.updateLastUpdate('Presentation has ended');
          this.sessionId = null;
        }
        else if (serverHasActiveSession && clientInLiveState &&
                 this.sessionId === serverPresentation.active_session &&
                 Math.abs((serverPresentation.progress || 0) - this.currentProgress) > 0.01) {
          this.updateProgress(serverPresentation.progress || 0);
        }
      }
    } catch (error) {
      console.warn('Failed to sync with server state:', error);
    }
  }

  private handleMessage(event: MessageEvent): void {
    try {
      const data = JSON.parse(event.data);

      if (data.collection === 'presentations' &&
          data.action === 'update' &&
          data.record &&
          data.record.id === this.presentationId) {
        this.handlePresentationStateChange(data.record);
      }

      if (data.collection === 'sync_sessions' &&
          data.action === 'update' &&
          data.record &&
          data.record.id === this.sessionId) {
        this.updateProgress(data.record.progress);
        this.updateLastUpdate();
      }
    } catch (error) {
      console.warn('Invalid SSE message:', error);
    }
  }

  private handlePresentationStateChange(presentationRecord: any): void {
    const newActiveSession = presentationRecord.active_session;

    // Waiting → Live transition
    if (this.currentState === 'waiting' && newActiveSession !== null) {
      this.sessionId = newActiveSession;
      this.transitionToState('starting', 'Presentation starting...');

      setTimeout(() => {
        this.transitionToState('live');
        this.updateConnectionStatus('connected');
        this.updateLastUpdate('Live presentation started');
      }, 1500);
    }
    // Live → Ended transition
    else if ((this.currentState === 'live' || this.currentState === 'starting') && newActiveSession === null) {
      this.transitionToState('ended');
      this.updateConnectionStatus('disconnected');
      this.updateLastUpdate('Presentation ended');
      this.sessionId = null;
    }
    // Ended → Live transition (new session started)
    else if (this.currentState === 'ended' && newActiveSession !== null) {
      this.sessionId = newActiveSession;
      this.transitionToState('starting', 'New presentation starting...');

      setTimeout(() => {
        this.transitionToState('live');
        this.updateConnectionStatus('connected');
        this.updateLastUpdate('Live presentation started');
      }, 1500);
    }
  }

  private handleError(): void {
    this.updateConnectionStatus('disconnected');
    this.updateLastUpdate('Connection lost');

    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      const delay = Math.pow(2, this.reconnectAttempts) * 1000;
      this.reconnectAttempts++;

      setTimeout(() => {
        this.updateConnectionStatus('connecting');
        this.updateLastUpdate('Reconnecting...');
        this.connect();
      }, delay);
    }
  }

  private updateProgress(progress: number): void {
    this.currentProgress = progress;

    const newStep = progressToStep(progress, this.stepCount);

    if (this.progressBar) {
      this.progressBar.value = progress;
    }

    if (this.progressLabel) {
      this.progressLabel.textContent = `${Math.round(progress * 100)}%`;
    }

    if (newStep !== this.currentStep) {
      this.updateStepDisplay(newStep);
      this.announceStepChange(newStep);
    }
  }

  private updateStepDisplay(newStep: number): void {
    this.currentStep = newStep;

    if (this.stepHeading) {
      const stepDisplay = formatStepDisplay(newStep, this.stepCount, this.stepLabels);
      this.stepHeading.textContent = stepDisplay;
    }

    if (this.stepIndicators) {
      const dots = this.stepIndicators.querySelectorAll('.step-dot');
      dots.forEach((dot, index) => {
        if (index === newStep) {
          dot.classList.add('active');
        } else {
          dot.classList.remove('active');
        }
      });
    }
  }

  private announceStepChange(newStep: number): void {
    const announcement = document.createElement('div');
    announcement.setAttribute('aria-live', 'polite');
    announcement.setAttribute('aria-atomic', 'true');
    announcement.style.position = 'absolute';
    announcement.style.left = '-10000px';
    announcement.style.width = '1px';
    announcement.style.height = '1px';
    announcement.style.overflow = 'hidden';

    const stepDisplay = formatStepDisplay(newStep, this.stepCount, this.stepLabels);
    announcement.textContent = `Now on ${stepDisplay}`;

    document.body.appendChild(announcement);

    setTimeout(() => {
      document.body.removeChild(announcement);
    }, 1000);
  }

  public destroy(): void {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
  }
}