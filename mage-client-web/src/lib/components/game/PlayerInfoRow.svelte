<script lang="ts">
	import Heart from '@lucide/svelte/icons/heart';
	import Skull from '@lucide/svelte/icons/skull';
	import LibraryZone from './LibraryZone.svelte';
	import Graveyard from './Graveyard.svelte';
	import ExileZone from './ExileZone.svelte';
	import ManaPoolComponent from './ManaPool.svelte';
	import type { CardView } from '$lib/generated/mage/v1/models';
	import type { Action } from 'svelte/action';

	interface Player {
		name: string;
		life: number;
		poison: number;
		libraryCount: number;
	}

	interface ManaPoolData {
		white: number;
		blue: number;
		black: number;
		red: number;
		green: number;
		colorless: number;
	}

	interface Props {
		player: Player;
		graveyard: CardView[];
		exile: CardView[];
		mana: ManaPoolData;
		showLifeMenu: boolean;
		onLifeChange: (delta: number) => void;
		onPoisonChange: (delta: number) => void;
		onToggleLifeMenu: () => void;
		onSearchLibrary: () => void;
		onDeckContextMenu: (e: MouseEvent) => void;
		libraryDropZoneRef?: (el: HTMLElement | null) => void;
		graveyardDropZoneRef?: (el: HTMLElement | null) => void;
		exileDropZoneRef?: (el: HTMLElement | null) => void;
	}

	let {
		player,
		graveyard,
		exile,
		mana,
		showLifeMenu,
		onLifeChange,
		onPoisonChange,
		onToggleLifeMenu,
		onSearchLibrary,
		onDeckContextMenu,
		libraryDropZoneRef,
		graveyardDropZoneRef,
		exileDropZoneRef
	}: Props = $props();

	let lifeMenuEl: HTMLDivElement | null = $state(null);

	// Create actions for drop zones
	const libraryDropZone: Action<HTMLElement> = (node) => {
		if (libraryDropZoneRef) libraryDropZoneRef(node);
		return {
			destroy() {
				if (libraryDropZoneRef) libraryDropZoneRef(null);
			}
		};
	};

	const graveyardDropZone: Action<HTMLElement> = (node) => {
		if (graveyardDropZoneRef) graveyardDropZoneRef(node);
		return {
			destroy() {
				if (graveyardDropZoneRef) graveyardDropZoneRef(null);
			}
		};
	};

	const exileDropZone: Action<HTMLElement> = (node) => {
		if (exileDropZoneRef) exileDropZoneRef(node);
		return {
			destroy() {
				if (exileDropZoneRef) exileDropZoneRef(null);
			}
		};
	};
</script>

<div class="player-info-row">
	<div class="player-identity">
		<span class="player-name">{player.name}</span>
	</div>

	<div class="player-stats-inline">
		<div class="life-group">
			<button class="stat-btn minus" onclick={() => onLifeChange(-1)}>−</button>
			<button class="stat-display life" onclick={onToggleLifeMenu}>
				<span class="stat-icon"><Heart size={14} /></span>
				<span class="stat-value">{player.life}</span>
			</button>
			<button class="stat-btn plus" onclick={() => onLifeChange(1)}>+</button>
		</div>

		{#if player.poison > 0}
			<div class="stat-display poison">
				<span class="stat-icon"><Skull size={14} /></span>
				<span class="stat-value">{player.poison}</span>
			</div>
		{/if}

		<div class="library-drop-zone" use:libraryDropZone>
			<LibraryZone
				libraryCount={player.libraryCount}
				playerName="You"
				onSearch={onSearchLibrary}
				onContextMenu={onDeckContextMenu}
			/>
		</div>

		{#if showLifeMenu}
			<div bind:this={lifeMenuEl} class="quick-menu">
				<div class="menu-section">
					<span class="menu-label">Life</span>
					<div class="menu-row">
						<button onclick={() => onLifeChange(-5)}>−5</button>
						<button onclick={() => onLifeChange(-1)}>−1</button>
						<button onclick={() => onLifeChange(1)}>+1</button>
						<button onclick={() => onLifeChange(5)}>+5</button>
					</div>
				</div>
				<div class="menu-section">
					<span class="menu-label">Poison</span>
					<div class="menu-row">
						<button onclick={() => onPoisonChange(-1)}>−1</button>
						<span class="menu-value">{player.poison}</span>
						<button onclick={() => onPoisonChange(1)}>+1</button>
					</div>
				</div>
				<button class="menu-close" onclick={onToggleLifeMenu}>✕</button>
			</div>
		{/if}
	</div>

	<div class="player-zones">
		<div class="graveyard-drop-zone" use:graveyardDropZone>
			<Graveyard
				cards={graveyard.map((c) => ({
					id: c.id,
					name: c.name,
					manaCost: c.manaCost,
					cardType: c.type,
					power: c.power,
					toughness: c.toughness,
					imageUrl: '',
					isTapped: false,
					isSelected: false
				}))}
				playerName="You"
				isOpponent={false}
				canDrag={true}
				onCardClick={() => {}}
			/>
		</div>
		<div class="exile-drop-zone" use:exileDropZone>
			<ExileZone
				cards={exile.map((c) => ({
					id: c.id,
					name: c.name,
					manaCost: c.manaCost,
					cardType: c.type,
					power: c.power,
					toughness: c.toughness,
					imageUrl: '',
					isTapped: false,
					isSelected: false
				}))}
				playerName="You"
				isOpponent={false}
				canDrag={true}
				onCardClick={() => {}}
				compact={true}
			/>
		</div>
		<ManaPoolComponent {mana} showEmpty={false} size="small" onManaClick={() => {}} />
	</div>
</div>
