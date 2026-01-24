<script lang="ts">
	import '$lib/styles/variables.css';
	import '$lib/styles/global.css';
	import 'mana-font/css/mana.css';
	import favicon from '$lib/assets/favicon.svg';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import GlobalConfirmDialog from '$lib/components/GlobalConfirmDialog.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}

<!-- Global Toast Notification Container -->
<ToastContainer />

<!-- Global Confirmation Dialog -->
<GlobalConfirmDialog />
