export class SyncViewer {
  private sessionId: string;
  private initialProgress: number;
  private apiBase: string;
  private currentProgress: number;

  // Connection state
  private connectionStatus: 'connected' | 'connecting' | 'disconnected';
  private eventSource: EventSource | null;
  private reconnectAttempts: number;
  private maxReconnectAttempts: number;

  // DOM elements
  private statusIndicator: HTMLElement | null;
  private statusText: HTMLElement | null;
  private lastUpdate: HTMLElement | null;
  private progressBar: HTMLProgressElement | null;
  private progressLabel: HTMLElement | null;

  constructor(sessionId: string, initialProgress: number, apiBase: string) {
    this.sessionId = sessionId;
    this.initialProgress = initialProgress;
    this.apiBase = apiBase;
    this.currentProgress = initialProgress;

    // Connection state
    this.connectionStatus = 'disconnected';
    this.eventSource = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;

    // DOM elements (will be set in init)
    this.statusIndicator = null;
    this.statusText = null;
    this.lastUpdate = null;
    this.progressBar = null;
    this.progressLabel = null;

    // Initialize when DOM is ready
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
    this.statusIndicator = document.getElementById('status-indicator');
    this.statusText = document.getElementById('status-text');
    this.lastUpdate = document.getElementById('last-update');
    this.progressBar = document.getElementById('progress-bar') as HTMLProgressElement;
    this.progressLabel = document.getElementById('progress-label');

    if (!this.statusIndicator || !this.statusText || !this.progressBar) {
      console.warn('SyncViewer: Required DOM elements not found');
      return;
    }

    // Set initial state
    this.updateConnectionStatus('connecting');
    this.updateLastUpdate('Initializing...');

    // Start SSE connection
    this.connect();
  }

  private updateConnectionStatus(status: 'connected' | 'connecting' | 'disconnected'): void {
    if (!this.statusIndicator || !this.statusText) return;

    this.connectionStatus = status;

    // Remove all status classes
    this.statusIndicator.classList.remove('connected', 'connecting', 'disconnected');

    // Add current status class
    this.statusIndicator.classList.add(status);

    // Update text
    switch (status) {
      case 'connected':
        this.statusText.textContent = 'Connected';
        break;
      case 'connecting':
        this.statusText.textContent = 'Connecting...';
        break;
      case 'disconnected':
        this.statusText.textContent = 'Disconnected';
        break;
      default:
        this.statusText.textContent = 'Unknown';
    }
  }

  private updateLastUpdate(text: string): void {
    if (this.lastUpdate) {
      this.lastUpdate.textContent = text;
    }
  }

  private updateProgress(progress: number): void {
    if (!this.progressBar || !this.progressLabel) return;

    this.currentProgress = progress;
    this.progressBar.value = progress;
    this.progressLabel.textContent = (progress * 100).toFixed(1) + '%';

    this.updateLastUpdate(`Updated ${new Date().toLocaleTimeString()}`);
  }

  private connect(): void {
    try {
      // Close existing connection if any
      if (this.eventSource) {
        this.eventSource.close();
      }

      // Create new EventSource connection
      this.eventSource = new EventSource(`${this.apiBase}/api/realtime`);

      // Set up event handlers
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
    this.reconnectAttempts = 0; // Reset reconnect counter
  }

  private handleMessage(event: MessageEvent): void {
    try {
      const data = JSON.parse(event.data);

      // Filter for sync_sessions collection updates to our session
      if (data.collection === 'sync_sessions' &&
          data.action === 'update' &&
          data.record &&
          data.record.id === this.sessionId) {

        this.updateProgress(data.record.progress);
      }
    } catch (error) {
      console.warn('Invalid SSE message:', error);
    }
  }

  private handleError(): void {
    this.updateConnectionStatus('disconnected');
    this.updateLastUpdate('Connection lost');
  }
}