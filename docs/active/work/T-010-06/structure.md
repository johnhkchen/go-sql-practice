# Structure: Extract Inline Scripts to TypeScript

## Directory Structure

### New Directory
```
frontend/src/scripts/
├── syncViewer.ts
├── syncController.ts
├── presentationViewer.ts
├── presenterController.ts
├── statsController.ts
└── searchEnhancer.ts
```

## File Modifications

### 1. Create: `frontend/src/scripts/syncViewer.ts`
```typescript
import type { EventSource } from '../../types/global';

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

  constructor(sessionId: string, initialProgress: number, apiBase: string);
  init(): void;
  private setupDOM(): void;
  private updateConnectionStatus(status: string): void;
  private updateLastUpdate(text: string): void;
  private updateProgress(progress: number): void;
  private connect(): void;
  private handleConnect(): void;
  private handleMessage(event: MessageEvent): void;
  private handleError(): void;
}
```

### 2. Modify: `frontend/src/pages/sync/[id].astro`
- Remove lines 136-291 (entire class definition)
- Replace with:
```html
<script>
  import { SyncViewer } from '../../scripts/syncViewer';

  const sessionId = '<%= sessionId %>';
  const initialProgress = <%= initialProgress %>;
  const apiBase = import.meta.env.PUBLIC_API_URL || '';

  window.syncViewer = new SyncViewer(sessionId, initialProgress, apiBase);
</script>
```

### 3. Create: `frontend/src/scripts/syncController.ts`
```typescript
export class SyncController {
  private sessionId: string;
  private adminToken: string;
  private currentProgress: number;
  private apiBase: string;

  // DOM elements
  private progressSlider: HTMLInputElement | null;
  private progressValue: HTMLElement | null;
  private copyButton: HTMLButtonElement | null;
  private viewerUrlInput: HTMLInputElement | null;

  // State management
  private isUpdating: boolean;
  private throttleTimeout: number | null;
  private lastAnnouncedValue?: number;

  constructor(sessionId: string, adminToken: string, initialProgress: number, apiBase: string);
  init(): void;
  private setupDOM(): void;
  private setupEventListeners(): void;
  private handleSliderKeyboard(e: KeyboardEvent): void;
  private updateProgress(value: number): void;
  private updateProgressDisplay(value: number): void;
  private announceProgressChange(value: number): void;
  private validateProgress(value: number): number;
  private throttledUpdate(value: number): void;
  private async sendProgressUpdate(value: number): Promise<void>;
  private async copyToClipboard(): Promise<void>;
}
```

### 4. Modify: `frontend/src/pages/sync/[id]/control.astro`
- Remove lines 204-550+ (entire class definition)
- Replace with:
```html
<script>
  import { SyncController } from '../../../scripts/syncController';

  const sessionId = '<%= sessionId %>';
  const adminToken = '<%= adminToken %>';
  const initialProgress = <%= initialProgress %>;
  const apiBase = import.meta.env.PUBLIC_API_URL || '';

  new SyncController(sessionId, adminToken, initialProgress, apiBase);
</script>
```

### 5. Create: `frontend/src/scripts/presentationViewer.ts`
```typescript
import { progressToStep } from '../utils/stepConversion';

export class PresentationAutoViewer {
  private presentationId: string;
  private sessionId: string | null;
  private currentProgress: number;
  private stepCount: number;
  private stepLabels: string[];
  private apiBase: string;

  // State management
  private currentState: 'waiting' | 'starting' | 'live' | 'ended';
  private connectionStatus: 'connected' | 'connecting' | 'disconnected';
  private eventSource: EventSource | null;
  private currentStep: number;

  constructor(
    presentationId: string,
    sessionId: string | null,
    initialProgress: number,
    stepCount: number,
    stepLabels: string[],
    initialState: string,
    apiBase: string
  );

  private setupEventListeners(): void;
  private transitionToState(newState: string, transitionMessage?: string): void;
  private updateLastUpdate(message?: string): void;
  private updateConnectionStatus(status: string): void;
  private connect(): void;
  private handleConnect(): void;
  private async syncWithServerState(): Promise<void>;
  private handleMessage(event: MessageEvent): void;
  private handleError(): void;
  private updateProgress(progress: number): void;
}
```

### 6. Modify: `frontend/src/pages/watch/[id].astro`
- Remove lines 598-950+ (entire class definition)
- Import progressToStep at top of script
- Replace with:
```html
<script>
  import { PresentationAutoViewer } from '../../scripts/presentationViewer';

  const presentationId = '<%= presentationId %>';
  const sessionId = '<%= sessionId %>';
  const initialProgress = <%= initialProgress %>;
  const stepCount = <%= stepCount %>;
  const stepLabels = <%= JSON.stringify(stepLabels) %>;
  const initialState = '<%= initialState %>';
  const apiBase = import.meta.env.PUBLIC_API_URL || '';

  new PresentationAutoViewer(
    presentationId,
    sessionId,
    initialProgress,
    stepCount,
    stepLabels,
    initialState,
    apiBase
  );
</script>
```

