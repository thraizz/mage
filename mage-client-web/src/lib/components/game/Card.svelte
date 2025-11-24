<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	// Props
	let {
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
		size = 'normal'
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
	} = $props();

	// State
	let showPreview = $state(false);
	let previewPosition = $state({ x: 0, y: 0 });
	let cardElement: HTMLDivElement | null = $state(null);
	let imageLoaded = $state(false);
	let imageError = $state(false);
	let hoverTimeout: ReturnType<typeof setTimeout> | null = null;

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
	class="card {sizeClasses()} {isTapped ? 'tapped' : ''} {isSelected ? 'selected' : ''} {isCardBack ? 'card-back' : ''}"
	role="button"
	tabindex="0"
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
						<span class="mana-symbol">{symbol}</span>
					{/each}
				</div>
			{/if}
		</div>
	{:else}
		<!-- Real Card -->
		<div class="card-inner">
			{#if imageUrl && !imageError}
				<img
					src={imageUrl}
					alt={cardName}
					class="card-image {imageLoaded ? 'loaded' : ''}"
					onload={handleImageLoad}
					onerror={handleImageError}
				/>
			{:else}
				<!-- Fallback when no image -->
				<div class="card-fallback">
					<div class="card-name-text">{cardName}</div>
					{#if manaCost}
						<div class="mana-cost">
							{#each manaSymbols as symbol}
								<span class="mana-symbol mana-{symbol.toLowerCase()}">{symbol}</span>
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

			<!-- Loading Spinner -->
			{#if imageUrl && !imageLoaded && !imageError}
				<div class="loading-spinner"></div>
			{/if}
		</div>
	{/if}
</div>

<!-- Hover Preview Portal -->
{#if showPreview && !isCardBack && !isPlaceholder}
	<div
		class="card-preview"
		style="left: {previewPosition.x}px; top: {previewPosition.y}px;"
		role="tooltip"
	>
		{#if imageUrl && !imageError}
			<img src={imageUrl} alt="{cardName} (preview)" class="preview-image" />
		{:else}
			<div class="preview-fallback">
				<div class="preview-name">{cardName}</div>
				{#if manaCost}
					<div class="preview-mana">
						{#each manaSymbols as symbol}
							<span class="mana-symbol mana-{symbol.toLowerCase()}">{symbol}</span>
						{/each}
					</div>
				{/if}
				{#if cardType}
					<div class="preview-type">{cardType}</div>
				{/if}
				{#if hasPowerToughness}
					<div class="preview-pt">{power}/{toughness}</div>
				{/if}
			</div>
		{/if}
	</div>
{/if}

<style>
	/* Base Card Styles */
	.card {
		position: relative;
		border-radius: 8px;
		cursor: pointer;
		transition:
			transform 0.2s,
			box-shadow 0.2s,
			filter 0.2s;
		user-select: none;
		background: #1a1f2e;
		border: 2px solid #3a4451;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
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

	.card.tapped {
		transform: rotate(90deg);
		transform-origin: center center;
	}

	.card.tapped:hover {
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
		background: repeating-linear-gradient(
			45deg,
			#3a4451,
			#3a4451 10px,
			#2a3441 10px,
			#2a3441 20px
		);
		border-radius: 4px;
		opacity: 0.5;
	}

	/* Card Inner */
	.card-inner,
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
		width: 100%;
		height: 100%;
		object-fit: cover;
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

	/* Card Preview (Hover) */
	.card-preview {
		position: fixed;
		width: 250px;
		height: 350px;
		z-index: 9999;
		pointer-events: none;
		animation: fadeIn 0.2s;
		border-radius: 12px;
		box-shadow:
			0 20px 25px -5px rgba(0, 0, 0, 0.5),
			0 10px 10px -5px rgba(0, 0, 0, 0.3);
		border: 3px solid #667eea;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: scale(0.95);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}

	.preview-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
		border-radius: 10px;
	}

	.preview-fallback {
		width: 100%;
		height: 100%;
		padding: 1.5rem;
		background: #1a1f2e;
		border-radius: 10px;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		color: white;
	}

	.preview-name {
		font-size: 1.25rem;
		font-weight: 700;
		line-height: 1.3;
	}

	.preview-mana {
		display: flex;
		gap: 0.25rem;
		flex-wrap: wrap;
	}

	.preview-mana .mana-symbol {
		width: 1.5rem;
		height: 1.5rem;
		font-size: 0.875rem;
	}

	.preview-type {
		font-size: 0.875rem;
		color: #9ca3af;
		margin-top: 1rem;
	}

	.preview-pt {
		font-size: 1.5rem;
		font-weight: 700;
		color: #fbbf24;
		margin-top: auto;
		text-align: center;
	}
</style>
