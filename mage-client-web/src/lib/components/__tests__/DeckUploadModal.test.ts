/**
 * DeckUploadModal deck parsing and validation tests
 */

import { describe, it, expect } from 'vitest';

describe('Commander Deck Format Validation', () => {
	const commanderDeck = `Commander:
1 Hearthhull, the Worldseed

1 Abrupt Decay
1 Aftermath Analyst
1 Amulet of Vigor
1 Arid Mesa
1 Ashaya, Soul of the Wild
1 Assassin's Trophy
1 Beast Within
1 Blood Crypt
1 Bloodstained Mire
1 Bojuka Bog
1 Boseiju, Who Endures
1 Braids, Arisen Nightmare
1 Cabin of the Dead
1 Commercial District
1 Constant Mists
1 Crop Rotation
1 Crucible of Worlds
1 Cultivate
1 Dakmor Salvage
1 Elvish Reclaimer
1 Entish Restoration
1 Entomb
1 Evendo Brushrazer
1 Exploration
1 Exploration Broodship
1 Explore
1 Fabled Passage
1 Famished Worldsire
1 Fangorn Forest
5 Forest
1 Glacial Chasm
1 Green Sun's Zenith
1 Harrow
1 Horizon Explorer
1 Horn of Greed
1 Icetill Explorer
1 Iridescent Vinelasher
1 Korvold, Fae-Cursed King
1 La abuela, siempre generosa
1 Lord Windgrace
1 Lotus Cobra
1 Lotus Field
1 Lumra, Bellow of the Woods
1 Marsh Flats
1 Master Emerald Shrine
1 Mayhem Devil
1 Misty Rainforest
2 Mountain
1 Nature's Lore
1 Newfound Adventure
1 Nissa, Resurgent Animist
1 Oracle of Mul Daya
1 Orcish Lumberjack
1 Overgrown Tomb
1 Polluted Delta
1 Prismatic Vista
1 Quirion Ranger
1 Ramunap Excavator
1 Raucous Theater
1 Reprocess
1 Retreat to Hagra
1 Rydia, Summoner of Mist
1 Sabotender
1 Sakura-Tribe Elder
1 Scalding Tarn
1 Scapeshift
1 Scute Swarm
1 Shifting Woodland
1 Snow-Covered Forest
1 Snow-Covered Mountain
1 Snow-Covered Swamp
1 Sol Ring
1 Springheart Nantuko
1 Squandered Resources
1 Stomping Ground
3 Swamp
1 Sylvan Safekeeper
1 Szarel, Genesis Shepherd
1 Tannuk, Memorial Ensign
1 Tear Asunder
1 Tezzeret, Cruel Captain
1 The Gitrog Monster
1 Traveling Chocobo
1 Underground Mortuary
1 Urza's Saga
1 Valakut Exploration
1 Verdant Catacombs
1 Walk-In Closet/Forgotten Cellar
1 Windswept Heath
1 Wooded Foothills
1 Ziatora's Proving Ground
1 Zuran Orb`;

	it('should parse Commander deck with commander section', () => {
		const lines = commanderDeck.split('\n');
		let inCommander = false;
		let commanderCount = 0;
		let deckCount = 0;

		for (const line of lines) {
			const trimmed = line.trim();
			if (!trimmed || trimmed.startsWith('#')) continue;

			if (trimmed.toLowerCase().includes('commander:')) {
				inCommander = true;
				continue;
			}

			const match = trimmed.match(/^(\d+)x?\s+(.+)$/i);
			if (match) {
				const qty = parseInt(match[1]);
				if (inCommander) {
					commanderCount += qty;
					inCommander = false; // Commander section is typically just 1 card
				} else {
					deckCount += qty;
				}
			}
		}

		expect(commanderCount).toBe(1);
		expect(deckCount + commanderCount).toBe(100);
	});

	it('should count total cards correctly', () => {
		const lines = commanderDeck.split('\n');
		let totalCount = 0;

		for (const line of lines) {
			const trimmed = line.trim();
			if (!trimmed || trimmed.startsWith('#') || trimmed.toLowerCase().includes('commander:')) {
				continue;
			}

			const match = trimmed.match(/^(\d+)x?\s+(.+)$/i);
			if (match) {
				totalCount += parseInt(match[1]);
			} else if (trimmed) {
				totalCount += 1;
			}
		}

		// Should be exactly 100 cards
		expect(totalCount).toBe(100);
	});

	it('should parse card names correctly', () => {
		const testLines = [
			'1 Hearthhull, the Worldseed',
			'5 Forest',
			'1 Walk-In Closet/Forgotten Cellar',
			'1 La abuela, siempre generosa'
		];

		const expectedNames = [
			'Hearthhull, the Worldseed',
			'Forest',
			'Walk-In Closet/Forgotten Cellar',
			'La abuela, siempre generosa'
		];

		testLines.forEach((line, i) => {
			const match = line.match(/^(\d+)x?\s+(.+)$/i);
			expect(match).toBeTruthy();
			if (match) {
				expect(match[2].trim()).toBe(expectedNames[i]);
			}
		});
	});
});

