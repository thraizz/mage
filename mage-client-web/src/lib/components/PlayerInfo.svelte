<script lang="ts">
	import type { Player } from '../types';

	interface Props {
		player: Player | undefined;
		isActive?: boolean;
	}

	let { player, isActive = false }: Props = $props();
</script>

{#if player}
	<div class="player-info" class:active={isActive}>
		<div class="player-name">{player.name}</div>
		<div class="player-life">❤️ {player.life}</div>
		<div class="player-zones">
			<span>📚 Library: {player.library_count}</span>
			<span>🎴 Hand: {player.hand_count}</span>
			<span>⚰️ Graveyard: {player.graveyard_count}</span>
		</div>
	</div>
{:else}
	<div class="player-info">
		<div class="player-name">Waiting for player...</div>
	</div>
{/if}

<style>
	.player-info {
		background: white;
		border: 2px solid #ddd;
		border-radius: 8px;
		padding: 16px;
		transition: all 0.2s ease;
	}

	.player-info.active {
		border-color: #4CAF50;
		box-shadow: 0 0 12px rgba(76, 175, 80, 0.3);
	}

	.player-name {
		font-size: 20px;
		font-weight: bold;
		margin-bottom: 8px;
	}

	.player-life {
		font-size: 32px;
		font-weight: bold;
		color: #e74c3c;
		margin-bottom: 12px;
	}

	.player-zones {
		display: flex;
		gap: 16px;
		font-size: 14px;
		color: #666;
	}
</style>
