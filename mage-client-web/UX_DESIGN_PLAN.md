# UX Design Plan: Mage Web Client

## Design Vision

A dark, immersive interface inspired by Magic: The Gathering's visual identity. The design should evoke the feeling of an ancient library or wizard's sanctum - sophisticated, mysterious, and powerful.

---

## 1. Design System

### 1.1 Color Palette

```css
:root {
  /* Core Background Hierarchy */
  --bg-void: #0a0b0d;           /* Deepest black - page background */
  --bg-obsidian: #12141a;       /* Primary surface */
  --bg-slate: #1a1d26;          /* Elevated surface */
  --bg-iron: #242833;           /* Interactive surface */
  --bg-steel: #2e3340;          /* Hover state */

  /* Text Hierarchy */
  --text-bright: #f4f4f5;       /* Primary text */
  --text-muted: #a1a1aa;        /* Secondary text */
  --text-dim: #71717a;          /* Tertiary/disabled text */
  --text-ghost: #52525b;        /* Placeholder text */

  /* Accent - Arcane Gold (primary brand) */
  --accent-gold: #c9a227;       /* Primary accent */
  --accent-gold-bright: #e4b82a; /* Hover state */
  --accent-gold-dim: #9a7b1e;   /* Pressed/muted */
  --accent-gold-glow: rgba(201, 162, 39, 0.15); /* Glow effect */

  /* Semantic Colors */
  --status-success: #22c55e;
  --status-warning: #f59e0b;
  --status-error: #ef4444;
  --status-info: #3b82f6;

  /* MTG Mana Colors */
  --mana-white: #f8f6e3;
  --mana-blue: #0e68ab;
  --mana-black: #150b00;
  --mana-red: #d3202a;
  --mana-green: #00733e;
  --mana-colorless: #9ca3af;

  /* Table Status Colors */
  --table-waiting: #f59e0b;
  --table-ready: #22c55e;
  --table-playing: #3b82f6;
  --table-finished: #71717a;

  /* Borders & Dividers */
  --border-subtle: #27272a;
  --border-default: #3f3f46;
  --border-strong: #52525b;
  --border-accent: var(--accent-gold-dim);

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.4);
  --shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.5);
  --shadow-glow: 0 0 20px var(--accent-gold-glow);
}
```

### 1.2 Typography Scale

```css
:root {
  /* Font Families */
  --font-display: 'Cinzel', 'Times New Roman', serif;  /* Headers, branding */
  --font-body: 'Inter', -apple-system, sans-serif;     /* Body text, UI */
  --font-mono: 'JetBrains Mono', monospace;            /* Code, numbers */

  /* Font Sizes */
  --text-xs: 0.75rem;     /* 12px - Tiny labels */
  --text-sm: 0.875rem;    /* 14px - Secondary text */
  --text-base: 1rem;      /* 16px - Body text */
  --text-lg: 1.125rem;    /* 18px - Emphasized body */
  --text-xl: 1.25rem;     /* 20px - Small headers */
  --text-2xl: 1.5rem;     /* 24px - Section headers */
  --text-3xl: 1.875rem;   /* 30px - Page headers */
  --text-4xl: 2.25rem;    /* 36px - Hero text */

  /* Font Weights */
  --weight-normal: 400;
  --weight-medium: 500;
  --weight-semibold: 600;
  --weight-bold: 700;

  /* Line Heights */
  --leading-tight: 1.25;
  --leading-normal: 1.5;
  --leading-relaxed: 1.75;
}
```

### 1.3 Spacing System

```css
:root {
  --space-1: 0.25rem;   /* 4px */
  --space-2: 0.5rem;    /* 8px */
  --space-3: 0.75rem;   /* 12px */
  --space-4: 1rem;      /* 16px */
  --space-5: 1.25rem;   /* 20px */
  --space-6: 1.5rem;    /* 24px */
  --space-8: 2rem;      /* 32px */
  --space-10: 2.5rem;   /* 40px */
  --space-12: 3rem;     /* 48px */
  --space-16: 4rem;     /* 64px */
}
```

### 1.4 Border Radius

```css
:root {
  --radius-sm: 0.25rem;   /* 4px - Small elements */
  --radius-md: 0.5rem;    /* 8px - Buttons, inputs */
  --radius-lg: 0.75rem;   /* 12px - Cards, panels */
  --radius-xl: 1rem;      /* 16px - Modals */
  --radius-full: 9999px;  /* Pills, avatars */
}
```

