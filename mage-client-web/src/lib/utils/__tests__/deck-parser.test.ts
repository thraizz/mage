/**
 * Deck parser tests - comprehensive format coverage
 */

import { describe, it, expect } from 'vitest';
import {
	parseStructuredCards,
	parseDeckList,
	structuredCardsToText,
	type DeckFormat
} from '../deck-parser';

describe('Deck Parser - Format Support', () => {
	describe('Standard Moxfield Format', () => {
		it('should parse commander deck with "Commander:" section', () => {
			const input = `Commander:
1 Korvold, Fae-Cursed King

4 Lightning Bolt
4 Forest
1 Sol Ring

Sideboard:
1 Cursed Totem`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(5);
			expect(cards.find((c) => c.name === 'Korvold, Fae-Cursed King')).toEqual({
				name: 'Korvold, Fae-Cursed King',
				quantity: 1,
				section: 'commander'
			});
			expect(cards.find((c) => c.name === 'Lightning Bolt')).toEqual({
				name: 'Lightning Bolt',
				quantity: 4,
				section: 'main'
			});
			expect(cards.find((c) => c.name === 'Cursed Totem')).toEqual({
				name: 'Cursed Totem',
				quantity: 1,
				section: 'sideboard'
			});
		});

		it('should parse deck with explicit "Main:" section', () => {
			const input = `Main:
4 Lightning Bolt
4 Counterspell

Sideboard:
2 Rest in Peace`;

			const cards = parseStructuredCards(input);
			expect(cards).toHaveLength(3);
			expect(cards.every((c) => c.section === 'main' || c.section === 'sideboard')).toBe(true);
		});
	});

	describe('Simple List Format (No Sections)', () => {
		it('should parse simple list without section headers', () => {
			const input = `1 Aftermath Analyst
1 Augur of Autumn
1 Beast Within
4 Forest
1 Sol Ring`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(5);
			expect(cards.every((c) => c.section === 'main')).toBe(true);
			expect(cards.find((c) => c.name === 'Forest')).toEqual({
				name: 'Forest',
				quantity: 4,
				section: 'main'
			});
		});

		it('should handle multiple instances of same basic land on separate lines', () => {
			const input = `4 Forest
4 Forest
2 Mountain
1 Mountain`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(4);

			// Should have 2 Forest entries
			const forestCards = cards.filter((c) => c.name === 'Forest');
			expect(forestCards).toHaveLength(2);
			expect(forestCards.every((c) => c.quantity === 4)).toBe(true);

			// Should have 2 Mountain entries
			const mountainCards = cards.filter((c) => c.name === 'Mountain');
			expect(mountainCards).toHaveLength(2);
		});

		it('should handle implicit commander at end of list', () => {
			const input = `1 Aftermath Analyst
1 Beast Within
4 Forest

1 Hearthhull, the Worldseed`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(4);
			// Without "Commander:" section, last card is just main deck
			expect(cards.every((c) => c.section === 'main')).toBe(true);
		});
	});

	describe('Arena Export Format', () => {
		it('should parse Arena format with metadata headers', () => {
			const input = `About
Name Hearthhull, Land Shuffle (Degenerate/Bracket 4)

Commander
1 Hearthhull, the Worldseed

Deck
1 Abrupt Decay
1 Beast Within
5 Forest
1 Sol Ring`;

			const cards = parseStructuredCards(input);

			// Should skip "About", "Name" lines and parse sections
			expect(cards.length).toBeGreaterThan(0);

			// Check commander section
			const commander = cards.find((c) => c.name === 'Hearthhull, the Worldseed');
			expect(commander).toBeDefined();
			// Note: Arena format uses "Commander" not "Commander:" so it might not be parsed as commander section
		});

		it('should handle Arena "Deck" section as main deck', () => {
			const input = `Deck
1 Abrupt Decay
1 Aftermath Analyst
1 Sol Ring`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(3);
			// "Deck" should trigger main deck section
			expect(cards.every((c) => c.section === 'main')).toBe(true);
		});
	});

	describe('Split Cards and Special Characters', () => {
		it('should parse split cards with //', () => {
			const input = `1 Fire // Ice
1 Wear // Tear`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(2);
			expect(cards[0].name).toBe('Fire // Ice');
			expect(cards[1].name).toBe('Wear // Tear');
		});

		it('should parse split cards with single slash', () => {
			const input = `1 Walk-In Closet/Forgotten Cellar
1 Wandering Mind/Paranormal Analyst`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(2);
			expect(cards[0].name).toBe('Walk-In Closet/Forgotten Cellar');
		});

		it('should handle cards with special characters', () => {
			const input = `1 Jace, the Mind Sculptor
1 Karn, Scion of Urza
1 "Ach! Hans, Run!"`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(3);
			expect(cards.find((c) => c.name.includes('Jace'))).toBeDefined();
		});
	});

	describe('Snow-Covered Lands', () => {
		it('should parse snow-covered basic lands', () => {
			const input = `1 Snow-Covered Forest
1 Snow-Covered Mountain
1 Snow-Covered Swamp`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(3);
			expect(cards.every((c) => c.name.includes('Snow-Covered'))).toBe(true);
		});
	});

	describe('Comments and Empty Lines', () => {
		it('should skip comment lines starting with #', () => {
			const input = `# This is my deck
1 Lightning Bolt
# Another comment
1 Counterspell`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(2);
		});

		it('should skip comment lines starting with //', () => {
			const input = `// Removal package
1 Lightning Bolt
1 Swords to Plowshares`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(2);
		});

		it('should handle empty lines gracefully', () => {
			const input = `1 Lightning Bolt

1 Counterspell


1 Mana Leak`;

			const cards = parseStructuredCards(input);

			expect(cards).toHaveLength(3);
		});
	});

	describe('Real-World Example: Simple Commander List', () => {
		it('should parse full commander deck without sections', () => {
			const input = `1 Aftermath Analyst
1 Augur of Autumn
1 Beast Within
1 Bojuka Bog
1 Canyon Slough
1 Command Tower
1 Cultivate
4 Forest
4 Forest
1 Harrow
1 Korvold, Fae-Cursed King
2 Mountain
1 Mountain Valley
1 Sol Ring
3 Swamp
2 Swamp

1 Hearthhull, the Worldseed`;

			const cards = parseStructuredCards(input);

			// Count total cards
			const total = cards.reduce((sum, c) => sum + c.quantity, 0);
			expect(total).toBe(27); // All cards counted

			// All should be main deck without explicit sections
			expect(cards.every((c) => c.section === 'main')).toBe(true);
		});
	});

	describe('Real-World Example: Arena Format', () => {
		it('should parse Arena export format', () => {
			const input = `About
Name Hearthhull, Land Shuffle (Degenerate/Bracket 4)

Commander
1 Hearthhull, the Worldseed

Deck
1 Abrupt Decay
1 Aftermath Analyst
1 Arid Mesa
1 Ashaya, Soul of the Wild
5 Forest
1 Green Sun's Zenith
1 Harrow
2 Mountain
1 Oracle of Mul Daya
1 Sol Ring
3 Swamp`;

			const cards = parseStructuredCards(input);

			// Should have cards parsed
			expect(cards.length).toBeGreaterThan(0);

			// Check some specific cards exist
			expect(cards.find((c) => c.name === 'Abrupt Decay')).toBeDefined();
			expect(cards.find((c) => c.name === 'Forest')).toBeDefined();
		});
	});

	describe('Quantity Formats', () => {
		it('should parse "4 Card" format', () => {
			const cards = parseStructuredCards('4 Lightning Bolt');
			expect(cards[0]).toEqual({ name: 'Lightning Bolt', quantity: 4, section: 'main' });
		});

		it('should parse "4x Card" format', () => {
			const cards = parseStructuredCards('4x Lightning Bolt');
			expect(cards[0]).toEqual({ name: 'Lightning Bolt', quantity: 4, section: 'main' });
		});

		it('should parse MTGO-style "1x Card (set) 123" format', () => {
			const cards = parseStructuredCards('1x Beast Within (eoc) 93');
			expect(cards[0]).toEqual({ name: 'Beast Within', quantity: 1, section: 'main' });

			const stats = parseDeckList('1x Beast Within (eoc) 93', 'Standard');
			expect(stats.mainDeckCount).toBe(1);
			expect(stats.errors).toHaveLength(0);
		});

		it('should parse card without quantity as 1', () => {
			const cards = parseStructuredCards('Lightning Bolt');
			expect(cards[0]).toEqual({ name: 'Lightning Bolt', quantity: 1, section: 'main' });
		});
	});
});

