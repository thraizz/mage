<script lang="ts">
  import ManaSymbol from '$lib/components/mtg/ManaSymbol.svelte';
  import {
    getScryfallImageUrl,
    getScryfallTokenSearchUrl,
    getScryfallVersionForSize
  } from '$lib/utils/scryfall';
  import { onDestroy, onMount } from 'svelte';

  // Props
  let {
    cardId,
    cardName,
    manaCost = '',
    cardType = '',
    power = '',
    toughness = '',
    color = '',
    imageUrl = '',
    isTapped = false,
    isSelected = false,
    counters = [],
    isPlaceholder = false,
    isCardBack = false,
    onclick = () => {},
    onhover = () => {},
    oncontextmenu = undefined as ((e: MouseEvent) => void) | undefined,
    size = 'normal',
    // Drag and play animation props
    isDragging = false,
    isBeingPlayed = false,
    isPending = false,
    // Ability indicator
    hasActivatedAbilities = false,
    // Combat mode props
    isAttacking = false,
    canAttack = false,
    attackTarget = '',
    isBlocking = false,
    canBlock = false,
    blockingWhat = [] as string[],
    // Summoning sickness
    summoningSickness = false
  }: {
    cardId?: string;
    cardName: string;
    manaCost?: string;
    cardType?: string;
    power?: string;
    toughness?: string;
    color?: string;
    imageUrl?: string;
    isTapped?: boolean;
    isSelected?: boolean;
    counters?: Array<{ name: string; count: number }>;
    isPlaceholder?: boolean;
    isCardBack?: boolean;
    onclick?: () => void;
    onhover?: () => void;
    oncontextmenu?: (e: MouseEvent) => void;
    size?: 'small' | 'normal' | 'large';
    // Drag and play animation props
    isDragging?: boolean;
    isBeingPlayed?: boolean;
    isPending?: boolean;
    // Ability indicator
    hasActivatedAbilities?: boolean;
    // Combat mode props
    isAttacking?: boolean;
    canAttack?: boolean;
    attackTarget?: string;
    isBlocking?: boolean;
    canBlock?: boolean;
    blockingWhat?: string[];
    // Summoning sickness
    summoningSickness?: boolean;
  } = $props();

  // Track tap state changes for animation
  let prevTapped = $state(isTapped);
  let isAnimatingTap = $state(false);

  // Detect tap state changes and trigger animation
  $effect(() => {
    if (isTapped !== prevTapped) {
      isAnimatingTap = true;
      prevTapped = isTapped;
      // Reset animation state after animation completes
      const timeout = setTimeout(() => {
        isAnimatingTap = false;
      }, 200); // Match animation duration
      return () => clearTimeout(timeout);
    }
  });

  // State
  let showPreview = $state(false);
  let cardElement: HTMLDivElement | null = $state(null);
  let imageLoaded = $state(false);
  let imageError = $state(false);
  let hoverTimeout: ReturnType<typeof setTimeout> | null = null;
  let tokenImageUrl = $state<string | null>(null);
  let tokenImageFetching = $state(false);
  let lastFetchedTokenKey = $state<string>('');

  // Check if this is a token based on card ID
  const isToken = $derived(cardId?.startsWith('token-') ?? false);

  // Create a unique key for this token based on its characteristics
  const tokenKey = $derived(isToken ? `${cardName}|${power}|${toughness}|${color}` : '');

  // Fetch token image from Scryfall search API
  async function fetchTokenImage(): Promise<void> {
    if (!isToken || !cardName || tokenImageFetching) return;

    // Don't refetch if we already have this token
    if (lastFetchedTokenKey === tokenKey && tokenImageUrl) return;

    tokenImageFetching = true;
    lastFetchedTokenKey = tokenKey;

    try {
      const searchUrl = getScryfallTokenSearchUrl(cardName, power, toughness, color);
      const response = await fetch(searchUrl);

      if (!response.ok) {
        console.warn(`Failed to fetch token image for ${cardName}:`, response.statusText);
        tokenImageUrl = null;
        tokenImageFetching = false;
        return;
      }

      const data = await response.json();

      // Get the first matching card's image
      if (data.data && data.data.length > 0) {
        const card = data.data[0];
        const versionKey = getScryfallVersionForSize(size);
        // Extract image URL from the card object
        tokenImageUrl = card.image_uris?.[versionKey] || card.image_uris?.normal || null;
      } else {
        console.warn(
          `No token found for ${cardName} with power=${power} toughness=${toughness} color=${color}`
        );
        tokenImageUrl = null;
      }
    } catch (error) {
      console.error(`Error fetching token image for ${cardName}:`, error);
      tokenImageUrl = null;
    } finally {
      tokenImageFetching = false;
    }
  }

  // Fetch token image when token identity changes
  $effect(() => {
    if (isToken && cardName && !imageUrl && tokenKey !== lastFetchedTokenKey) {
      fetchTokenImage();
    }
  });

  // Derive the effective image URL - use Scryfall if no explicit imageUrl provided
  const effectiveImageUrl = $derived(
    imageUrl ||
      (isToken && tokenImageUrl
        ? tokenImageUrl
        : !isCardBack && !isPlaceholder && cardName
          ? getScryfallImageUrl(cardName, getScryfallVersionForSize(size))
          : '')
  );

  // Create portal container for preview (to avoid transform issues with fixed positioning)
  let portalContainer: HTMLDivElement | null = null;

  onMount(() => {
    // Create a container at the body level for the preview portal
    portalContainer = document.createElement('div');
    portalContainer.className = 'card-preview-portal';
    document.body.appendChild(portalContainer);
  });

  /**
   * Handle mouse enter - show preview after delay
   */
  function handleMouseEnter(): void {
    if (isCardBack || isPlaceholder) return;

    // Call hover callback
    onhover();

    // Delay preview to avoid flickering
    hoverTimeout = setTimeout(() => {
      showPreview = true;
    }, 300);
  }

  /**
   * Handle mouse leave - hide preview
   */
  function handleMouseLeave(): void {
    if (hoverTimeout) {
      clearTimeout(hoverTimeout);
      hoverTimeout = null;
    }
    showPreview = false;
  }

  /**
   * Handle card click
   */
  function handleClick(): void {
    onclick();
  }

  /**
   * Handle image load
   */
  function handleImageLoad(): void {
    imageLoaded = true;
    imageError = false;
  }

  /**
   * Handle image error
   */
  function handleImageError(): void {
    imageLoaded = false;
    imageError = true;
  }

  /**
   * Parse mana cost to display mana symbols
   */
  function parseManaCost(cost: string): string[] {
    if (!cost) return [];
    // Simple regex to extract mana symbols: {W}, {U}, {B}, {R}, {G}, {1}, {2}, etc.
    const symbols = cost.match(/\{[^}]+\}/g) || [];
    return symbols.map((s) => s.replace(/[{}]/g, ''));
  }

  // Cleanup on destroy
  onDestroy(() => {
    if (hoverTimeout) {
      clearTimeout(hoverTimeout);
    }
    // Remove portal container
    if (portalContainer && portalContainer.parentNode) {
      portalContainer.parentNode.removeChild(portalContainer);
    }
  });

  // Update preview element when showPreview changes
  $effect(() => {
    if (!portalContainer) return;

    if (showPreview && !isCardBack && !isPlaceholder && effectiveImageUrl) {
      // Create/update preview element - centered in viewport
      portalContainer.innerHTML = `
				<div class="card-preview">
					<img src="${effectiveImageUrl}" alt="${cardName} (preview)" class="preview-image" />
				</div>
			`;
    } else {
      portalContainer.innerHTML = '';
    }
  });

  // Derived values
  const sizeClasses = $derived.by(() => {
    switch (size) {
      case 'small':
        return 'card-small';
      case 'large':
        return 'card-large';
      default:
        return 'card-normal';
    }
  });

  const manaSymbols = $derived(parseManaCost(manaCost));
  const hasCounters = $derived(counters && counters.length > 0);
  const hasPowerToughness = $derived(!!power && !!toughness && power !== '' && toughness !== '');

  // Debug the power and toughness values whenever they change
  $effect(() => {
    console.log(
      power !== '' || toughness !== ''
        ? `Card ${cardName} has P/T ${power}/${toughness}`
        : `Card ${cardName} has no P/T, its a ${cardType}`
    );
  });
