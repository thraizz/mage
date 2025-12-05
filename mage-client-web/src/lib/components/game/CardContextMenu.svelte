<script lang="ts">
	import {
		moveCard,
		tapUntap,
		flipCard,
		transformCard,
		setCardCounter,
		modifyCardCounter,
		destroyToken
	} from '$lib/api/direct-actions';
	import { toast } from '$lib/stores/toast';

	interface Props {
		gameId: string;
		cardId: string;
		cardName: string;
		isTapped: boolean;
		isFaceDown: boolean;
		isToken: boolean;
		currentZone: string;
		onClose: () => void;
		position: { x: number; y: number };
	}

	let {
		gameId,
		cardId,
		cardName,
		isTapped,
		isFaceDown,
		isToken,
		currentZone,
		onClose,
		position
	}: Props = $props();

	let showCounterInput = $state(false);
	let counterType = $state('+1/+1');
	let counterAmount = $state(1);

	// Available zones for movement
	const zones = [
		{ id: 'HAND', label: 'Hand', icon: '✋' },
		{ id: 'BATTLEFIELD', label: 'Battlefield', icon: '⚔️' },
		{ id: 'GRAVEYARD', label: 'Graveyard', icon: '💀' },
		{ id: 'EXILE', label: 'Exile', icon: '🌀' },
		{ id: 'LIBRARY', label: 'Library (Top)', icon: '📚' },
		{ id: 'COMMAND', label: 'Command Zone', icon: '👑' }
	].filter((z) => z.id !== currentZone);

	// Common counter types
	const counterTypes = ['+1/+1', '-1/-1', 'loyalty', 'charge', 'time', 'verse', 'lore', 'custom'];

	async function handleTapUntap() {
		try {
			await tapUntap(gameId, cardId, !isTapped);
			onClose();
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handleFlip() {
		try {
			await flipCard(gameId, cardId, !isFaceDown);
			onClose();
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handleTransform() {
		try {
			await transformCard(gameId, cardId);
			onClose();
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handleMoveToZone(targetZone: string) {
		try {
			await moveCard(gameId, cardId, targetZone);
			toast.success(`Moved ${cardName} to ${targetZone.toLowerCase()}`);
			onClose();
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handleSetCounter() {
		try {
			await setCardCounter(gameId, cardId, counterType, counterAmount);
			toast.success(`Set ${counterAmount} ${counterType} counters`);
			showCounterInput = false;
			onClose();
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handleModifyCounter(delta: number) {
		try {
			await modifyCardCounter(gameId, cardId, counterType, delta);
			onClose();
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handleDestroyToken() {
		try {
			await destroyToken(gameId, cardId);
			toast.success(`Destroyed token: ${cardName}`);
			onClose();
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.context-menu')) {
			onClose();
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<div
	class="context-menu"
	style="left: {position.x}px; top: {position.y}px;"
	role="menu"
	onkeydown={(e) => e.key === 'Escape' && onClose()}
>
	<div class="menu-header">
		<span class="card-name">{cardName}</span>
		<button class="close-btn" onclick={onClose}>×</button>
	</div>

	<div class="menu-section">
		<span class="section-label">Card State</span>

		<button class="menu-item" onclick={handleTapUntap}>
			{#if isTapped}
				<span class="icon">↻</span> Untap
			{:else}
				<span class="icon">↺</span> Tap
			{/if}
		</button>

		<button class="menu-item" onclick={handleFlip}>
			{#if isFaceDown}
				<span class="icon">👁</span> Turn Face Up
			{:else}
				<span class="icon">🔒</span> Turn Face Down
			{/if}
		</button>

		<button class="menu-item" onclick={handleTransform}>
			<span class="icon">🔄</span> Transform
		</button>
	</div>

	<div class="menu-section">
		<span class="section-label">Counters</span>

		{#if showCounterInput}
			<div class="counter-input">
				<select bind:value={counterType} class="counter-select">
					{#each counterTypes as type}
						<option value={type}>{type}</option>
					{/each}
				</select>
				<input type="number" bind:value={counterAmount} min="0" max="99" class="counter-number" />
				<button class="counter-set-btn" onclick={handleSetCounter}>Set</button>
			</div>
			<div class="counter-quick">
				<button class="quick-counter" onclick={() => handleModifyCounter(-1)}>-1</button>
				<button class="quick-counter" onclick={() => handleModifyCounter(1)}>+1</button>
				<button class="quick-counter" onclick={() => handleModifyCounter(5)}>+5</button>
			</div>
		{:else}
			<button class="menu-item" onclick={() => (showCounterInput = true)}>
				<span class="icon">🔢</span> Add/Remove Counters
			</button>
		{/if}
	</div>

	<div class="menu-section">
		<span class="section-label">Move To</span>

		{#each zones as zone}
			<button class="menu-item" onclick={() => handleMoveToZone(zone.id)}>
				<span class="icon">{zone.icon}</span>
				{zone.label}
			</button>
		{/each}
	</div>

	{#if isToken}
		<div class="menu-section danger">
			<button class="menu-item danger" onclick={handleDestroyToken}>
				<span class="icon">💥</span> Destroy Token
			</button>
		</div>
	{/if}
</div>

<style>
	.context-menu {
		position: fixed;
		background: var(--surface-1, #1a1a2e);
		border: 1px solid var(--border-color, #444);
		border-radius: 8px;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
		min-width: 200px;
		max-width: 280px;
		z-index: 1000;
		font-family: var(--font-sans, system-ui);
		overflow: hidden;
	}

	.menu-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 10px 12px;
		background: var(--surface-2, #252540);
		border-bottom: 1px solid var(--border-color, #333);
	}

	.card-name {
		font-weight: 600;
		color: var(--text-color, #fff);
		font-size: 13px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.close-btn {
		background: none;
		border: none;
		color: var(--text-muted, #888);
		font-size: 18px;
		cursor: pointer;
		padding: 0;
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 4px;
	}

	.close-btn:hover {
		background: var(--surface-3, #353550);
		color: var(--text-color, #fff);
	}

	.menu-section {
		padding: 8px 0;
		border-bottom: 1px solid var(--border-color, #333);
	}

	.menu-section:last-child {
		border-bottom: none;
	}

	.menu-section.danger {
		background: rgba(255, 68, 68, 0.1);
	}

	.section-label {
		display: block;
		padding: 4px 12px 6px;
		font-size: 10px;
		text-transform: uppercase;
		color: var(--text-muted, #888);
		font-weight: 600;
	}

	.menu-item {
		display: flex;
		align-items: center;
		gap: 8px;
		width: 100%;
		padding: 8px 12px;
		border: none;
		background: transparent;
		color: var(--text-color, #fff);
		font-size: 13px;
		text-align: left;
		cursor: pointer;
		transition: background 0.1s ease;
	}

	.menu-item:hover {
		background: var(--surface-3, #353550);
	}

	.menu-item.danger {
		color: var(--danger-color, #ff4444);
	}

	.menu-item.danger:hover {
		background: rgba(255, 68, 68, 0.2);
	}

	.icon {
		font-size: 14px;
		width: 20px;
		text-align: center;
	}

	.counter-input {
		display: flex;
		gap: 6px;
		padding: 6px 12px;
	}

	.counter-select {
		flex: 1;
		padding: 6px 8px;
		border: 1px solid var(--border-color, #333);
		border-radius: 4px;
		background: var(--surface-2, #252540);
		color: var(--text-color, #fff);
		font-size: 12px;
	}

	.counter-number {
		width: 50px;
		padding: 6px 8px;
		border: 1px solid var(--border-color, #333);
		border-radius: 4px;
		background: var(--surface-2, #252540);
		color: var(--text-color, #fff);
		font-size: 12px;
		text-align: center;
	}

	.counter-set-btn {
		padding: 6px 12px;
		border: none;
		border-radius: 4px;
		background: var(--accent-color, #00d4ff);
		color: #000;
		font-size: 12px;
		font-weight: 600;
		cursor: pointer;
	}

	.counter-set-btn:hover {
		background: var(--accent-hover, #33ddff);
	}

	.counter-quick {
		display: flex;
		gap: 6px;
		padding: 4px 12px 8px;
	}

	.quick-counter {
		flex: 1;
		padding: 6px;
		border: 1px solid var(--border-color, #333);
		border-radius: 4px;
		background: var(--surface-2, #252540);
		color: var(--text-color, #fff);
		font-size: 12px;
		cursor: pointer;
	}

	.quick-counter:hover {
		background: var(--surface-3, #353550);
	}
</style>