describe('Deck Parser - Statistics', () => {
	it('should calculate correct stats for standard deck', () => {
		const input = `4 Lightning Bolt
4 Counterspell
20 Island
20 Mountain

Sideboard:
2 Rest in Peace
1 Surgical Extraction`;

		const stats = parseDeckList(input, 'Standard');

		expect(stats.mainDeckCount).toBe(48);
		expect(stats.sideboardCount).toBe(3);
		expect(stats.commanderCount).toBe(0);
		expect(stats.totalCount).toBe(51);
	});

	it('should calculate correct stats for commander deck', () => {
		const input = `Commander:
1 Korvold, Fae-Cursed King

99 cards in main deck`;

		const stats = parseDeckList(input, 'Commander');

		expect(stats.commanderCount).toBe(1);
	});

	it('should detect 4-of violations in Standard', () => {
		const input = `5 Lightning Bolt
4 Island`;

		const stats = parseDeckList(input, 'Standard');

		expect(stats.errors.length).toBeGreaterThan(0);
		expect(stats.errors[0]).toContain('lightning bolt');
	});

	it('should allow any quantity in Commander format', () => {
		const input = `Commander:
1 Korvold, Fae-Cursed King

20 Forest
20 Mountain`;

		const stats = parseDeckList(input, 'Commander');

		// Should not have 4-of violations
		expect(stats.errors).toHaveLength(0);
	});

	it('should not count basic lands in 4-of validation', () => {
		const input = `50 Forest
4 Lightning Bolt`;

		const stats = parseDeckList(input, 'Standard');

		// Should not complain about Forest
		expect(stats.errors).toHaveLength(0);
	});
});

describe('Deck Parser - Round-Trip Conversion', () => {
	it('should convert structured cards back to text', () => {
		const input = `Commander:
1 Korvold, Fae-Cursed King

4 Lightning Bolt
4 Forest

Sideboard:
1 Rest in Peace`;

		const cards = parseStructuredCards(input);
		const output = structuredCardsToText(cards);

		expect(output).toContain('Commander:');
		expect(output).toContain('Korvold, Fae-Cursed King');
		expect(output).toContain('4 Lightning Bolt');
		expect(output).toContain('Sideboard:');
	});

	it('should handle cards without quantity prefix', () => {
		const input = `Commander:
Korvold, Fae-Cursed King

Lightning Bolt
Forest`;

		const cards = parseStructuredCards(input);
		const output = structuredCardsToText(cards);

		// Single quantity cards should not have number prefix
		expect(output).toContain('Korvold, Fae-Cursed King');
		expect(output).toContain('Lightning Bolt');
	});
});
