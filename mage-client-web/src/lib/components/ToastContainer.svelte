<script lang="ts">
	import { toast } from '$lib/stores/toast';
	import Toast from './Toast.svelte';
</script>

<!--
  Toast container positioned at top-right of viewport
  Uses ARIA live region for screen reader accessibility
-->
<div class="toast-container" aria-live="polite" aria-atomic="false">
	{#each $toast as toastData (toastData.id)}
		<Toast data={toastData} />
	{/each}
</div>

<style>
	.toast-container {
		position: fixed;
		top: 1rem;
		right: 1rem;
		z-index: 9999;
		pointer-events: none;
		display: flex;
		flex-direction: column;
		align-items: flex-end;
	}

	.toast-container :global(.toast) {
		pointer-events: auto;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.toast-container {
			top: auto;
			bottom: 1rem;
			left: 1rem;
			right: 1rem;
			align-items: stretch;
		}
	}
</style>