### 1.5 Transitions

```css
:root {
  --transition-fast: 150ms ease;
  --transition-base: 200ms ease;
  --transition-slow: 300ms ease;
}
```

---

## 2. Base UI Components

### 2.1 Component Directory Structure

```
src/lib/components/
├── ui/                          # Base UI primitives
│   ├── Button.svelte
│   ├── Input.svelte
│   ├── Select.svelte
│   ├── Textarea.svelte
│   ├── Checkbox.svelte
│   ├── Badge.svelte
│   ├── Avatar.svelte
│   ├── Tooltip.svelte
│   ├── Tabs.svelte
│   ├── Panel.svelte
│   ├── Divider.svelte
│   └── Icon.svelte
│
├── layout/                      # Layout components
│   ├── PageHeader.svelte
│   ├── Sidebar.svelte
│   ├── Grid.svelte
│   └── Stack.svelte
│
├── feedback/                    # Feedback components
│   ├── Modal.svelte             # Existing, refactor
│   ├── Toast.svelte             # Existing, refactor
│   ├── LoadingSpinner.svelte    # Existing, refactor
│   ├── Skeleton.svelte          # New
│   ├── Alert.svelte             # New
│   └── Progress.svelte          # New
│
├── data/                        # Data display components
│   ├── Table.svelte
│   ├── List.svelte
│   └── EmptyState.svelte
│
├── mtg/                         # MTG-specific components
│   ├── ManaSymbol.svelte
│   ├── ManaCost.svelte
│   ├── CardFrame.svelte
│   ├── FormatBadge.svelte
│   └── PlayerAvatar.svelte
│
├── lobby/                       # Lobby view components
│   ├── TableCard.svelte         # Existing, refactor
│   ├── TableFilters.svelte      # New
│   ├── PlayerList.svelte        # Refactor from OnlinePlayersList
│   └── LobbyStats.svelte        # New
│
├── deck/                        # Deck builder components
│   ├── DeckCard.svelte          # Existing, refactor
│   ├── DeckList.svelte          # New
│   ├── DeckStats.svelte         # New
│   ├── CardSearch.svelte        # New
│   ├── ManaCurve.svelte         # New
│   └── DeckImport.svelte        # New
│
├── table/                       # Table view components
│   ├── SeatLayout.svelte        # New
│   ├── PlayerSeat.svelte        # New
│   ├── TableInfo.svelte         # New
│   └── ReadyIndicator.svelte    # New
│
├── chat/                        # Chat components
│   ├── ChatPanel.svelte         # Unified chat
│   ├── ChatMessage.svelte       # New
│   └── ChatInput.svelte         # New
│
└── game/                        # Game view components (existing)
    └── ...
```

---

### 2.2 Button Component

**File:** `src/lib/components/ui/Button.svelte`

```svelte
<script lang="ts">
  interface Props {
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
    size?: 'sm' | 'md' | 'lg';
    disabled?: boolean;
    loading?: boolean;
    fullWidth?: boolean;
    type?: 'button' | 'submit';
    onclick?: (e: MouseEvent) => void;
  }

  let {
    variant = 'primary',
    size = 'md',
    disabled = false,
    loading = false,
    fullWidth = false,
    type = 'button',
    onclick,
    children
  }: Props = $props();
</script>

<button
  {type}
  class="btn btn-{variant} btn-{size}"
  class:btn-full={fullWidth}
  class:btn-loading={loading}
  disabled={disabled || loading}
  {onclick}
>
  {#if loading}
    <span class="btn-spinner"></span>
  {/if}
  <span class="btn-content" class:invisible={loading}>
    {@render children?.()}
  </span>
</button>

<style>
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-2);
    font-family: var(--font-body);
    font-weight: var(--weight-medium);
    border: 1px solid transparent;
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all var(--transition-fast);
    position: relative;
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Variants */
  .btn-primary {
    background: var(--accent-gold);
    color: var(--bg-void);
    border-color: var(--accent-gold);
  }

  .btn-primary:hover:not(:disabled) {
    background: var(--accent-gold-bright);
    box-shadow: var(--shadow-glow);
  }

  .btn-secondary {
    background: var(--bg-iron);
    color: var(--text-bright);
    border-color: var(--border-default);
  }

  .btn-secondary:hover:not(:disabled) {
    background: var(--bg-steel);
    border-color: var(--border-strong);
  }

  .btn-ghost {
    background: transparent;
    color: var(--text-muted);
  }

  .btn-ghost:hover:not(:disabled) {
    background: var(--bg-iron);
    color: var(--text-bright);
  }

  .btn-danger {
    background: var(--status-error);
    color: white;
    border-color: var(--status-error);
  }

  .btn-danger:hover:not(:disabled) {
    background: #dc2626;
  }

  /* Sizes */
  .btn-sm {
    padding: var(--space-1) var(--space-3);
    font-size: var(--text-sm);
  }

  .btn-md {
    padding: var(--space-2) var(--space-4);
    font-size: var(--text-base);
  }

  .btn-lg {
    padding: var(--space-3) var(--space-6);
    font-size: var(--text-lg);
  }

  .btn-full {
    width: 100%;
  }

  /* Loading spinner */
  .btn-spinner {
    position: absolute;
    width: 1em;
    height: 1em;
    border: 2px solid currentColor;
    border-top-color: transparent;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }

  .invisible {
    visibility: hidden;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
```

