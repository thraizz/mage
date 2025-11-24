/**
 * Deck parser tests - Real user examples
 */

import { describe, it, expect } from 'vitest';
import { parseStructuredCards, parseDeckList } from '../deck-parser';

describe('Deck Parser - Real User Examples', () => {
	it('should parse simple commander list without sections', () => {
		const input = `1 Aftermath Analyst
1 Augur of Autumn
1 Baloth Prime
1 Beast Within
1 Binding the Old Gods
1 Blasphemous Act
1 Bojuka Bog
1 Braids, Arisen Nightmare
1 Cabaretti Courtyard
1 Canyon Slough
1 Cinder Glade
1 Command Tower
1 Cultivate
1 Dakmor Salvage
1 Escape to the Wilds
1 Escape Tunnel
1 Eumidian Hatchery
1 Eumidian Wastewaker
1 Evendo Brushrazer
1 Evolving Wilds
1 Exploration Broodship
1 Fabled Passage
1 Farseek
1 Festering Thicket
4 Forest
4 Forest
1 Formless Genesis
1 Gaze of Granite
1 Greater Gargadon
1 Harrow
1 Horizon Explorer
1 Infernal Grasp
1 Juri, Master of the Revue
1 Karplusan Forest
1 Korvold, Fae-Cursed King
1 Llanowar Wastes
1 Lotus Cobra
1 Maestros Theater
1 Manifold Key
1 Mayhem Devil
1 Mazirek, Kraul Death Priest
1 Moraug, Fury of Akoum
1 Mountain
2 Mountain
1 Mountain Valley
1 Myriad Landscape
1 Nature's Lore
1 Night's Whisper
1 Omnath, Locus of Rage
1 Oracle of Mul Daya
1 Pest Infestation
1 Planetary Annihilation
1 Putrefy
1 Rakdos Charm
1 Rampaging Baloths
1 Ramunap Excavator
1 Riveteers Overlook
1 Rocky Tar Pit
1 Roiling Regrowth
1 Satyr Wayfinder
1 Scouring Swarm
1 Sheltered Thicket
1 Smoldering Marsh
1 Sol Ring
1 Soul of Windgrace
1 Splendid Reclamation
1 Springbloom Druid
1 Sprouting Goblin
1 Sulfurous Springs
3 Swamp
2 Swamp
1 Sylvan Safekeeper
1 Szarel, Genesis Shepherd
1 Tear Asunder
1 Terramorphic Expanse
1 The Gitrog Monster
1 Tiller Engine
1 Tireless Tracker
1 Titania, Protector of Argoth
1 Twilight Mire
1 Valakut Exploration
1 Vernal Fen
1 Viridescent Bog
1 Walk-In Closet/Forgotten Cellar
1 Wastes
1 Windgrace's Judgment
1 Worldsoul's Rage
1 Zask, Skittering Swarmlord
1 Zuran Orb

1 Hearthhull, the Worldseed`;

		const cards = parseStructuredCards(input);

		// Check that all cards are parsed
		expect(cards.length).toBeGreaterThan(80);

		// Check specific cards
		expect(cards.find((c) => c.name === 'Korvold, Fae-Cursed King')).toBeDefined();
		expect(cards.find((c) => c.name === 'Hearthhull, the Worldseed')).toBeDefined();
		expect(cards.find((c) => c.name === 'Walk-In Closet/Forgotten Cellar')).toBeDefined();

		// Check Forest entries
		const forestCards = cards.filter((c) => c.name === 'Forest');
		expect(forestCards).toHaveLength(2);
		expect(forestCards.every((c) => c.quantity === 4)).toBe(true);

		// Check Mountain entries
		const mountainCards = cards.filter((c) => c.name === 'Mountain');
		expect(mountainCards).toHaveLength(2);
		expect(mountainCards[0].quantity).toBe(1);
		expect(mountainCards[1].quantity).toBe(2);

		// Check Swamp entries
		const swampCards = cards.filter((c) => c.name === 'Swamp');
		expect(swampCards).toHaveLength(2);
		expect(swampCards[0].quantity).toBe(3);
		expect(swampCards[1].quantity).toBe(2);

		// Calculate total card count
		const totalCards = cards.reduce((sum, c) => sum + c.quantity, 0);
		expect(totalCards).toBe(100); // Should be 100 for Commander

		// All should be main deck without explicit sections
		expect(cards.every((c) => c.section === 'main')).toBe(true);
	});

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
1 Assassin's Trophy
1 Beast Within
1 Blood Crypt
1 Bloodstained Mire
1 Bojuka Bog
1 Boseiju, Who Endures
1 Braids, Arisen Nightmare
1 Command Tower
1 Commercial District
1 Crop Rotation
1 Crucible of Worlds
1 Cultivate
1 Elvish Reclaimer
1 Entish Restoration
1 Entomb
1 Exploration
1 Explore
1 Fabled Passage
1 Famished Worldsire
1 Farseek
1 Field of the Dead
5 Forest
1 Green Sun's Zenith
1 Harrow
1 Icetill Explorer
1 Iridescent Vinelasher
1 Korvold, Fae-Cursed King
1 Lotus Cobra
1 Lotus Field
1 Lumra, Bellow of the Woods
1 Marsh Flats
1 Mayhem Devil
1 Misty Rainforest
2 Mountain
1 Nissa, Resurgent Animist
1 Oracle of Mul Daya
1 Overgrown Tomb
1 Polluted Delta
1 Prismatic Vista
1 Ramunap Excavator
1 Raucous Theater
1 Rydia, Summoner of Mist
1 Sabotender
1 Scalding Tarn
1 Scapeshift
1 Scute Swarm
1 Shifting Woodland
1 Snow-Covered Forest
1 Snow-Covered Mountain
1 Snow-Covered Swamp
1 Sol Ring
1 Springheart Nantuko
1 Stomping Ground
3 Swamp
1 Sylvan Safekeeper
1 Szarel, Genesis Shepherd
1 Tannuk, Memorial Ensign
1 Tear Asunder
1 Tezzeret, Cruel Captain
1 The Gitrog Monster
1 Tireless Provisioner
1 Traveling Chocobo
1 Underground Mortuary
1 Valakut Exploration
1 Verdant Catacombs
1 Walk-In Closet // Forgotten Cellar
1 Windswept Heath
1 Wooded Foothills
1 Yavimaya, Cradle of Growth
1 Ziatora's Proving Ground`;

		const cards = parseStructuredCards(input);

		// Check that all cards are parsed
		expect(cards.length).toBeGreaterThan(70);

		// Check commander section
		const commander = cards.find((c) => c.name === 'Hearthhull, the Worldseed');
		expect(commander).toBeDefined();
		expect(commander?.section).toBe('commander');

		// Check main deck cards
		const abruptDecay = cards.find((c) => c.name === 'Abrupt Decay');
		expect(abruptDecay).toBeDefined();
		expect(abruptDecay?.section).toBe('main');

		// Check split cards with //
		const walkIn = cards.find((c) => c.name === 'Walk-In Closet // Forgotten Cellar');
		expect(walkIn).toBeDefined();

		// Check snow-covered lands
		const snowForest = cards.find((c) => c.name === 'Snow-Covered Forest');
		expect(snowForest).toBeDefined();

		// Calculate total card count (note: this is a partial deck example)
		const totalCards = cards.reduce((sum, c) => sum + c.quantity, 0);
		expect(totalCards).toBeGreaterThan(70); // Partial deck for testing

		// Check sections distribution
		const commanderCards = cards.filter((c) => c.section === 'commander');
		const mainCards = cards.filter((c) => c.section === 'main');
		expect(commanderCards.length).toBe(1);
		expect(mainCards.length).toBeGreaterThan(70);
	});

	it('should parse simple list and compute stats correctly', () => {
		const input = `1 Aftermath Analyst
1 Beast Within
4 Forest
1 Hearthhull, the Worldseed`;

		const stats = parseDeckList(input, 'Commander');

		expect(stats.totalCount).toBe(7);
		expect(stats.mainDeckCount).toBe(7); // All main without commander section
		expect(stats.commanderCount).toBe(0); // No explicit commander section
	});

	it('should parse Arena format and compute stats correctly', () => {
		const input = `Commander
1 Hearthhull, the Worldseed

Deck
1 Abrupt Decay
1 Beast Within
5 Forest`;

		const stats = parseDeckList(input, 'Commander');

		expect(stats.commanderCount).toBe(1);
		expect(stats.mainDeckCount).toBe(7);
		expect(stats.totalCount).toBe(8);
	});
});
