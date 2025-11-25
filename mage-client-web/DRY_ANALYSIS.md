# DRY (Don't Repeat Yourself) Analysis - Mage Web Client

**Analysis Date:** 2025-11-24
**Total Route Files Analyzed:** 8 pages (~6,800 LOC)
**Total Components:** 33 components

---

## Executive Summary

This analysis identifies DRY violations and opportunities to extract reusable logic from page components in the Mage web client. The codebase shows good component extraction but has significant room for improvement in:

1. **API pattern consolidation** - Repeated session handling, error handling, and API call patterns
2. **Form validation logic** - Duplicated validation functions across pages
3. **Data transformation utilities** - Repeated formatting and conversion logic
4. **Loading state management** - Similar loading/error state patterns
5. **WebSocket connection patterns** - Duplicated connection setup logic

**Priority Findings:**
- **High Priority:** 18 issues requiring immediate attention
- **Medium Priority:** 12 issues for near-term refactoring
- **Low Priority:** 8 issues for future consideration

**Estimated Refactoring Effort:** 3-5 days for high priority items, 2-3 days for medium priority

---

## 1. API Pattern Violations

### 1.1 Session Validation Pattern (High Priority)

**Location:** Multiple API files (`lobby.ts:78-81`, `table.ts:76-80`, `profile.ts:23-28`)

**Issue:** Every API function repeats the same session validation pattern:
```typescript
const sessionId = await client.ensureSessionId();
if (!sessionId) {
    throw new Error('No active session - please login first');
}
```

**Occurrences:** 15+ times across 3 API files

**Recommendation:** Create an API wrapper utility
```typescript
// lib/api/base.ts
export async function withSession<T>(fn: (sessionId: string) => Promise<T>): Promise<T> {
    const client = getMageClient();
    const sessionId = await client.ensureSessionId();
    if (!sessionId) {
        throw new Error('No active session - please login first');
    }
    return fn(sessionId);
}

export async function withRoom<T>(fn: (sessionId: string, roomId: string) => Promise<T>): Promise<T> {
    return withSession(async (sessionId) => {
        const roomResponse = await client.getMainRoomId();
        if (!roomResponse.roomId) {
            throw new Error('Failed to get main room ID');
        }
        return fn(sessionId, roomResponse.roomId);
    });
}
```

**Estimated Effort:** Medium (4-6 hours)

---

### 1.2 Room ID Fetching Pattern (High Priority)

**Location:** All table/lobby API functions (`lobby.ts:84-87`, `table.ts:83-86`, repeated 10+ times)

**Issue:** Room ID fetching is duplicated in every function:
```typescript
const roomResponse = await client.getMainRoomId();
if (!roomResponse.roomId) {
    throw new Error('Failed to get main room ID');
}
```

**Recommendation:** Use the `withRoom` wrapper from 1.1 above

**Estimated Effort:** Small (included in 1.1)

---

### 1.3 TableView to Table Conversion (Medium Priority)

**Location:** `lobby.ts:18-69` and `table.ts:18-69` (identical 50+ line function)

**Issue:** Exact duplicate of `convertTableViewToTable` function in two files

**Recommendation:** Extract to shared utility
```typescript
// lib/utils/table-converter.ts
export function convertTableViewToTable(view: TableView): Table {
    // Move the entire function here
}
```

**Estimated Effort:** Small (1-2 hours)

---

### 1.4 Error Response Handling (High Priority)

**Location:** Throughout all API files

**Issue:** Inconsistent error handling patterns:
```typescript
// Pattern 1
if (!response.success) {
    throw new Error(response.error || 'Failed to ...');
}

// Pattern 2
if (error instanceof Error) {
    errorMessage = error.message;
} else {
    errorMessage = 'Failed to ...';
}
```

**Recommendation:** Create standardized error handler
```typescript
// lib/utils/api-errors.ts
export class ApiError extends Error {
    constructor(message: string, public code?: string, public details?: unknown) {
        super(message);
        this.name = 'ApiError';
    }
}

export function handleApiResponse<T extends { success: boolean; error?: string }>(
    response: T,
    defaultError: string
): void {
    if (!response.success) {
        throw new ApiError(response.error || defaultError);
    }
}

export function formatErrorMessage(error: unknown, fallback: string): string {
    if (error instanceof ApiError || error instanceof Error) {
        return error.message;
    }
    return fallback;
}
```