---

### 2.3 Input Component

**File:** `src/lib/components/ui/Input.svelte`

```svelte
<script lang="ts">
  interface Props {
    type?: 'text' | 'password' | 'email' | 'search' | 'number';
    value?: string;
    placeholder?: string;
    label?: string;
    error?: string;
    disabled?: boolean;
    required?: boolean;
    id?: string;
  }

  let {
    type = 'text',
    value = $bindable(''),
    placeholder = '',
    label = '',
    error = '',
    disabled = false,
    required = false,
    id = crypto.randomUUID()
  }: Props = $props();
</script>

<div class="input-group" class:has-error={error}>
  {#if label}
    <label for={id} class="input-label">
      {label}
      {#if required}<span class="required">*</span>{/if}
    </label>
  {/if}

  <input
    {id}
    {type}
    bind:value
    {placeholder}
    {disabled}
    {required}
    class="input"
  />

  {#if error}
    <span class="input-error">{error}</span>
  {/if}
</div>

<style>
  .input-group {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }

  .input-label {
    font-size: var(--text-sm);
    font-weight: var(--weight-medium);
    color: var(--text-muted);
  }

  .required {
    color: var(--status-error);
    margin-left: var(--space-1);
  }

  .input {
    width: 100%;
    padding: var(--space-2) var(--space-3);
    font-family: var(--font-body);
    font-size: var(--text-base);
    color: var(--text-bright);
    background: var(--bg-iron);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-md);
    outline: none;
    transition: all var(--transition-fast);
  }

  .input::placeholder {
    color: var(--text-ghost);
  }

  .input:focus {
    border-color: var(--accent-gold);
    box-shadow: 0 0 0 2px var(--accent-gold-glow);
  }

  .input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .has-error .input {
    border-color: var(--status-error);
  }

  .has-error .input:focus {
    box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.15);
  }

  .input-error {
    font-size: var(--text-sm);
    color: var(--status-error);
  }
</style>
```

---

### 2.4 Badge Component

**File:** `src/lib/components/ui/Badge.svelte`

```svelte
<script lang="ts">
  interface Props {
    variant?: 'default' | 'success' | 'warning' | 'error' | 'info' | 'muted';
    size?: 'sm' | 'md';
  }

  let {
    variant = 'default',
    size = 'md',
    children
  }: Props = $props();
</script>

<span class="badge badge-{variant} badge-{size}">
  {@render children?.()}
</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    font-family: var(--font-body);
    font-weight: var(--weight-medium);
    border-radius: var(--radius-full);
    white-space: nowrap;
  }

  .badge-sm {
    padding: 0.125rem var(--space-2);
    font-size: var(--text-xs);
  }

  .badge-md {
    padding: var(--space-1) var(--space-3);
    font-size: var(--text-sm);
  }

  .badge-default {
    background: var(--bg-iron);
    color: var(--text-muted);
  }

  .badge-success {
    background: rgba(34, 197, 94, 0.15);
    color: var(--status-success);
  }

  .badge-warning {
    background: rgba(245, 158, 11, 0.15);
    color: var(--status-warning);
  }

  .badge-error {
    background: rgba(239, 68, 68, 0.15);
    color: var(--status-error);
  }

  .badge-info {
    background: rgba(59, 130, 246, 0.15);
    color: var(--status-info);
  }

  .badge-muted {
    background: var(--bg-slate);
    color: var(--text-dim);
  }
</style>
```