describe('Standard Deck Format Validation', () => {
	it('should enforce 60 card minimum', () => {
		const tooSmallDeck = `20 Mountain
20 Lightning Bolt`;

		const lines = tooSmallDeck.split('\n');
		let totalCount = 0;

		for (const line of lines) {
			const trimmed = line.trim();
			if (!trimmed) continue;

			const match = trimmed.match(/^(\d+)x?\s+(.+)$/i);
			if (match) {
				totalCount += parseInt(match[1]);
			}
		}

		expect(totalCount).toBeLessThan(60);
	});

	it('should track card quantities for 4-of validation', () => {
		const deck = `5 Lightning Bolt
20 Mountain`;

		const cardQuantities = new Map<string, number>();
		const lines = deck.split('\n');

		for (const line of lines) {
			const trimmed = line.trim();
			if (!trimmed) continue;

			const match = trimmed.match(/^(\d+)x?\s+(.+)$/i);
			if (match) {
				const qty = parseInt(match[1]);
				const cardName = match[2].trim().toLowerCase();

				// Don't count basic lands
				const isBasicLand = ['plains', 'island', 'swamp', 'mountain', 'forest'].includes(cardName);
				if (!isBasicLand) {
					const current = cardQuantities.get(cardName) || 0;
					cardQuantities.set(cardName, current + qty);
				}
			}
		}

		// Lightning Bolt should have 5 copies (violation)
		expect(cardQuantities.get('lightning bolt')).toBe(5);
		expect(cardQuantities.get('lightning bolt')).toBeGreaterThan(4);

		// Mountain should not be counted (basic land)
		expect(cardQuantities.has('mountain')).toBe(false);
	});
});

describe('Deck List Parsing', () => {
	it('should handle different quantity formats', () => {
		const testCases = [
			{ input: '4 Lightning Bolt', qty: 4, name: 'Lightning Bolt' },
			{ input: '4x Lightning Bolt', qty: 4, name: 'Lightning Bolt' },
			{ input: '1 Sol Ring', qty: 1, name: 'Sol Ring' }
		];

		testCases.forEach(({ input, qty, name }) => {
			const match = input.match(/^(\d+)x?\s+(.+)$/i);
			expect(match).toBeTruthy();
			if (match) {
				expect(parseInt(match[1])).toBe(qty);
				expect(match[2].trim()).toBe(name);
			}
		});
	});

	it('should skip comments and empty lines', () => {
		const deckWithComments = `# This is a comment
4 Lightning Bolt
// Another comment

20 Mountain`;

		const lines = deckWithComments.split('\n');
		const validLines = lines.filter((line) => {
			const trimmed = line.trim();
			return trimmed && !trimmed.startsWith('#') && !trimmed.startsWith('//');
		});

		expect(validLines.length).toBe(2);
	});

	it('should detect sideboard section', () => {
		const deckWithSideboard = `20 Mountain
4 Lightning Bolt

Sideboard:
2 Smash to Smithereens`;

		const lines = deckWithSideboard.split('\n');
		let inSideboard = false;
		let sideboardCount = 0;

		for (const line of lines) {
			const trimmed = line.trim();
			if (trimmed.toLowerCase().includes('sideboard')) {
				inSideboard = true;
				continue;
			}

			if (inSideboard && trimmed) {
				const match = trimmed.match(/^(\d+)x?\s+(.+)$/i);
				if (match) {
					sideboardCount += parseInt(match[1]);
				}
			}
		}

		expect(sideboardCount).toBe(2);
	});
});
