<script lang="ts">
	import type { Card } from '../types';

	interface Props {
		card: Card;
		onclick?: () => void;
	}

	let { card, onclick }: Props = $props();
</script>

<button
	class="card"
	class:tapped={card.tapped}
	class:attacking={card.attacking}
	class:blocking={card.blocking}
	onclick={onclick}
>
	<div class="card-header">
		<span class="card-name">{card.name}</span>
	</div>
	
	<div class="card-type">{card.type}</div>
	
	{#if card.power && card.toughness}
		<div class="card-pt">
			{card.power}/{card.toughness}
		</div>
	{/if}
	
	{#if card.damage > 0}
		<div class="card-damage">
			💔 {card.damage}
		</div>
	{/if}
	
	{#if card.abilities.length > 0}
		<div class="card-abilities">
			{#each card.abilities as ability}
				<span class="ability">{ability.id}</span>
			{/each}
		</div>
	{/if}
</button>

<style>
	.card {
		width: 150px;
		height: 210px;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border: 2px solid #333;
		border-radius: 8px;
		padding: 8px;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		flex-direction: column;
		position: relative;
		color: white;
		font-family: system-ui, -apple-system, sans-serif;
	}

	.card:hover {
		transform: translateY(-4px);
		box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
		border-color: #ffd700;
	}

	.card.tapped {
		transform: rotate(90deg);
		opacity: 0.7;
	}

	.card.attacking {
		border-color: #ff4444;
		box-shadow: 0 0 12px rgba(255, 68, 68, 0.6);
	}

	.card.blocking {
		border-color: #44ff44;
		box-shadow: 0 0 12px rgba(68, 255, 68, 0.6);
	}

	.card-header {
		font-weight: bold;
		font-size: 12px;
		margin-bottom: 4px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.card-type {
		font-size: 10px;
		opacity: 0.8;
		margin-bottom: 8px;
	}

	.card-pt {
		position: absolute;
		bottom: 8px;
		right: 8px;
		font-size: 18px;
		font-weight: bold;
		background: rgba(0, 0, 0, 0.5);
		padding: 4px 8px;
		border-radius: 4px;
	}

	.card-damage {
		position: absolute;
		top: 8px;
		right: 8px;
		font-size: 14px;
		font-weight: bold;
		background: rgba(255, 0, 0, 0.8);
		padding: 2px 6px;
		border-radius: 4px;
	}

	.card-abilities {
		display: flex;
		flex-wrap: wrap;
		gap: 4px;
		margin-top: auto;
	}

	.ability {
		font-size: 9px;
		background: rgba(255, 255, 255, 0.2);
		padding: 2px 4px;
		border-radius: 3px;
		text-transform: uppercase;
	}
</style>