---

### 2.5 Panel Component

**File:** `src/lib/components/ui/Panel.svelte`

```svelte
<script lang="ts">
  interface Props {
    title?: string;
    padding?: 'none' | 'sm' | 'md' | 'lg';
    variant?: 'default' | 'elevated' | 'bordered';
  }

  let {
    title = '',
    padding = 'md',
    variant = 'default',
    children
  }: Props = $props();
</script>

<div class="panel panel-{variant} panel-pad-{padding}">
  {#if title}
    <div class="panel-header">
      <h3 class="panel-title">{title}</h3>
    </div>
  {/if}
  <div class="panel-content">
    {@render children?.()}
  </div>
</div>

<style>
  .panel {
    border-radius: var(--radius-lg);
    overflow: hidden;
  }

  .panel-default {
    background: var(--bg-obsidian);
  }

  .panel-elevated {
    background: var(--bg-slate);
    box-shadow: var(--shadow-md);
  }

  .panel-bordered {
    background: var(--bg-obsidian);
    border: 1px solid var(--border-subtle);
  }

  .panel-header {
    padding: var(--space-4) var(--space-4) 0;
  }

  .panel-title {
    font-family: var(--font-display);
    font-size: var(--text-lg);
    font-weight: var(--weight-semibold);
    color: var(--text-bright);
    margin: 0;
  }

  .panel-pad-none .panel-content { padding: 0; }
  .panel-pad-sm .panel-content { padding: var(--space-3); }
  .panel-pad-md .panel-content { padding: var(--space-4); }
  .panel-pad-lg .panel-content { padding: var(--space-6); }

  .panel-pad-none .panel-header { padding: var(--space-3) var(--space-3) 0; }
  .panel-pad-sm .panel-header { padding: var(--space-3) var(--space-3) 0; }
  .panel-pad-md .panel-header { padding: var(--space-4) var(--space-4) 0; }
  .panel-pad-lg .panel-header { padding: var(--space-6) var(--space-6) 0; }
</style>
```

---

### 2.6 MTG Format Badge

**File:** `src/lib/components/mtg/FormatBadge.svelte`

```svelte
<script lang="ts">
  interface Props {
    format: string;
    size?: 'sm' | 'md';
  }

  let { format, size = 'md' }: Props = $props();

  const formatColors: Record<string, string> = {
    'standard': '#f59e0b',
    'pioneer': '#8b5cf6',
    'modern': '#ef4444',
    'legacy': '#3b82f6',
    'vintage': '#10b981',
    'commander': '#ec4899',
    'edh': '#ec4899',
    'pauper': '#71717a',
    'limited': '#06b6d4',
    'draft': '#06b6d4',
    'sealed': '#14b8a6'
  };

  const color = $derived(formatColors[format.toLowerCase()] || '#71717a');
</script>

<span
  class="format-badge format-badge-{size}"
  style="--format-color: {color}"
>
  {format}
</span>

<style>
  .format-badge {
    display: inline-flex;
    align-items: center;
    font-family: var(--font-body);
    font-weight: var(--weight-semibold);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    border-radius: var(--radius-sm);
    background: color-mix(in srgb, var(--format-color) 15%, transparent);
    color: var(--format-color);
    border: 1px solid color-mix(in srgb, var(--format-color) 30%, transparent);
  }

  .format-badge-sm {
    padding: 0.125rem var(--space-2);
    font-size: 0.625rem;
  }

  .format-badge-md {
    padding: var(--space-1) var(--space-2);
    font-size: var(--text-xs);
  }
</style>
```

---

### 2.7 Mana Symbol Component

**File:** `src/lib/components/mtg/ManaSymbol.svelte`