### 7. Create: `frontend/src/scripts/presenterController.ts`
```typescript
import type { PresentationData } from '../types/api';
import {
  stepToProgress,
  progressToStep,
  getNextStep,
  getPreviousStep,
  formatStepDisplay
} from '../utils/stepConversion';

interface PresentationDataInternal {
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
  private presentationData: PresentationDataInternal;
  private apiBase: string;
  private currentStep: number;
  private currentProgress: number;
  private isUpdating: boolean = false;
  private throttleTimeout: number | null = null;

  // DOM elements (all typed)

  constructor(
    sessionId: string,
    adminToken: string,
    presentationData: PresentationDataInternal,
    apiBase: string
  );

  private init(): void;
  private setupDOM(): void;
  private validateDOMElements(): boolean;
  private setupEventListeners(): void;
  // ... other methods
}
```

### 8. Modify: `frontend/src/components/PresenterController.astro`
- Remove lines 804-1100+ (entire class definition)
- Keep the imports at the top
- Replace with:
```html
<script>
  import { PresenterController } from '../scripts/presenterController';

  const sessionId = '<%= sessionId %>';
  const adminToken = '<%= adminToken %>';
  const presentationData = <%= JSON.stringify(presentationData) %>;
  const apiBase = import.meta.env.PUBLIC_API_URL || '';

  new PresenterController(sessionId, adminToken, presentationData, apiBase);
</script>
```

### 9. Create: `frontend/src/scripts/statsController.ts`
```typescript
import type { StatsData, StatsState } from '../types/api';

function getApiBase(): string {
  if (typeof import.meta !== 'undefined' && import.meta.env?.PUBLIC_API_URL) {
    return import.meta.env.PUBLIC_API_URL;
  }
  if (typeof window !== 'undefined' && (window as any).PUBLIC_API_URL) {
    return (window as any).PUBLIC_API_URL;
  }
  return '';
}

export class StatsController {
  private state: StatsState;
  private container: HTMLElement | null;

  constructor();
  async init(): Promise<void>;
  private setupEventListeners(): void;
  private async fetchStats(): Promise<void>;
  private updateRefreshButton(loading: boolean): void;
  private isValidStatsData(data: any): data is StatsData;
  private setState(newState: Partial<StatsState>): void;
  private updateDOM(): void;
  private renderLoadingState(): void;
  private renderErrorState(error: string): void;
  private renderData(data: StatsData): void;
  // ... other methods
}
```

### 10. Modify: `frontend/src/components/StatsSummary.astro`
- Remove lines 53-380 (entire script block with class)
- Remove duplicate type definitions
- Replace with:
```html
<script>
  import { StatsController } from '../scripts/statsController';

  const controller = new StatsController();
  controller.init();

  // Global refresh function for retry buttons
  window.refreshStats = () => controller.refreshStats();
</script>
```

### 11. Create: `frontend/src/scripts/searchEnhancer.ts`
```typescript
import type { LinkItem, SearchResponse, SearchState } from '../types/api';

export class SearchEnhancer {
  private apiBase: string;
  private searchForm: HTMLFormElement | null;
  private searchInput: HTMLInputElement | null;
  private debounceTimeout: number | null = null;
  private currentController: AbortController | null = null;
  private state: SearchState;

  constructor(apiBase: string);

  private init(): void;
  private setupEventHandlers(): void;
  private async performSearch(query: string): Promise<void>;
  private updateURL(query: string): void;
}
```

### 12. Modify: `frontend/src/components/SearchInterface.astro`
- Remove lines 431-700+ (entire script block)
- Replace with simplified version:
```html
<script>
  import { SearchEnhancer } from '../scripts/searchEnhancer';

  const apiBase = import.meta.env.PUBLIC_API_URL || '';

  document.addEventListener('DOMContentLoaded', () => {
    new SearchEnhancer(apiBase);
  });
</script>
```

## Type Definitions Updates

### No changes needed to `frontend/src/types/api.ts`
- Already contains all required interfaces
- Classes will import from here

### Consider adding `frontend/src/types/global.d.ts`
```typescript
// Global type augmentations if needed
interface Window {
  syncViewer?: any;
  refreshStats?: () => void;
  PUBLIC_API_URL?: string;
}
```

## Import Path Conventions

- From scripts to types: `import type { ... } from '../types/api'`
- From scripts to utils: `import { ... } from '../utils/stepConversion'`
- From Astro to scripts: `import { ... } from '../../scripts/className'`
- From components to scripts: `import { ... } from '../scripts/className'`

## Build Considerations

- TypeScript will compile `.ts` files in `src/scripts/`
- Astro will bundle the imports in `<script>` tags
- No changes needed to `tsconfig.json` or build configuration
- Type checking will catch errors at build time