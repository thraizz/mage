<!--
  BlockArrows.svelte
  
  SVG overlay that draws curved arrows from blocking creatures to the attackers they're blocking.
  Used during combat to visualize block assignments.
  
  Features:
  - Curved arrows with animated drawing
  - Color-coded by blocker
  - Responsive to card positions
-->
<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	// Props
	let {
		blockAssignments,
		blockerElements,
		attackerElements
	}: {
		// Map of blockerId -> attackerId
		blockAssignments: Map<string, string>;
		// Map of cardId -> HTMLElement for blockers
		blockerElements: Map<string, HTMLElement>;
		// Map of cardId -> HTMLElement for attackers
		attackerElements: Map<string, HTMLElement>;
	} = $props();

	// State
	let containerRef: SVGSVGElement | null = $state(null);
	let paths = $state<
		Array<{
			id: string;
			d: string;
			color: string;
			animationDelay: number;
		}>
	>([]);

	// Color palette for different blockers
	const colors = [
		'#3b82f6', // blue
		'#8b5cf6', // purple
		'#10b981', // emerald
		'#f59e0b', // amber
		'#ec4899', // pink
		'#06b6d4', // cyan
		'#84cc16', // lime
		'#f97316' // orange
	];

	/**
	 * Calculate SVG path from blocker to attacker
	 */
	function calculatePath(blockerEl: HTMLElement, attackerEl: HTMLElement): string {
		const blockerRect = blockerEl.getBoundingClientRect();
		const attackerRect = attackerEl.getBoundingClientRect();

		// Get center points
		const startX = blockerRect.left + blockerRect.width / 2;
		const startY = blockerRect.top;
		const endX = attackerRect.left + attackerRect.width / 2;
		const endY = attackerRect.bottom;

		// Calculate control point for quadratic bezier curve
		const midX = (startX + endX) / 2;
		const midY = (startY + endY) / 2;
		const curveOffset = Math.abs(endX - startX) * 0.3;

		// Create curved path
		const controlX = midX;
		const controlY = midY - curveOffset;

		return `M ${startX} ${startY} Q ${controlX} ${controlY} ${endX} ${endY}`;
	}

	/**
	 * Update all paths based on current assignments and element positions
	 */
	function updatePaths() {
		const newPaths: typeof paths = [];
		let colorIndex = 0;

		for (const [blockerId, attackerId] of blockAssignments) {
			const blockerEl = blockerElements.get(blockerId);
			const attackerEl = attackerElements.get(attackerId);

			if (!blockerEl || !attackerEl) continue;

			const pathData = calculatePath(blockerEl, attackerEl);
			newPaths.push({
				id: `${blockerId}-${attackerId}`,
				d: pathData,
				color: colors[colorIndex % colors.length],
				animationDelay: colorIndex * 0.1
			});

			colorIndex++;
		}

		paths = newPaths;
	}

	// Update paths when assignments or elements change
	$effect(() => {
		// Dependencies
		blockAssignments;
		blockerElements;
		attackerElements;

		// Small delay to ensure elements are rendered
		const timeout = setTimeout(updatePaths, 50);
		return () => clearTimeout(timeout);
	});

	// Update paths on window resize
	let resizeObserver: ResizeObserver | null = null;

	onMount(() => {
		// Update on resize
		resizeObserver = new ResizeObserver(updatePaths);
		if (containerRef?.parentElement) {
			resizeObserver.observe(containerRef.parentElement);
		}

		// Initial update
		updatePaths();

		// Update on scroll
		window.addEventListener('scroll', updatePaths, true);
	});

	onDestroy(() => {
		resizeObserver?.disconnect();
		window.removeEventListener('scroll', updatePaths, true);
	});
</script>

<svg bind:this={containerRef} class="block-arrows-overlay" aria-hidden="true">
	<defs>
		<!-- Arrow marker for path ends -->
		{#each paths as path}
			<marker
				id="arrowhead-{path.id}"
				markerWidth="10"
				markerHeight="7"
				refX="9"
				refY="3.5"
				orient="auto"
			>
				<polygon points="0 0, 10 3.5, 0 7" fill={path.color} />
			</marker>
		{/each}

		<!-- Glow filter -->
		<filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
			<feGaussianBlur stdDeviation="3" result="coloredBlur" />
			<feMerge>
				<feMergeNode in="coloredBlur" />
				<feMergeNode in="SourceGraphic" />
			</feMerge>
		</filter>
	</defs>

	<!-- Background glow paths -->
	{#each paths as path}
		<path
			class="arrow-glow"
			d={path.d}
			stroke={path.color}
			stroke-width="6"
			fill="none"
			stroke-linecap="round"
			filter="url(#glow)"
			style="animation-delay: {path.animationDelay}s"
		/>
	{/each}

	<!-- Main arrow paths -->
	{#each paths as path}
		<path
			class="arrow-path"
			d={path.d}
			stroke={path.color}
			stroke-width="3"
			fill="none"
			stroke-linecap="round"
			marker-end="url(#arrowhead-{path.id})"
			style="animation-delay: {path.animationDelay}s"
		/>
	{/each}
</svg>

<style>
	.block-arrows-overlay {
		position: fixed;
		inset: 0;
		pointer-events: none;
		z-index: 90;
		overflow: visible;
	}

	.arrow-path {
		stroke-dasharray: 1000;
		stroke-dashoffset: 1000;
		animation: draw-arrow 0.5s ease-out forwards;
		opacity: 0.9;
	}

	.arrow-glow {
		stroke-dasharray: 1000;
		stroke-dashoffset: 1000;
		animation: draw-arrow 0.5s ease-out forwards;
		opacity: 0.3;
	}

	@keyframes draw-arrow {
		to {
			stroke-dashoffset: 0;
		}
	}

	/* Hover effect - controlled by parent */
	:global(.block-arrows-overlay:hover .arrow-path) {
		opacity: 1;
	}

	:global(.block-arrows-overlay:hover .arrow-glow) {
		opacity: 0.5;
	}
</style>