```svelte
<script lang="ts">
  interface Props {
    symbol: 'W' | 'U' | 'B' | 'R' | 'G' | 'C' | number;
    size?: 'sm' | 'md' | 'lg';
  }

  let { symbol, size = 'md' }: Props = $props();

  const colors: Record<string, { bg: string; text: string }> = {
    'W': { bg: 'var(--mana-white)', text: '#1a1a1a' },
    'U': { bg: 'var(--mana-blue)', text: '#ffffff' },
    'B': { bg: 'var(--mana-black)', text: '#aaaaaa' },
    'R': { bg: 'var(--mana-red)', text: '#ffffff' },
    'G': { bg: 'var(--mana-green)', text: '#ffffff' },
    'C': { bg: 'var(--mana-colorless)', text: '#1a1a1a' }
  };

  const isNumeric = $derived(typeof symbol === 'number');
  const style = $derived(
    isNumeric
      ? colors['C']
      : colors[symbol as string] || colors['C']
  );
</script>

<span
  class="mana-symbol mana-symbol-{size}"
  style="background: {style.bg}; color: {style.text}"
>
  {symbol}
</span>

<style>
  .mana-symbol {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-family: var(--font-body);
    font-weight: var(--weight-bold);
    border-radius: var(--radius-full);
    box-shadow: inset 0 -2px 0 rgba(0, 0, 0, 0.2);
  }

  .mana-symbol-sm {
    width: 1rem;
    height: 1rem;
    font-size: 0.625rem;
  }

  .mana-symbol-md {
    width: 1.25rem;
    height: 1.25rem;
    font-size: 0.75rem;
  }

  .mana-symbol-lg {
    width: 1.5rem;
    height: 1.5rem;
    font-size: 0.875rem;
  }
</style>
```

---

## 3. Lobby View Rework

### 3.1 Layout Structure

```
+------------------------------------------------------------------+
|  [Logo]           MAGE            [Connection] [User Menu]       |
+------------------------------------------------------------------+
|                                                                  |
|  +------------------------+  +--------------------------------+  |
|  |     ONLINE PLAYERS     |  |         ACTIVE TABLES          |  |
|  |     (Collapsible)      |  |                                |  |
|  |                        |  |  [Filters: Format | Search]    |  |
|  |  [Player List]         |  |  [Create Table]                |  |
|  |  - Player 1            |  |                                |  |
|  |  - Player 2            |  |  +---------------------------+ |  |
|  |  - Player 3            |  |  | Table Card               | |  |
|  |  ...                   |  |  | Format: Modern            | |  |
|  |                        |  |  | Host: Player1             | |  |
|  +------------------------+  |  | 2/4 Players | Waiting     | |  |
|                              |  +---------------------------+ |  |
|  +------------------------+  |                                |  |
|  |      LOBBY CHAT        |  |  +---------------------------+ |  |
|  |                        |  |  | Table Card               | |  |
|  |  [Message List]        |  |  | ...                       | |  |
|  |                        |  |  +---------------------------+ |  |
|  |  [Input] [Send]        |  |                                |  |
|  +------------------------+  +--------------------------------+  |
|                                                                  |
+------------------------------------------------------------------+
```

### 3.2 Key Components

#### TableFilters.svelte
```svelte
<!-- Filter bar with format dropdown, search input, open-only toggle -->
<script lang="ts">
  interface Props {
    format: string;
    search: string;
    openOnly: boolean;
    formats: string[];
    onchange: (filters: { format: string; search: string; openOnly: boolean }) => void;
  }
</script>
```

#### TableCard.svelte (Reworked)
- Cleaner layout with format badge
- Player slots visualization (filled/empty)
- Status indicator (waiting/playing/finished)
- Host indicator
- Password lock icon if protected
- Hover state with gold accent border

#### PlayerList.svelte
- Compact list with online indicators
- Status dot (green = active, yellow = in-game, gray = idle)
- Click to view profile tooltip

#### LobbyStats.svelte
- Total players online
- Active games count
- Tables waiting count

---

## 4. Deck Builder View Rework

### 4.1 Layout Structure