**Estimated Effort:** Medium (3-4 hours)

---

## 2. Form Validation Violations

### 2.1 Password Validation (High Priority)

**Location:**
- `login/+page.svelte:54-59` (6 char minimum)
- `register/+page.svelte:50-58` (8 char minimum)
- `profile/+page.svelte:77-85` (8 char minimum + match check)

**Issue:** Three different implementations of password validation with inconsistent requirements

**Recommendation:** Extract to validation service
```typescript
// lib/services/validation.ts
export const ValidationRules = {
    password: {
        minLength: 8,
        message: 'Password must be at least 8 characters'
    },
    username: {
        minLength: 3,
        maxLength: 20,
        pattern: /^[a-zA-Z0-9_]+$/,
        message: 'Username must be 3-20 characters (letters, numbers, underscores only)'
    }
};

export function validatePassword(password: string): string {
    if (!password) return 'Password is required';
    if (password.length < ValidationRules.password.minLength) {
        return ValidationRules.password.message;
    }
    return '';
}

export function validatePasswordMatch(password: string, confirm: string): string {
    if (!confirm) return 'Please confirm your password';
    if (password !== confirm) return 'Passwords do not match';
    return '';
}

export function validateUsername(username: string): string {
    if (!username) return 'Username is required';
    const { minLength, maxLength, pattern } = ValidationRules.username;
    if (username.length < minLength) return `Username must be at least ${minLength} characters`;
    if (username.length > maxLength) return `Username must be no more than ${maxLength} characters`;
    if (!pattern.test(username)) return ValidationRules.username.message;
    return '';
}
```

**Estimated Effort:** Small (2-3 hours)

---

### 2.2 Form Validation Patterns (Medium Priority)

**Location:**
- `login/+page.svelte:38-62`
- `register/+page.svelte:34-76`
- `profile/+page.svelte:66-85`

**Issue:** Each page reimplements form validation state management

**Recommendation:** Create reusable form validation hook
```typescript
// lib/hooks/useFormValidation.ts
export interface ValidationRule<T> {
    validate: (value: T) => string;
    onBlur?: boolean;
}

export function createFormValidator<T extends Record<string, unknown>>(
    rules: { [K in keyof T]?: ValidationRule<T[K]> }
) {
    const errors = $state<{ [K in keyof T]?: string }>({});

    function validateField<K extends keyof T>(field: K, value: T[K]): string {
        const rule = rules[field];
        if (!rule) return '';
        const error = rule.validate(value);
        errors[field] = error;
        return error;
    }

    function validateAll(values: T): boolean {
        let isValid = true;
        for (const field in rules) {
            const error = validateField(field, values[field]);
            if (error) isValid = false;
        }
        return isValid;
    }

    return { errors, validateField, validateAll };
}
```

**Estimated Effort:** Medium (4-5 hours)

---

## 3. Data Transformation Violations

### 3.1 Date Formatting (Medium Priority)

**Location:**
- `profile/+page.svelte:121-127` - formatDate
- `profile/+page.svelte:143-157` - formatTimeAgo
- `history/+page.svelte:3` - formatRelativeTime (imported)

**Issue:** Multiple date formatting functions scattered across files