</script>

<div
  bind:this={cardElement}
  class="card {sizeClasses} {isTapped ? 'tapped' : ''} {isSelected ? 'selected' : ''} {isCardBack
    ? 'card-back'
    : ''} {isAnimatingTap ? 'tap-animating' : ''} {isDragging ? 'dragging' : ''} {isBeingPlayed
    ? 'being-played'
    : ''} {isPending ? 'pending' : ''} {isAttacking ? 'attacking' : ''} {canAttack
    ? 'can-attack'
    : ''} {isBlocking ? 'blocking' : ''} {canBlock ? 'can-block' : ''} {summoningSickness
    ? 'summoning-sick'
    : ''}"
  role="button"
  tabindex="0"
  draggable="false"
  onclick={handleClick}
  onmouseenter={handleMouseEnter}
  onmouseleave={handleMouseLeave}
  onkeydown={(e) => e.key === 'Enter' && handleClick()}
  oncontextmenu={(e) => oncontextmenu?.(e)}
>
  {#if isCardBack}
    <!-- Card Back -->
    <div class="card-back-inner">
      <div class="card-back-pattern"></div>
    </div>
  {:else if isPlaceholder}
    <!-- Placeholder Card -->
    <div class="card-placeholder-inner">
      <div class="card-name-placeholder">{cardName}</div>
      {#if manaCost}
        <div class="mana-cost-placeholder">
          {#each manaSymbols as symbol}
            <ManaSymbol {symbol} size="sm" />
          {/each}
        </div>
      {/if}
    </div>
  {:else}
    <!-- Real Card -->
    {#if effectiveImageUrl && !imageError}
      <img
        src={effectiveImageUrl}
        alt={cardName}
        class="card-image {imageLoaded ? 'loaded' : ''}"
        draggable="false"
        onload={handleImageLoad}
        onerror={handleImageError}
      />
      <!-- Loading Spinner -->
      {#if !imageLoaded}
        <div class="loading-spinner"></div>
      {/if}
    {:else}
      <!-- Fallback when no image -->
      <div class="card-fallback">
        <div class="card-name-text">{cardName}</div>
        {#if manaCost}
          <div class="mana-cost">
            {#each manaSymbols as symbol}
              <ManaSymbol {symbol} size="sm" />
            {/each}
          </div>
        {/if}
        {#if cardType}
          <div class="card-type-text">{cardType}</div>
        {/if}
        {#if hasPowerToughness}
          <div class="power-toughness">
            {power}/{toughness}
          </div>
        {/if}
      </div>
    {/if}
    <!-- Real Card -->
    {#if effectiveImageUrl && !imageError}
      <img
        src={effectiveImageUrl}
        alt={cardName}
        class="card-image {imageLoaded ? 'loaded' : ''}"
        draggable="false"
        onload={handleImageLoad}
        onerror={handleImageError}
      />
      <!-- Loading Spinner -->
      {#if !imageLoaded}
        <div class="loading-spinner"></div>
      {/if}

      <!-- Text Overlay for Readability -->
      <div class="text-overlay">
        <!-- Card Name & Mana Cost -->
        <div class="overlay-header">
          <div class="overlay-card-name">{cardName}</div>
          {#if manaCost}
            <div class="overlay-mana-cost">
              {#each manaSymbols as symbol}
                <ManaSymbol {symbol} size="sm" />
              {/each}
            </div>
          {/if}
        </div>

        <!-- Card Type -->
        {#if cardType}
          <div class="overlay-type">{cardType}</div>
        {/if}

        <!-- Power/Toughness -->
        {#if hasPowerToughness}
          <div class="overlay-pt">
            {power}/{toughness}
          </div>
        {/if}
      </div>
    {/if}

    <!-- Counters Badge -->
    {#if hasCounters}
      <div class="counters-badge">
        {#each counters as counter}
          <div class="counter" title="{counter.count} {counter.name} counter(s)">
            <span class="counter-name">{counter.name}</span>
            <span class="counter-count"
              >{counter.count > 0 ? `x${counter.count}` : counter.count}</span
            >
          </div>
        {/each}
      </div>
    {/if}

    <!-- Activated Abilities Badge -->
    {#if hasActivatedAbilities}
      <div class="ability-badge" title="Has activated abilities (press A)">
        <span class="ability-icon">⚡</span>
      </div>
    {/if}

    <!-- Combat: Attacking Badge -->
    {#if isAttacking}
      <div
        class="combat-badge attacking-badge"
        title="Attacking{attackTarget ? ` ${attackTarget}` : ''}"
      >
        <span class="combat-icon">⚔️</span>
      </div>
      {#if attackTarget}
        <div class="attack-target-label">
          → {attackTarget}
        </div>
      {/if}
    {/if}

    <!-- Combat: Can Attack Indicator -->
    {#if canAttack && !isAttacking}
      <div class="can-attack-indicator" title="Click to declare as attacker">
        <span class="attack-hint">⚔️</span>
      </div>
    {/if}

    <!-- Combat: Blocking Badge -->
    {#if isBlocking}
      <div
        class="combat-badge blocking-badge"
        title="Blocking{blockingWhat.length > 0 ? `: ${blockingWhat.join(', ')}` : ''}"
      >
        <span class="combat-icon">🛡️</span>
      </div>
      {#if blockingWhat.length > 0}
        <div class="blocking-label">
          ↔ {blockingWhat.join(', ')}
        </div>
      {/if}
    {/if}

    <!-- Combat: Can Block Indicator -->
    {#if canBlock && !isBlocking}
      <div class="can-block-indicator" title="Click to select as blocker">
        <span class="block-hint">🛡️</span>
      </div>
    {/if}

    <!-- Pending Play Overlay -->
    {#if isPending}
      <div class="pending-overlay">
        <div class="pending-spinner"></div>
      </div>
    {/if}

    <!-- Summoning Sickness Overlay -->
    {#if summoningSickness}
      <div
        class="summoning-sickness-overlay"
        title="Summoning Sickness - Cannot attack or use tap abilities this turn"
      >
        <div class="zzz-indicator">💤</div>
      </div>
    {/if}
  {/if}
</div>

<!-- Preview is rendered via portal to document.body to avoid transform issues -->

<style>
  /* Base Card Styles */
  .card {
    position: relative;
    border-radius: 8px;
    cursor: pointer;
    /* Smooth transitions for all interactive states */
    transition:
      transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1),
      box-shadow 0.2s ease-out,
      filter 0.2s ease-out,
      border-color 0.15s ease-out;
    padding: 0;
    transform-origin: center center;
    user-select: none;
    -webkit-user-select: none;
    -webkit-user-drag: none;
    background: #1a1f2e;
    border: 2px solid #3a4451;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    /* Enable GPU acceleration for smooth animations */
    will-change: transform;
    /* Touch device optimization */
    touch-action: manipulation;
    -webkit-tap-highlight-color: transparent;
  }

  .card * {
    -webkit-user-drag: none;
    user-select: none;
  }

  /* Card Sizes */
  .card-small {
    height: calc(var(--card-height) * 0.8);
    width: calc(var(--card-width) * 0.8);
    font-size: 0.625rem;
  }

  .card-normal {
    height: var(--card-height);
    width: var(--card-width);
    font-size: 0.75rem;
  }

  .card-large {
    height: calc(var(--card-height) * 1.25);
    width: calc(var(--card-width) * 1.25);
    font-size: 0.875rem;
  }

  /* Card States */
  .card:hover {
    transform: translateY(-15px) scale(1.05);
    box-shadow: 0 8px 16px rgba(102, 126, 234, 0.3);
    border-color: #667eea;
    z-index: 10;
  }

  .card:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.3);
  }

  .card.selected {
    border-color: #fbbf24;
    box-shadow: 0 0 0 3px rgba(251, 191, 36, 0.4);
  }

  /* Tapped State - 90° clockwise rotation */
  .card.tapped {
    transform: rotate(90deg);
    transform-origin: center center;
  }

  .card.tapped:hover {
    transform: rotate(90deg) translateY(-15px) scale(1.05);
  }

  /* Tap Animation - Visual feedback during tap/untap */
  .card.tap-animating {
    animation: tap-glow 0.2s ease-out;
    /* Disable transform transition during animation to prevent jump after animation ends */
    transition:
      box-shadow 0.2s ease-out,
      filter 0.2s ease-out,
      border-color 0.15s ease-out;
  }

  .card.tap-animating.tapped {
    animation: tap-rotate-in 0.2s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
  }

  .card.tap-animating:not(.tapped) {
    animation: tap-rotate-out 0.2s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
  }

  /* Tap animation keyframes - Clockwise rotation with slight overshoot */
  @keyframes tap-rotate-in {
    0% {
      transform: rotate(0deg);
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    }
    50% {
      box-shadow:
        0 4px 12px rgba(102, 126, 234, 0.4),
        0 0 20px rgba(102, 126, 234, 0.2);
    }
    70% {
      transform: rotate(95deg);
    }
    100% {
      transform: rotate(90deg);
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    }
  }

  /* Untap animation keyframes - Counter-clockwise rotation with slight overshoot */
  @keyframes tap-rotate-out {
    0% {
      transform: rotate(90deg);
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    }
    50% {
      box-shadow:
        0 4px 12px rgba(16, 185, 129, 0.4),
        0 0 20px rgba(16, 185, 129, 0.2);
    }
    70% {
      transform: rotate(-5deg);
    }
    100% {
      transform: rotate(0deg);
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    }
  }

  /* Subtle glow pulse during any tap state change */
  @keyframes tap-glow {
    0%,
    100% {
      filter: brightness(1);
    }
    50% {
      filter: brightness(1.15);
    }
  }

  /* Drag State */
  .card.dragging {
    opacity: 0.5;
    transform: scale(0.95);
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.5);
    cursor: grabbing;
  }

  .card.dragging:hover {
    transform: scale(0.95);
  }

  /* Being Played Animation (fly out of hand) */
  .card.being-played {
    animation: card-fly-out 0.4s ease-out forwards;
    pointer-events: none;
  }

  @keyframes card-fly-out {
    0% {
      transform: translateY(0) scale(1);
      opacity: 1;
    }
    30% {
      transform: translateY(-30px) scale(1.1);
      opacity: 1;
    }
    100% {
      transform: translateY(-100px) scale(0.8);
      opacity: 0;
    }
  }

  /* Pending State (waiting for server confirmation) */
  .card.pending {
    opacity: 0.7;
    filter: brightness(0.9);
    pointer-events: none;
  }

  .pending-overlay {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 5;
  }

  .pending-spinner {
    width: 28px;
    height: 28px;
    border: 3px solid rgba(255, 255, 255, 0.3);
    border-top-color: #667eea;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  /* Card Back */
  .card-back {
    background: linear-gradient(135deg, #2a3441 0%, #1a1f2e 100%);
    border: 2px solid #3a4451;
  }

  .card-back:hover {
    transform: none;
    border-color: #3a4451;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    cursor: default;
  }

  .card-back-inner {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .card-back-pattern {
    width: 60%;
    height: 60%;
    background: repeating-linear-gradient(45deg, #3a4451, #3a4451 10px, #2a3441 10px, #2a3441 20px);
    border-radius: 4px;
    opacity: 0.5;
  }

  /* Card Placeholder Inner */
  .card-placeholder-inner {
    width: 100%;
    height: 100%;
    border-radius: 6px;
    overflow: hidden;
    position: relative;
    display: flex;
    flex-direction: column;
  }

  /* Card Image */
  .card-image {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 6px;
    opacity: 0;
    transition: opacity 0.3s;
  }

  .card-image.loaded {
    opacity: 1;
  }

  /* Card Fallback (No Image) */
  .card-fallback,
  .card-placeholder-inner {
    width: 100%;
    height: 100%;
    padding: 0.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    background: #1a1f2e;
    color: white;
  }

  .card-name-text,
  .card-name-placeholder {
    font-weight: 600;
    font-size: 0.75rem;
    line-height: 1.2;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .mana-cost,
  .mana-cost-placeholder {
    display: flex;
    gap: 0.125rem;
    flex-wrap: wrap;
  }

  .mana-symbol {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1rem;
    height: 1rem;
    border-radius: 50%;
    background: #3a4451;
    font-size: 0.625rem;
    font-weight: 700;
    color: white;
  }

  /* Mana Symbol Colors */
  .mana-w {
    background: #f0f0d8;
    color: #000;
  }
  .mana-u {
    background: #0e68ab;
  }
  .mana-b {
    background: #150b00;
  }
  .mana-r {
    background: #d3202a;
  }
  .mana-g {
    background: #00733e;
  }

  .card-type-text {
    font-size: 0.625rem;
    color: #9ca3af;
    margin-top: auto;
  }

  .power-toughness {
    font-weight: 700;
    font-size: 0.875rem;
    color: #fbbf24;
    margin-top: auto;
    text-align: right;
  }

  /* Counters Badge */
  .counters-badge {
    position: absolute;
    top: 0.25rem;
    right: 0.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
    z-index: 20;
  }

  .counter {
    padding: 0.125rem 0.375rem;
    background: #10b981;
    border-radius: 4px;
    font-size: 0.625rem;
    font-weight: 700;
    color: white;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }

  .counter-name {
    font-weight: 800;
  }

  .counter-count {
    font-weight: 700;
  }

  /* Ability Badge */
  .ability-badge {
    position: absolute;
    top: 0.25rem;
    left: 0.25rem;
    width: 1.125rem;
    height: 1.125rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(251, 191, 36, 0.85);
    border-radius: 50%;
    z-index: 2;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.4);
  }

  .ability-icon {
    font-size: 0.6875rem;
    line-height: 1;
  }

  /* Loading Spinner */
  .loading-spinner {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 24px;
    height: 24px;
    border: 3px solid #3a4451;
    border-top-color: #667eea;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: translate(-50%, -50%) rotate(360deg);
    }
  }

  /* Card Preview (Hover) - Global styles for portal-rendered preview */
  :global(.card-preview-portal) {
    position: fixed;
    top: 0;
    left: 0;
    width: 0;
    height: 0;
    overflow: visible;
    pointer-events: none;
    z-index: 99999;
  }

  :global(.card-preview) {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    height: auto;
    max-height: 90vh;
    z-index: 99999;
    pointer-events: none;
    border-radius: 39px;
    box-shadow:
      0 20px 25px -5px rgba(0, 0, 0, 0.5),
      0 10px 10px -5px rgba(0, 0, 0, 0.3);
    border: 3px solid #667eea;
    overflow: hidden;
    background: #1a1f2e;
  }

  :global(.card-preview .preview-image) {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 9px;
  }

  /* ============================================
	   COMBAT STATES
	   ============================================ */

  /* Attacking State - Red theme with elevation */
  .card.attacking {
    border-color: #ef4444;
    box-shadow:
      0 0 0 3px rgba(239, 68, 68, 0.4),
      0 8px 20px rgba(239, 68, 68, 0.5),
      0 0 40px rgba(239, 68, 68, 0.2);
    transform: translateY(-12px);
    z-index: 5;
    animation: attacking-pulse 1.5s ease-in-out infinite;
  }

  @keyframes attacking-pulse {
    0%,
    100% {
      box-shadow:
        0 0 0 3px rgba(239, 68, 68, 0.4),
        0 8px 20px rgba(239, 68, 68, 0.5),
        0 0 40px rgba(239, 68, 68, 0.2);
    }
    50% {
      box-shadow:
        0 0 0 4px rgba(239, 68, 68, 0.6),
        0 10px 25px rgba(239, 68, 68, 0.6),
        0 0 50px rgba(239, 68, 68, 0.3);
    }
  }

  .card.attacking:hover {
    transform: translateY(-18px) scale(1.05);
  }

  .card.attacking.tapped {
    transform: rotate(90deg) translateY(-12px);
  }

  .card.attacking.tapped:hover {
    transform: rotate(90deg) translateY(-18px) scale(1.05);
  }

  /* Can Attack State - Green dashed outline */
  .card.can-attack {
    border-style: dashed;
    border-color: #22c55e;
    box-shadow:
      0 0 0 2px rgba(34, 197, 94, 0.3),
      0 0 15px rgba(34, 197, 94, 0.2);
    animation: can-attack-shimmer 2s ease-in-out infinite;
  }

  @keyframes can-attack-shimmer {
    0%,
    100% {
      box-shadow:
        0 0 0 2px rgba(34, 197, 94, 0.3),
        0 0 15px rgba(34, 197, 94, 0.2);
    }
    50% {
      box-shadow:
        0 0 0 3px rgba(34, 197, 94, 0.5),
        0 0 25px rgba(34, 197, 94, 0.3);
    }
  }

  .card.can-attack:hover {
    border-color: #16a34a;
    transform: translateY(-20px) scale(1.08);
    box-shadow:
      0 0 0 3px rgba(34, 197, 94, 0.6),
      0 10px 20px rgba(34, 197, 94, 0.4);
  }

  /* Blocking State - Blue theme */
  .card.blocking {
    border-color: #3b82f6;
    box-shadow:
      0 0 0 3px rgba(59, 130, 246, 0.4),
      0 8px 20px rgba(59, 130, 246, 0.5),
      0 0 40px rgba(59, 130, 246, 0.2);
    z-index: 5;
    animation: blocking-pulse 1.5s ease-in-out infinite;
  }

  @keyframes blocking-pulse {
    0%,
    100% {
      box-shadow:
        0 0 0 3px rgba(59, 130, 246, 0.4),
        0 8px 20px rgba(59, 130, 246, 0.5),
        0 0 40px rgba(59, 130, 246, 0.2);
    }
    50% {
      box-shadow:
        0 0 0 4px rgba(59, 130, 246, 0.6),
        0 10px 25px rgba(59, 130, 246, 0.6),
        0 0 50px rgba(59, 130, 246, 0.3);
    }
  }

  .card.blocking:hover {
    transform: translateY(-15px) scale(1.05);
  }

  /* Can Block State - Blue dashed outline */
  .card.can-block {
    border-style: dashed;
    border-color: #3b82f6;
    box-shadow:
      0 0 0 2px rgba(59, 130, 246, 0.3),
      0 0 15px rgba(59, 130, 246, 0.2);
    animation: can-block-shimmer 2s ease-in-out infinite;
  }

  @keyframes can-block-shimmer {
    0%,
    100% {
      box-shadow:
        0 0 0 2px rgba(59, 130, 246, 0.3),
        0 0 15px rgba(59, 130, 246, 0.2);
    }
    50% {
      box-shadow:
        0 0 0 3px rgba(59, 130, 246, 0.5),
        0 0 25px rgba(59, 130, 246, 0.3);
    }
  }

  .card.can-block:hover {
    border-color: #2563eb;
    transform: translateY(-20px) scale(1.08);
    box-shadow:
      0 0 0 3px rgba(59, 130, 246, 0.6),
      0 10px 20px rgba(59, 130, 246, 0.4);
  }

  /* Combat Badges */
  .combat-badge {
    position: absolute;
    top: -0.5rem;
    left: 50%;
    transform: translateX(-50%);
    padding: 0.25rem 0.5rem;
    border-radius: 9999px;
    font-size: 0.875rem;
    font-weight: 700;
    z-index: 10;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
    white-space: nowrap;
  }

  .attacking-badge {
    background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
    color: white;
    border: 2px solid #fca5a5;
  }

  .blocking-badge {
    background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
    color: white;
    border: 2px solid #93c5fd;
  }

  .combat-icon {
    font-size: 0.75rem;
  }

  /* Attack Target Label */
  .attack-target-label {
    position: absolute;
    bottom: -1.5rem;
    left: 50%;
    transform: translateX(-50%);
    padding: 0.125rem 0.5rem;
    background: rgba(239, 68, 68, 0.9);
    color: white;
    font-size: 0.625rem;
    font-weight: 600;
    border-radius: 4px;
    white-space: nowrap;
    z-index: 10;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  }

  /* Blocking Label */
  .blocking-label {
    position: absolute;
    bottom: -1.5rem;
    left: 50%;
    transform: translateX(-50%);
    padding: 0.125rem 0.5rem;
    background: rgba(59, 130, 246, 0.9);
    color: white;
    font-size: 0.625rem;
    font-weight: 600;
    border-radius: 4px;
    white-space: nowrap;
    z-index: 10;
    max-width: 150px;
    overflow: hidden;
    text-overflow: ellipsis;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  }

  /* Can Attack/Block Indicators */
  .can-attack-indicator,
  .can-block-indicator {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 2.5rem;
    height: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    opacity: 0;
    transition: opacity 0.2s;
    z-index: 8;
    pointer-events: none;
  }

  .can-attack-indicator {
    background: rgba(34, 197, 94, 0.8);
    border: 2px solid rgba(255, 255, 255, 0.5);
  }

  .can-block-indicator {
    background: rgba(59, 130, 246, 0.8);
    border: 2px solid rgba(255, 255, 255, 0.5);
  }

  .card.can-attack:hover .can-attack-indicator,
  .card.can-block:hover .can-block-indicator {
    opacity: 1;
  }

  .attack-hint,
  .block-hint {
    font-size: 1.25rem;
  }

  /* ===== Summoning Sickness Styles ===== */
  .card.summoning-sick {
    /* Desaturation to indicate can't attack */
    filter: saturate(0.5) brightness(0.75);
  }

  .summoning-sickness-overlay {
    position: absolute;
    inset: 0;
    border-radius: 6px;
    pointer-events: none;
    z-index: 15;
    /* Simple semi-transparent purple overlay */
    background: rgba(88, 28, 135, 0.35);
  }

  /* ZZZ sleeping indicator - simple pulse animation */
  .zzz-indicator {
    position: absolute;
    top: 0.25rem;
    right: 0.25rem;
    font-size: 1.25rem;
    animation: zzz-pulse 2s ease-in-out infinite;
    /* Single lightweight shadow */
    filter: drop-shadow(0 0 4px rgba(255, 255, 255, 0.8));
    z-index: 25;
  }

  @keyframes zzz-pulse {
    0%,
    100% {
      opacity: 0.7;
      transform: scale(1);
    }
    50% {
      opacity: 1;
      transform: scale(1.1);
    }
  }

  /* Hide summoning sickness for card backs and placeholders */
  .card-back .summoning-sickness-overlay,
  .card-placeholder .summoning-sickness-overlay {
    display: none;
  }

  /* ============================================
   TEXT OVERLAY FOR READABILITY
   ============================================ */

  .text-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    pointer-events: none;
    z-index: 3;
    padding: 0.375rem;
    border-radius: 6px;
  }

  /* Header: Card Name + Mana Cost */
  .overlay-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 0.25rem;
    margin-bottom: auto;
  }

  .overlay-card-name {
    flex: 1;
    font-size: 0.6875rem;
    font-weight: 700;
    line-height: 1.2;
    color: white;
    background: rgba(0, 0, 0, 0.75);
    padding: 0.25rem 0.375rem;
    border-radius: 4px;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
    background: rgba(0, 0, 0, 0.85); /* Darker background instead of blur */
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    line-clamp: 2;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .overlay-mana-cost {
    display: flex;
    gap: 0.125rem;
    flex-wrap: wrap;
    background: rgba(0, 0, 0, 0.75);
    padding: 0.25rem;
    border-radius: 4px;
    background: rgba(0, 0, 0, 0.85);
  }

  /* Card Type */
  .overlay-type {
    font-size: 0.625rem;
    font-weight: 600;
    color: white;
    background: rgba(0, 0, 0, 0.7);
    padding: 0.125rem 0.375rem;
    border-radius: 4px;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
    background: rgba(0, 0, 0, 0.85);
    align-self: flex-start;
    max-width: 90%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    bottom: 36%;
    position: relative;
  }

  /* Power/Toughness */
  .overlay-pt {
    font-size: 0.875rem;
    font-weight: 700;
    color: #fbbf24;
    background: rgba(0, 0, 0, 0.8);
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    text-shadow:
      0 1px 2px rgba(0, 0, 0, 0.8),
      0 0 8px rgba(251, 191, 36, 0.5);
    background: rgba(0, 0, 0, 0.85);
    align-self: flex-end;
    min-width: 2rem;
    text-align: center;
    position: absolute;
    bottom: 0;
  }

  /* Responsive sizing for different card sizes */
  .card-small .overlay-card-name {
    font-size: 0.5rem;
    padding: 0.125rem 0.25rem;
  }

  .card-small .overlay-type {
    font-size: 0.5rem;
    padding: 0.125rem 0.25rem;
  }

  .card-small .overlay-pt {
    font-size: 0.625rem;
    padding: 0.125rem 0.25rem;
  }

  .card-large .overlay-card-name {
    font-size: 0.75rem;
  }

  .card-large .overlay-type {
    font-size: 0.6875rem;
  }

  .card-large .overlay-pt {
    font-size: 1rem;
  }
</style>