```
+------------------------------------------------------------------+
|  [Logo]           MAGE            [Connection] [User Menu]       |
+------------------------------------------------------------------+
|                                                                  |
|  MY DECKS                                      [Import Deck]     |
|                                                                  |
|  [Format Filter: All | Standard | Modern | Commander | ...]      |
|                                                                  |
|  +-------------------+  +-------------------+  +---------------+ |
|  |                   |  |                   |  |               | |
|  |   [Format Badge]  |  |   [Format Badge]  |  |  + NEW DECK   | |
|  |                   |  |                   |  |               | |
|  |   Deck Name       |  |   Deck Name       |  |  Click to     | |
|  |                   |  |                   |  |  create       | |
|  |   60 cards        |  |   100 cards       |  |               | |
|  |   [Mana Curve]    |  |   [Mana Curve]    |  |               | |
|  |                   |  |                   |  |               | |
|  |   [Edit] [Delete] |  |   [Edit] [Delete] |  |               | |
|  +-------------------+  +-------------------+  +---------------+ |
|                                                                  |
+------------------------------------------------------------------+
```

### 4.2 Deck Detail View

```
+------------------------------------------------------------------+
|  < Back to Decks       DECK NAME               [Export] [Delete] |
+------------------------------------------------------------------+
|                                                                  |
|  +---------------------------+  +-----------------------------+  |
|  |      DECK STATS           |  |       MANA CURVE            |  |
|  |  Format: Commander        |  |                             |  |
|  |  Cards: 100               |  |       [Bar Chart]           |  |
|  |  Colors: WUB              |  |    0 1 2 3 4 5 6 7+         |  |
|  +---------------------------+  +-----------------------------+  |
|                                                                  |
|  +------------------------------------------------------------+  |
|  |  COMMANDER                                                  |  |
|  |  +--------+                                                 |  |
|  |  | Card   |                                                 |  |
|  |  +--------+                                                 |  |
|  +------------------------------------------------------------+  |
|                                                                  |
|  +------------------------------------------------------------+  |
|  |  MAIN DECK (60)                                             |  |
|  |                                                             |  |
|  |  Creatures (24)           Instants (8)        Lands (24)    |  |
|  |  4x Card Name             4x Card Name        4x Card Name  |  |
|  |  4x Card Name             4x Card Name        4x Card Name  |  |
|  |  ...                      ...                 ...           |  |
|  +------------------------------------------------------------+  |
|                                                                  |
|  +------------------------------------------------------------+  |
|  |  SIDEBOARD (15)                                             |  |
|  |  4x Card Name   4x Card Name   4x Card Name   3x Card Name  |  |
|  +------------------------------------------------------------+  |
|                                                                  |
+------------------------------------------------------------------+
```

### 4.3 Key Components

#### DeckCard.svelte (Reworked)
- Large format badge
- Color identity indicators (mana symbols)
- Mini mana curve visualization
- Card count
- Last modified date
- Action buttons on hover

#### ManaCurve.svelte
- Horizontal bar chart
- 0-7+ CMC buckets
- Color-coded by mana type
- Responsive sizing

#### DeckStats.svelte
- Format badge
- Card count breakdown
- Color identity display
- Average CMC

#### DeckImport.svelte
- Multi-format support (MTGO, Arena, text)
- Paste area with syntax highlighting
- Validation feedback
- Format auto-detection

---

## 5. Table View Rework

### 5.1 Layout Structure

```
+------------------------------------------------------------------+
|  [Logo]           MAGE            [Connection] [User Menu]       |
+------------------------------------------------------------------+
|                                                                  |
|  +------------------------------------------------------------+  |
|  |                      TABLE INFO                             |  |
|  |  [Format Badge]  Table Name                    [Leave]      |  |
|  |  Host: Player1   Password Protected                         |  |
|  +------------------------------------------------------------+  |
|                                                                  |
|  +------------------------------------------------------------+  |
|  |                     PLAYER SEATS                            |  |
|  |                                                             |  |
|  |  +---------------+              +---------------+           |  |
|  |  |   SEAT 1      |              |   SEAT 2      |           |  |
|  |  |               |              |               |           |  |
|  |  |  [Avatar]     |              |  [Empty]      |           |  |
|  |  |  Player1      |              |  Waiting...   |           |  |
|  |  |  [READY]      |              |               |           |  |
|  |  |  Deck: Modern |              |               |           |  |
|  |  |  Control      |              |               |           |  |
|  |  +---------------+              +---------------+           |  |
|  |                                                             |  |
|  |  +---------------+              +---------------+           |  |
|  |  |   SEAT 3      |              |   SEAT 4      |           |  |
|  |  |   [Empty]     |              |   [Empty]     |           |  |
|  |  +---------------+              +---------------+           |  |
|  |                                                             |  |
|  +------------------------------------------------------------+  |
|                                                                  |
|  +---------------------------+  +-----------------------------+  |
|  |      HOST CONTROLS        |  |       TABLE CHAT            |  |
|  |  (if host)                |  |                             |  |
|  |                           |  |  [Messages]                 |  |
|  |  [Kick Player v]          |  |                             |  |
|  |  [START GAME]             |  |  [Input]                    |  |
|  +---------------------------+  +-----------------------------+  |
|                                                                  |
+------------------------------------------------------------------+
```