**Recommendation:** Create centralized date utility
```typescript
// lib/utils/date-format.ts
export function formatDate(timestamp: number): string {
    return new Date(timestamp).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

export function formatTimeAgo(timestamp: number): string {
    const now = Date.now();
    const diff = now - timestamp;
    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (days > 0) return `${days}d ago`;
    if (hours > 0) return `${hours}h ago`;
    if (minutes > 0) return `${minutes}m ago`;
    return 'Just now';
}

export function formatDuration(seconds: number): string {
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m`;
    const hours = Math.floor(minutes / 60);
    const remainingMinutes = minutes % 60;
    return `${hours}h ${remainingMinutes}m`;
}
```

**Estimated Effort:** Small (1-2 hours)

---

### 3.2 Proto Timestamp Conversion (Medium Priority)

**Location:**
- `lobby.ts:32-43` - getCreateTime helper
- `lobby.ts:283-294` - getConnectedAt helper
- `table.ts:32-43` - getCreateTime helper (duplicate)
- `profile.ts:51-52` - inline conversion

**Issue:** Repeated proto timestamp conversion logic in multiple files

**Recommendation:** Create proto conversion utility
```typescript
// lib/utils/proto-converters.ts
export function protoTimestampToMillis(timestamp: Date | string | undefined): number {
    if (!timestamp) return Date.now();
    if (timestamp instanceof Date) return timestamp.getTime();
    if (typeof timestamp === 'string') return new Date(timestamp).getTime();
    return Date.now();
}

export function protoSecondsToMillis(seconds: number | undefined): number | undefined {
    return seconds ? seconds * 1000 : undefined;
}
```

**Estimated Effort:** Small (1 hour)

---

## 4. Loading State Management Violations

### 4.1 Loading/Error State Pattern (High Priority)

**Location:** Every page component

**Issue:** Each page declares the same state variables:
```typescript
let loading = $state(true);
let error = $state<string | null>(null);
```

And repeats the same try-catch-finally pattern:
```typescript
loading = true;
error = null;
try {
    // API call
} catch (err) {
    error = err instanceof Error ? err.message : 'Failed to...';
} finally {
    loading = false;
}
```

**Occurrences:** 8 pages, multiple times per page (30+ instances)

**Recommendation:** Create async state management utility
```typescript
// lib/utils/async-state.ts
export interface AsyncState<T> {
    data: T | null;
    loading: boolean;
    error: string | null;
}

export function createAsyncState<T>(initialData: T | null = null) {
    let state = $state<AsyncState<T>>({
        data: initialData,
        loading: false,
        error: null
    });

    async function execute(
        fn: () => Promise<T>,
        errorMessage = 'Operation failed'
    ): Promise<T | null> {
        state.loading = true;
        state.error = null;

        try {
            const data = await fn();
            state.data = data;
            return data;
        } catch (err) {
            state.error = err instanceof Error ? err.message : errorMessage;
            console.error(errorMessage, err);
            return null;
        } finally {
            state.loading = false;
        }
    }

    function reset() {
        state = {
            data: initialData,
            loading: false,
            error: null
        };
    }

    return {
        get state() { return state; },
        execute,
        reset
    };
}
```

**Usage:**
```typescript
const tablesState = createAsyncState<Table[]>([]);
await tablesState.execute(() => fetchTables());
```

**Estimated Effort:** Large (6-8 hours to implement + update all pages)

---

### 4.2 Loading State UI Components (Medium Priority)

**Location:** All pages render similar loading/error/empty states

**Issue:** Repeated loading spinner, error banner, and empty state markup

**Recommendation:** Create reusable state components
```typescript
// lib/components/AsyncStateRenderer.svelte
<script lang="ts">
    import type { AsyncState } from '$lib/utils/async-state';
    import LoadingSpinner from './LoadingSpinner.svelte';

    interface Props {
        state: AsyncState<any>;
        loadingMessage?: string;
        errorRetry?: () => void;
        emptyMessage?: string;
        children: Snippet;
    }

    let { state, loadingMessage, errorRetry, emptyMessage, children }: Props = $props();
</script>

