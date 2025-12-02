<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { getScryfallImageUrl, getScryfallVersionForSize } from '$lib/utils/scryfall';
	import ManaSymbol from '$lib/components/mtg/ManaSymbol.svelte';

	// Props
	let {
		// eslint-disable-next-line no-unused-vars
		cardId,
		cardName,
		manaCost = '',
		cardType = '',
		power = '',
		toughness = '',
		imageUrl = '',
		isTapped = false,
		isSelected = false,
		counters = [],
		isPlaceholder = false,
		isCardBack = false,
		onclick = () => {},
		onhover = () => {},
		size = 'normal',
		// Targeting mode props
		isValidTarget = false,
		isTargetSelected = false,
		isTargetingActive = false,
		// Drag and play animation props
		isDragging = false,
		isBeingPlayed = false,
		isPending = false,
		// Ability indicator
		hasActivatedAbilities = false
	}: {
		cardId?: string;
		cardName: string;
		manaCost?: string;
		cardType?: string;
		power?: string;
		toughness?: string;
		imageUrl?: string;
		isTapped?: boolean;
		isSelected?: boolean;
		counters?: Array<{ type: string; count: number }>;
		isPlaceholder?: boolean;
		isCardBack?: boolean;
		onclick?: () => void;
		onhover?: () => void;
		size?: 'small' | 'normal' | 'large';
		// Targeting mode props
		isValidTarget?: boolean;
		isTargetSelected?: boolean;
		isTargetingActive?: boolean;
		// Drag and play animation props
		isDragging?: boolean;
		isBeingPlayed?: boolean;
		isPending?: boolean;
		// Ability indicator
		hasActivatedAbilities?: boolean;
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

	// Derive the effective image URL - use Scryfall if no explicit imageUrl provided
	const effectiveImageUrl = $derived(
		imageUrl || (!isCardBack && !isPlaceholder && cardName
			? getScryfallImageUrl(cardName, getScryfallVersionForSize(size))
			: '')
	);

	// Larger image URL for the hover preview
	const previewImageUrl = $derived(
		imageUrl || (!isCardBack && !isPlaceholder && cardName
			? getScryfallImageUrl(cardName, 'large')
			: '')
	);

	// State
	let showPreview = $state(false);
	let previewPosition = $state({ x: 0, y: 0 });
	let cardElement: HTMLDivElement | null = $state(null);
	let imageLoaded = $state(false);
	let imageError = $state(false);
	let hoverTimeout: ReturnType<typeof setTimeout> | null = null;

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
	function handleMouseEnter(event: MouseEvent): void {
		if (isCardBack || isPlaceholder) return;

		// Call hover callback
		onhover();

		// Delay preview to avoid flickering
		hoverTimeout = setTimeout(() => {
			showPreview = true;
			updatePreviewPosition(event);
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
	 * Handle mouse move - update preview position
	 */
	function handleMouseMove(event: MouseEvent): void {
		if (showPreview) {
			updatePreviewPosition(event);
		}
	}

	/**
	 * Update preview position to avoid going off-screen
	 */
	function updatePreviewPosition(event: MouseEvent): void {
		const padding = 20;
		const previewWidth = 250;
		const previewHeight = 350;

		let x = event.clientX + padding;
		let y = event.clientY + padding;

		// Check if preview goes off right edge
		if (x + previewWidth > window.innerWidth) {
			x = event.clientX - previewWidth - padding;
		}

		// Check if preview goes off bottom edge
		if (y + previewHeight > window.innerHeight) {
			y = window.innerHeight - previewHeight - padding;
		}

		// Ensure preview doesn't go off top or left
		x = Math.max(padding, x);
		y = Math.max(padding, y);

		previewPosition = { x, y };
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
		
		if (showPreview && !isCardBack && !isPlaceholder && previewImageUrl) {
			// Create/update preview element
			portalContainer.innerHTML = `
				<div class="card-preview" style="left: ${previewPosition.x}px; top: ${previewPosition.y}px;">
					<img src="${previewImageUrl}" alt="${cardName} (preview)" class="preview-image" />
				</div>
			`;
		} else {
			portalContainer.innerHTML = '';
		}
	});

	// Derived values
	const sizeClasses = $derived(() => {
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
	const hasPowerToughness = $derived(power && toughness);
</script>

<div
	bind:this={cardElement}
	class="card {sizeClasses()} {isTapped ? 'tapped' : ''} {isSelected ? 'selected' : ''} {isCardBack
		? 'card-back'
		: ''} {isTargetingActive ? 'targeting-mode' : ''} {isValidTarget ? 'valid-target' : ''} {isTargetSelected ? 'target-selected' : ''} {isTargetingActive && !isValidTarget ? 'invalid-target' : ''} {isAnimatingTap ? 'tap-animating' : ''} {isDragging ? 'dragging' : ''} {isBeingPlayed ? 'being-played' : ''} {isPending ? 'pending' : ''}"
	role="button"
	tabindex="0"
	draggable="false"
	onclick={handleClick}
	onmouseenter={handleMouseEnter}
	onmouseleave={handleMouseLeave}
	onmousemove={handleMouseMove}
	onkeydown={(e) => e.key === 'Enter' && handleClick()}
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

		<!-- Counters Badge -->
		{#if hasCounters}
			<div class="counters-badge">
				{#each counters as counter}
					<div class="counter" title="{counter.count} {counter.type} counter(s)">
						{counter.count > 0 ? `+${counter.count}` : counter.count}
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

		<!-- Pending Play Overlay -->
		{#if isPending}
			<div class="pending-overlay">
				<div class="pending-spinner"></div>
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
		width: 80px;
		height: 112px;
		font-size: 0.625rem;
	}

	.card-normal {
		width: 100px;
		height: 140px;
		font-size: 0.75rem;
	}

	.card-large {
		width: 120px;
		height: 168px;
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
	}

	.card.tap-animating.tapped {
		animation: tap-rotate-in 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
	}

	.card.tap-animating:not(.tapped) {
		animation: tap-rotate-out 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
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
		0%, 100% {
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

	/* Targeting Mode States */
	.card.targeting-mode {
		cursor: crosshair;
	}

	.card.valid-target {
		cursor: pointer;
		border-color: #fbbf24;
		box-shadow:
			0 0 0 3px rgba(251, 191, 36, 0.4),
			0 0 20px rgba(251, 191, 36, 0.6),
			0 0 40px rgba(251, 191, 36, 0.3);
		animation: target-pulse 1.5s ease-in-out infinite;
	}

	@keyframes target-pulse {
		0%, 100% {
			box-shadow:
				0 0 0 3px rgba(251, 191, 36, 0.4),
				0 0 20px rgba(251, 191, 36, 0.6),
				0 0 40px rgba(251, 191, 36, 0.3);
		}
		50% {
			box-shadow:
				0 0 0 4px rgba(251, 191, 36, 0.6),
				0 0 30px rgba(251, 191, 36, 0.8),
				0 0 50px rgba(251, 191, 36, 0.5);
		}
	}

	.card.valid-target:hover {
		transform: translateY(-20px) scale(1.1);
		border-color: #f59e0b;
		box-shadow:
			0 0 0 4px rgba(245, 158, 11, 0.6),
			0 0 35px rgba(245, 158, 11, 0.8),
			0 8px 16px rgba(0, 0, 0, 0.3);
	}

	.card.valid-target.tapped:hover {
		transform: rotate(90deg) translateY(-20px) scale(1.1);
	}

	.card.invalid-target {
		opacity: 0.4;
		filter: grayscale(60%);
		cursor: not-allowed;
		pointer-events: none;
	}

	.card.invalid-target:hover {
		transform: none;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
		border-color: #3a4451;
	}

	.card.target-selected {
		border-color: #22c55e;
		box-shadow:
			0 0 0 4px rgba(34, 197, 94, 0.5),
			0 0 25px rgba(34, 197, 94, 0.7),
			inset 0 0 30px rgba(34, 197, 94, 0.15);
		animation: target-selected-glow 1s ease-in-out infinite;
	}

	@keyframes target-selected-glow {
		0%, 100% {
			box-shadow:
				0 0 0 4px rgba(34, 197, 94, 0.5),
				0 0 25px rgba(34, 197, 94, 0.7),
				inset 0 0 30px rgba(34, 197, 94, 0.15);
		}
		50% {
			box-shadow:
				0 0 0 5px rgba(34, 197, 94, 0.7),
				0 0 35px rgba(34, 197, 94, 0.9),
				inset 0 0 40px rgba(34, 197, 94, 0.2);
		}
	}

	.card.target-selected::after {
		content: '✓';
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		font-size: 2rem;
		color: #22c55e;
		text-shadow: 0 0 10px rgba(34, 197, 94, 0.8);
		z-index: 10;
		pointer-events: none;
	}

	.card.target-selected:hover {
		transform: translateY(-15px) scale(1.05);
	}

	.card.target-selected.tapped:hover {
		transform: rotate(90deg) translateY(-15px) scale(1.05);
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
		z-index: 2;
	}

	.counter {
		padding: 0.125rem 0.375rem;
		background: #10b981;
		border-radius: 4px;
		font-size: 0.625rem;
		font-weight: 700;
		color: white;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.5);
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
		width: 250px;
		height: 350px;
		z-index: 99999;
		pointer-events: none;
		animation: cardPreviewFadeIn 0.2s ease-out;
		border-radius: 12px;
		box-shadow:
			0 20px 25px -5px rgba(0, 0, 0, 0.5),
			0 10px 10px -5px rgba(0, 0, 0, 0.3);
		border: 3px solid #667eea;
		overflow: hidden;
		background: #1a1f2e;
	}

	@keyframes cardPreviewFadeIn {
		from {
			opacity: 0;
			transform: scale(0.95);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}

	:global(.card-preview .preview-image) {
		width: 100%;
		height: 100%;
		object-fit: cover;
		border-radius: 9px;
	}
</style>