### 5.2 Key Components

#### SeatLayout.svelte
- Grid layout adapting to player count (2/4/6/8)
- Visual connection between seats
- Centered table representation

#### PlayerSeat.svelte
- Empty state: dashed border, "Waiting for player"
- Filled state: Player avatar, name, ready status
- Ready indicator: checkmark badge
- Deck info: format/archetype if visible
- Host crown indicator
- Kick button (host only)

#### TableInfo.svelte
- Format badge
- Table name
- Host indicator
- Password indicator
- Leave button

#### ReadyIndicator.svelte
- Animated ready checkmark
- Pulsing not-ready state
- Color-coded status

---

## 6. Global CSS Variables

**File:** `src/app.css` or `src/lib/styles/variables.css`

This file should be imported in the root layout and contain all CSS custom properties defined in section 1.

---

## 7. Implementation Priority

### Phase 1: Foundation
1. Create CSS variables file with design tokens
2. Implement Button, Input, Badge, Panel components
3. Add ManaSymbol, FormatBadge MTG components

### Phase 2: Lobby Rework
4. Refactor TableCard with new design
5. Create TableFilters component
6. Refactor PlayerList
7. Update lobby page layout

### Phase 3: Deck Builder Rework
8. Refactor DeckCard with new design
9. Create ManaCurve component
10. Create DeckStats component
11. Create DeckImport modal
12. Update deck list and detail pages

### Phase 4: Table View Rework
13. Create SeatLayout component
14. Create PlayerSeat component
15. Create TableInfo component
16. Unify ChatPanel across views
17. Update table page layout

### Phase 5: Polish
18. Add transitions and animations
19. Implement skeleton loading states
20. Add tooltips throughout
21. Responsive design audit
22. Accessibility audit

---

## 8. File Changes Summary

### New Files
- `src/lib/styles/variables.css`
- `src/lib/components/ui/Button.svelte`
- `src/lib/components/ui/Input.svelte`
- `src/lib/components/ui/Select.svelte`
- `src/lib/components/ui/Badge.svelte`
- `src/lib/components/ui/Panel.svelte`
- `src/lib/components/ui/Tabs.svelte`
- `src/lib/components/ui/Skeleton.svelte`
- `src/lib/components/mtg/ManaSymbol.svelte`
- `src/lib/components/mtg/ManaCost.svelte`
- `src/lib/components/mtg/FormatBadge.svelte`
- `src/lib/components/lobby/TableFilters.svelte`
- `src/lib/components/lobby/LobbyStats.svelte`
- `src/lib/components/deck/ManaCurve.svelte`
- `src/lib/components/deck/DeckStats.svelte`
- `src/lib/components/deck/DeckImport.svelte`
- `src/lib/components/table/SeatLayout.svelte`
- `src/lib/components/table/PlayerSeat.svelte`
- `src/lib/components/table/TableInfo.svelte`
- `src/lib/components/chat/ChatPanel.svelte`
- `src/lib/components/chat/ChatMessage.svelte`
- `src/lib/components/chat/ChatInput.svelte`

### Refactored Files
- `src/lib/components/TableCard.svelte` -> `src/lib/components/lobby/TableCard.svelte`
- `src/lib/components/DeckCard.svelte` -> `src/lib/components/deck/DeckCard.svelte`
- `src/lib/components/OnlinePlayersList.svelte` -> `src/lib/components/lobby/PlayerList.svelte`
- `src/lib/components/Modal.svelte` -> integrate design tokens
- `src/lib/components/Toast.svelte` -> integrate design tokens
- `src/routes/(protected)/lobby/+page.svelte`
- `src/routes/(protected)/decks/+page.svelte`
- `src/routes/(protected)/decks/[id]/+page.svelte`
- `src/routes/(protected)/table/[id]/+page.svelte`