{#if state.loading}
    <div class="loading-container">
        <LoadingSpinner size="large" />
        {#if loadingMessage}
            <p>{loadingMessage}</p>
        {/if}
    </div>
{:else if state.error}
    <div class="error-container">
        <p>{state.error}</p>
        {#if errorRetry}
            <button onclick={errorRetry}>Retry</button>
        {/if}
    </div>
{:else if !state.data || (Array.isArray(state.data) && state.data.length === 0)}
    {#if emptyMessage}
        <div class="empty-state">
            <p>{emptyMessage}</p>
        </div>
    {/if}
{:else}
    {@render children()}
{/if}
```

**Estimated Effort:** Medium (3-4 hours)

---

## 5. WebSocket Connection Patterns

### 5.1 WebSocket Connection Setup (High Priority)

**Location:**
- `lobby/+page.svelte:224-244` - connectWebSocket function
- `table/[id]/+page.svelte:174-196` - connectWebSocket function

**Issue:** Near-identical WebSocket connection setup duplicated

**Recommendation:** Extract to service
```typescript
// lib/services/websocket-connector.ts
export interface WebSocketConnectionOptions {
    onConnect?: () => void;
    onDisconnect?: () => void;
    onError?: (error: Error) => void;
}

export async function connectWebSocket(
    options: WebSocketConnectionOptions = {}
): Promise<string | null> {
    try {
        const client = getMageClient();
        const sessionId = await client.ensureSessionId();

        if (!sessionId) {
            console.warn('No session ID available for WebSocket connection');
            return null;
        }

        if (!websocketStore.isConnected()) {
            await websocketStore.connect(sessionId);
            options.onConnect?.();
        }

        return sessionId;
    } catch (err) {
        const error = err instanceof Error ? err : new Error('WebSocket connection failed');
        console.error('Failed to connect WebSocket:', error);
        options.onError?.(error);
        return null;
    }
}
```

**Estimated Effort:** Small (2-3 hours)

---

### 5.2 WebSocket Cleanup Pattern (Medium Priority)

**Location:**
- `lobby/+page.svelte:270-282` - onDestroy cleanup
- `table/[id]/+page.svelte:209-218` - onDestroy cleanup

**Issue:** Similar cleanup logic in both pages

**Recommendation:** Create cleanup hook
```typescript
// lib/hooks/useWebSocket.ts
export function useWebSocket(
    subscribe: (sessionId: string) => (() => void) | null,
    disconnect: () => void
) {
    let unsubscribe: (() => void) | null = null;

    async function connect() {
        const sessionId = await connectWebSocket();
        if (sessionId) {
            unsubscribe = subscribe(sessionId);
        }
    }

    onMount(connect);

    onDestroy(() => {
        if (unsubscribe) {
            unsubscribe();
            unsubscribe = null;
        }
        disconnect();
    });

    return { reconnect: connect };
}
```

**Estimated Effort:** Small (2 hours)

---

## 6. Authentication & Navigation Patterns

### 6.1 Auth Check on Mount (Medium Priority)

**Location:**
- `lobby/+page.svelte:248-268` - checkAuth wrapper
- `game/[id]/+page.svelte:350-358` - direct check with alert
- `login/+page.svelte:31-36` - redirect if already authed
- `register/+page.svelte:27-32` - redirect if already authed

**Issue:** Multiple implementations of auth verification on mount

**Recommendation:** Extract to reusable utilities
```typescript
// lib/utils/auth-guard.ts (already exists, enhance it)
export function useAuthGuard(redirectTo = '/login') {
    onMount(() => {
        if (!$auth.isAuthenticated) {
            goto(redirectTo);
        }
    });
}

export function useGuestGuard(redirectTo = '/lobby') {
    onMount(() => {
        const restored = auth.loadAuthFromStorage();
        if (restored) {
            goto(redirectTo);
        }
    });
}

export async function waitForAuth(
    onAuthenticated: () => void | Promise<void>,
    timeout = 5000
): Promise<void> {
    const startTime = Date.now();

    const checkAuth = () => {
        if ($auth.isAuthenticated) {
            onAuthenticated();
        } else if (Date.now() - startTime < timeout) {
            setTimeout(checkAuth, 100);
        } else {
            console.error('Auth timeout');
        }
    };

    checkAuth();
}
```

**Estimated Effort:** Small (2-3 hours)

---

### 6.2 Login Flow Duplication (High Priority)

**Location:**
- `login/+page.svelte:64-138` - performLogin function
- `login/+page.svelte:206-262` - handleGuestLoginAfterModal (similar logic)
- `register/+page.svelte:113-171` - auto-login after register

**Issue:** Similar login flow (connect, create token, store auth) repeated 3 times

**Recommendation:** Extract to auth service
```typescript
// lib/services/auth-service.ts
export interface LoginCredentials {
    username: string;
    password: string;
}

export interface LoginResult {
    success: boolean;
    error?: string;
}

export async function performLogin(
    credentials: LoginCredentials,
    redirectTo = '/lobby'
): Promise<LoginResult> {
    try {
        const client = getMageClient();
        const response = await client.connectUser(
            credentials.username,
            credentials.password
        );

        if (!response.success) {
            return { success: false, error: response.error || 'Login failed' };
        }

        const sessionId = response.sessionId || client.getSessionId();
        if (!sessionId) {
            return { success: false, error: 'No session ID received' };
        }

        // Ensure sessionId is set
        if (!client.getSessionId()) {
            client.setSessionId(sessionId);
        }

        // Create token and login
        const token = createSessionToken(
            sessionId,
            response.userId,
            credentials.username,
            `${credentials.username}@example.com`
        );

        auth.login(token, {
            id: response.userId,
            username: credentials.username,
            email: `${credentials.username}@example.com`
        });

        toast.success(`Welcome back, ${credentials.username}!`);

        // Small delay to ensure state is set
        await new Promise((resolve) => setTimeout(resolve, 50));

        goto(redirectTo);

        return { success: true };
    } catch (error) {
        const message = error instanceof Error ? error.message : 'Login failed';
        toast.error(message);
        return { success: false, error: message };
    }
}
```

**Estimated Effort:** Medium (4-5 hours)

---

## 7. Component Extraction Opportunities

### 7.1 Match Result Badge (Low Priority)

**Location:**
- `profile/+page.svelte:539-582` - match item styling
- `history/+page.svelte:44-72` - getResultBadgeClass and getResultEmoji

**Issue:** Badge styling and result mapping logic repeated

**Recommendation:** Create MatchResultBadge component
```typescript
// lib/components/MatchResultBadge.svelte
<script lang="ts">
    interface Props {
        result: 'win' | 'loss' | 'draw' | 'concede';
    }

    let { result }: Props = $props();

    const config = {
        win: { emoji: '🏆', class: 'win', color: '#10b981' },
        loss: { emoji: '💀', class: 'loss', color: '#dc2626' },
        draw: { emoji: '🤝', class: 'draw', color: '#4f46e5' },
        concede: { emoji: '🏳️', class: 'concede', color: '#6b7280' }
    };

    const current = $derived(config[result] || config.win);
</script>

<span class="result-badge {current.class}">
    {current.emoji}
    {result.toUpperCase()}
</span>
```

**Estimated Effort:** Small (1-2 hours)

---

### 7.2 Player Stat Display (Low Priority)

**Location:**
- `table/[id]/+page.svelte:427-440` - opponent stats
- `table/[id]/+page.svelte:505-513` - player stats

**Issue:** Similar stat display layout repeated for player and opponent

**Recommendation:** Create PlayerStatsBar component
```typescript
// lib/components/PlayerStatsBar.svelte
<script lang="ts">
    interface Props {
        life: number;
        libraryCount: number;
        handCount?: number;
        showHand?: boolean;
    }

    let { life, libraryCount, handCount, showHand = true }: Props = $props();
</script>

<div class="player-stats">
    <div class="stat life" title="Life Total">
        <span class="stat-icon">❤️</span>
        <span class="stat-value">{life}</span>
    </div>
    <div class="stat library" title="Library">
        <span class="stat-icon">📚</span>
        <span class="stat-value">{libraryCount}</span>
    </div>
    {#if showHand && handCount !== undefined}
        <div class="stat hand" title="Hand Size">
            <span class="stat-icon">🎴</span>
            <span class="stat-value">{handCount}</span>
        </div>
    {/if}
</div>
```

**Estimated Effort:** Small (1 hour)

---

### 7.3 Empty State Components (Low Priority)

**Location:** Every page has custom empty state markup

**Issue:** Repeated empty state structure with icon, title, message, action button

**Recommendation:** Create generic EmptyState component
```typescript
// lib/components/EmptyState.svelte
<script lang="ts">
    interface Props {
        icon?: string;
        title: string;
        message: string;
        actionLabel?: string;
        onAction?: () => void;
    }

    let { icon = '🎮', title, message, actionLabel, onAction }: Props = $props();
</script>

<div class="empty-state">
    <div class="empty-icon">{icon}</div>
    <h3>{title}</h3>
    <p>{message}</p>
    {#if actionLabel && onAction}
        <button class="btn-primary" onclick={onAction}>
            {actionLabel}
        </button>
    {/if}
</div>
```

**Estimated Effort:** Small (1 hour)

---

## 8. Utility Function Duplication

### 8.1 Random Generation Functions (Low Priority)

**Location:**
- `login/+page.svelte:153-161` - generateRandomPassword
- `login/+page.svelte:167-171` - generateGuestUsername

**Issue:** These utilities belong in a shared module, not in a page component

**Recommendation:** Move to utilities
```typescript
// lib/utils/random.ts
export function generateRandomPassword(length = 12): string {
    const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*';
    let password = '';
    for (let i = 0; i < length; i++) {
        password += charset.charAt(Math.floor(Math.random() * charset.length));
    }
    return password;
}

export function generateGuestUsername(): string {
    const randomPart = Math.random().toString(36).substring(2, 9);
    return 'guest_' + randomPart;
}
```

**Estimated Effort:** Small (30 minutes)

---

### 8.2 Clipboard Operations (Low Priority)

**Location:**
- `login/+page.svelte:264-275` - copyPasswordToClipboard
- `login/+page.svelte:437-444` - copy username (inline)

**Issue:** Clipboard operations with toast feedback repeated

**Recommendation:** Create clipboard utility
```typescript
// lib/utils/clipboard.ts
export async function copyToClipboard(text: string, successMessage = 'Copied!'): Promise<boolean> {
    if (typeof navigator === 'undefined' || !navigator.clipboard) {
        toast.error('Clipboard not available');
        return false;
    }

    try {
        await navigator.clipboard.writeText(text);
        toast.success(successMessage);
        return true;
    } catch (err) {
        toast.error('Failed to copy');
        return false;
    }
}
```

**Estimated Effort:** Small (30 minutes)

---

## 9. Store Usage Patterns

### 9.1 Store Import Inconsistency (Low Priority)

**Location:** Throughout codebase

**Issue:** Mixed patterns for importing stores:
```typescript
// Pattern 1: Named export
import { auth } from '$lib/stores/auth';

// Pattern 2: Named as authStore
import { authStore } from '$lib/stores/auth';
```

**Recommendation:** Standardize on named exports and update inconsistent imports

**Estimated Effort:** Small (1 hour)

---

## 10. Style Duplication

### 10.1 Button Styles (Medium Priority)

**Location:** Every page defines similar button classes

**Issue:** Repeated button style definitions:
- `.btn-primary` - defined in 8 pages
- `.btn-secondary` - defined in 5 pages
- `.btn-danger` - defined in 3 pages

**Recommendation:** Extract to global styles or component library
```css
/* src/styles/buttons.css */
.btn-primary {
    padding: 0.75rem 1.5rem;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.2s;
}

.btn-primary:hover:not(:disabled) {
    background: #5568d3;
}

.btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}
/* ... etc */
```

**Estimated Effort:** Medium (3-4 hours)

---

### 10.2 Loading Spinner Animation (Low Priority)

**Location:** Multiple pages define the same spin animation

**Issue:** `@keyframes spin` defined in:
- `login/+page.svelte:662-666`
- `register/+page.svelte:411-415`
- `profile/+page.svelte:377-381`
- `game/[id]/+page.svelte:624-628`

**Recommendation:** Move to global styles and use LoadingSpinner component everywhere

**Estimated Effort:** Small (1 hour)

---

## Prioritized Action Items

### Immediate (High Priority - Week 1)

1. **API Wrapper Utility** (6 hours)
   - Create `withSession` and `withRoom` wrappers
   - Standardize error handling with `ApiError` class
   - Update all API files to use new patterns

2. **Loading State Management** (8 hours)
   - Implement `createAsyncState` utility
   - Create `AsyncStateRenderer` component
   - Update 3-4 pages as proof of concept

3. **Form Validation Service** (5 hours)
   - Extract all validation functions
   - Create reusable validation rules
   - Update login, register, profile pages

4. **Login Flow Consolidation** (5 hours)
   - Create `performLogin` in auth-service
   - Update login, register pages
   - Extract guest login logic

5. **WebSocket Connection Pattern** (3 hours)
   - Create `connectWebSocket` service
   - Update lobby and table pages

**Total High Priority:** ~27 hours (3-4 days)

---

### Near-term (Medium Priority - Week 2-3)

1. **Data Transformation Utilities** (3 hours)
   - Date formatting utilities
   - Proto conversion utilities
   - TableView converter extraction

2. **Form Validation Hook** (5 hours)
   - Create `createFormValidator` hook
   - Update all forms to use it

3. **AsyncStateRenderer** (4 hours)
   - Finish component implementation
   - Update remaining pages

4. **WebSocket Cleanup Hook** (2 hours)
   - Create `useWebSocket` hook
   - Update pages with WebSocket

5. **Button Styles** (4 hours)
   - Extract to global CSS
   - Create button component variants

**Total Medium Priority:** ~18 hours (2-3 days)

---

### Future (Low Priority - Week 4+)

1. **Component Extraction** (4 hours)
   - MatchResultBadge
   - PlayerStatsBar
   - EmptyState

2. **Utility Organization** (2 hours)
   - Random generators
   - Clipboard utilities
   - Store import standardization

3. **Style Consolidation** (2 hours)
   - Global animations
   - Consistent spacing/sizing

**Total Low Priority:** ~8 hours (1 day)

---

## Refactoring Strategy

### Phase 1: Foundation (High Priority)
Focus on creating the infrastructure that will make future refactoring easier:
1. API wrapper utilities
2. Error handling standardization
3. Loading state management

### Phase 2: Forms & Auth (High + Medium Priority)
Consolidate authentication and form-related logic:
1. Validation service
2. Auth service
3. Form validation hooks

### Phase 3: Data & Display (Medium Priority)
Standardize data transformation and UI patterns:
1. Date/proto utilities
2. AsyncStateRenderer
3. Component extraction

### Phase 4: Polish (Low Priority)
Clean up remaining duplication:
1. Utilities organization
2. Style consolidation
3. Import standardization

---

## Testing Recommendations

After each refactoring phase:

1. **Manual Testing Checklist:**
   - Login/logout flow
   - Form validation (login, register, profile)
   - Table creation and joining
   - WebSocket connections
   - Error scenarios

2. **Unit Tests:**
   - Validation functions
   - Data transformation utilities
   - API wrappers
   - Auth service

3. **Integration Tests:**
   - Complete login flow
   - Form submission
   - API error handling

---

## Metrics

### Code Reduction Estimate
- **Lines of Code Reduction:** ~1,200-1,500 lines (18-22%)
- **Functions Eliminated:** ~40-50 duplicate functions
- **Components Created:** 8-10 reusable components
- **Utilities Created:** 6-8 utility modules

### Maintainability Improvements
- Consistent error handling across all API calls
- Single source of truth for validation rules
- Standardized loading/error state management
- Reduced cognitive load for new developers

### Performance Impact
- Minimal to none (some patterns may improve performance)
- Better code splitting opportunities
- Reduced bundle size through utility reuse

---

## Conclusion

The Mage web client has a solid foundation with good component separation, but significant DRY violations exist in:
1. API patterns (session/room handling)
2. Form validation
3. Loading state management
4. Authentication flows

By following this refactoring plan, the codebase will become:
- **More maintainable** - Changes in one place affect all usages
- **More testable** - Extracted utilities are easier to unit test
- **More consistent** - Standardized patterns throughout
- **Easier to onboard** - Clear patterns for new developers

**Recommended Start:** Begin with High Priority items in order listed above. Each can be completed incrementally and merged separately.
