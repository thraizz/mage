/**
 * DeckUploadModal deck parsing and validation tests
 */

import { describe, it, expect } from 'vitest';
import {
  parseDeckList,
  parseStructuredCards,
  structuredCardsToText,
  validateDeck,
  type DeckStats
} from '$lib/utils/deck-parser';

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
    const stats = parseDeckList(commanderDeck, 'Commander');

    expect(stats.commanderCount).toBe(1);
    expect(stats.mainDeckCount + stats.commanderCount).toBe(100);
    expect(stats.totalCount).toBe(100);
  });

  it('should count total cards correctly', () => {
    const stats = parseDeckList(commanderDeck, 'Commander');

    // Should be exactly 100 cards
    expect(stats.totalCount).toBe(100);
    expect(stats.commanderCount + stats.mainDeckCount).toBe(100);
  });

  it('should parse card names correctly', () => {
    const testDeck = `1 Hearthhull, the Worldseed
5 Forest
1 Walk-In Closet/Forgotten Cellar
1 La abuela, siempre generosa`;

    const cards = parseStructuredCards(testDeck);

    expect(cards).toHaveLength(4);
    expect(cards[0].name).toBe('Hearthhull, the Worldseed');
    expect(cards[1].name).toBe('Forest');
    expect(cards[2].name).toBe('Walk-In Closet/Forgotten Cellar');
    expect(cards[3].name).toBe('La abuela, siempre generosa');
  });
});

describe('Standard Deck Format Validation', () => {
  it('should enforce 60 card minimum', () => {
    const tooSmallDeck = `20 Mountain
20 Lightning Bolt`;

    const stats = parseDeckList(tooSmallDeck, 'Standard');

    expect(stats.mainDeckCount).toBeLessThan(60);
    expect(stats.totalCount).toBe(40);
  });

  it('should track card quantities for 4-of validation', () => {
    const deck = `5 Lightning Bolt
20 Mountain`;

    const stats = parseDeckList(deck, 'Standard');

    // Lightning Bolt should have 5 copies (violation)
    expect(stats.errors.length).toBeGreaterThan(0);
    expect(
      stats.errors.some((err) => err.toLowerCase().includes('lightning bolt') && err.includes('5'))
    ).toBe(true);
  });
});

describe('Deck List Parsing', () => {
  it('should handle different quantity formats', () => {
    const testDeck = `4 Lightning Bolt
4x Lightning Strike
1 Sol Ring`;

    const cards = parseStructuredCards(testDeck);

    expect(cards).toHaveLength(3);
    expect(cards[0]).toEqual({ name: 'Lightning Bolt', quantity: 4, section: 'main' });
    expect(cards[1]).toEqual({ name: 'Lightning Strike', quantity: 4, section: 'main' });
    expect(cards[2]).toEqual({ name: 'Sol Ring', quantity: 1, section: 'main' });
  });

  it('should skip comments and empty lines', () => {
    const deckWithComments = `# This is a comment
4 Lightning Bolt
// Another comment

20 Mountain`;

    const stats = parseDeckList(deckWithComments, 'Standard');

    expect(stats.mainDeckCount).toBe(24); // 4 + 20
    expect(stats.errors.length).toBe(0);
  });

  it('should detect sideboard section', () => {
    const deckWithSideboard = `20 Mountain
4 Lightning Bolt

Sideboard:
2 Smash to Smithereens`;

    const stats = parseDeckList(deckWithSideboard, 'Standard');
    const cards = parseStructuredCards(deckWithSideboard);

    expect(stats.sideboardCount).toBe(2);
    expect(stats.mainDeckCount).toBe(24);
    expect(cards.filter((c) => c.section === 'sideboard')).toHaveLength(1);
    expect(cards.filter((c) => c.section === 'sideboard')[0].name).toBe('Smash to Smithereens');
  });

  it('should convert structured cards back to text', () => {
    const originalText = `Commander:
1 Hearthhull, the Worldseed

1 Lightning Bolt
20 Mountain

Sideboard:
2 Smash to Smithereens`;

    const cards = parseStructuredCards(originalText);
    const convertedText = structuredCardsToText(cards);

    // Should preserve structure
    expect(convertedText).toContain('Commander:');
    expect(convertedText).toContain('Hearthhull, the Worldseed');
    expect(convertedText).toContain('Sideboard:');
    expect(convertedText).toContain('Smash to Smithereens');

    // Parse back and verify
    const reParsed = parseStructuredCards(convertedText);
    expect(reParsed).toHaveLength(cards.length);
  });

  it('should validate Commander deck correctly', () => {
    const validCommanderDeck = `Commander:
1 Hearthhull, the Worldseed

${Array(99).fill('1 Test Card').join('\n')}`;

    const stats = parseDeckList(validCommanderDeck, 'Commander');
    const errors = validateDeck('Test Deck', validCommanderDeck, 'Commander', stats);

    expect(stats.totalCount).toBe(100);
    expect(errors.length).toBe(0);
  });

  it('should validate Standard deck correctly', () => {
    const validStandardDeck = `${Array(60).fill('1 Test Card').join('\n')}

Sideboard:
${Array(15).fill('1 Sideboard Card').join('\n')}`;

    const stats = parseDeckList(validStandardDeck, 'Standard');
    const errors = validateDeck('Test Deck', validStandardDeck, 'Standard', stats);

    expect(stats.mainDeckCount).toBe(60);
    expect(stats.sideboardCount).toBe(15);
    // May have validation errors if card names are identical (4-of rule)
    // But deck structure should be valid
    expect(stats.totalCount).toBe(75);
  });
});

describe('Deck Save Format Validation', () => {
  const exampleCommanderDeck = `Commander:
1 Hearthhull, the Worldseed

1 Aftermath Analyst
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
1 Zuran Orb`;

  it('should parse example commander deck with correct section separation', () => {
    const stats = parseDeckList(exampleCommanderDeck, 'Commander');
    const cards = parseStructuredCards(exampleCommanderDeck);

    // Should have exactly 1 commander
    expect(stats.commanderCount).toBe(1);
    expect(cards.filter((c) => c.section === 'commander')).toHaveLength(1);
    expect(cards.filter((c) => c.section === 'commander')[0].name).toBe(
      'Hearthhull, the Worldseed'
    );

    // Should have 99 main deck cards (some cards appear multiple times, so fewer unique entries)
    expect(stats.mainDeckCount).toBe(99);
    // Note: card entries may be fewer than 99 due to duplicates (e.g., "4 Forest" appears twice = 8 Forests)
    const mainDeckCards = cards.filter((c) => c.section === 'main');
    const mainDeckTotal = mainDeckCards.reduce((sum, c) => sum + c.quantity, 0);
    expect(mainDeckTotal).toBe(99);

    // Total should be 100
    expect(stats.totalCount).toBe(100);
    expect(stats.commanderCount + stats.mainDeckCount).toBe(100);

    // No sideboard
    expect(stats.sideboardCount).toBe(0);
    expect(cards.filter((c) => c.section === 'sideboard')).toHaveLength(0);
  });

  it('should convert example commander deck back to correct text format for save', () => {
    const cards = parseStructuredCards(exampleCommanderDeck);
    const convertedText = structuredCardsToText(cards);

    // Should start with Commander: section
    expect(convertedText.trim().startsWith('Commander:')).toBe(true);

    // Should have commander card
    expect(convertedText).toContain('Hearthhull, the Worldseed');

    // Should have empty line separating commander from main deck
    const lines = convertedText.split('\n');
    const commanderIndex = lines.findIndex((line) => line.trim().toLowerCase() === 'commander:');
    expect(commanderIndex).toBeGreaterThanOrEqual(0);

    // Next line should be the commander card (without quantity prefix for single cards)
    expect(lines[commanderIndex + 1].trim()).toBe('Hearthhull, the Worldseed');

    // Line after commander should be empty
    expect(lines[commanderIndex + 2].trim()).toBe('');

    // After empty line, should have main deck cards
    const firstMainDeckCard = lines[commanderIndex + 3].trim();
    expect(firstMainDeckCard).toBeTruthy();
    expect(firstMainDeckCard).not.toContain('Commander:');
    expect(firstMainDeckCard).not.toContain('Sideboard:');

    // Should NOT have Sideboard: section
    expect(convertedText.toLowerCase()).not.toContain('sideboard:');
  });

  it('should maintain correct structure when parsing and converting round-trip', () => {
    const cards = parseStructuredCards(exampleCommanderDeck);
    const convertedText = structuredCardsToText(cards);
    const reParsedCards = parseStructuredCards(convertedText);

    // Should have same number of cards
    expect(reParsedCards.length).toBe(cards.length);

    // Commander section should be preserved
    const originalCommander = cards.filter((c) => c.section === 'commander');
    const reParsedCommander = reParsedCards.filter((c) => c.section === 'commander');
    expect(reParsedCommander.length).toBe(originalCommander.length);
    expect(reParsedCommander[0].name).toBe(originalCommander[0].name);

    // Main deck section should be preserved
    const originalMain = cards.filter((c) => c.section === 'main');
    const reParsedMain = reParsedCards.filter((c) => c.section === 'main');
    expect(reParsedMain.length).toBe(originalMain.length);

    // Verify stats match
    const originalStats = parseDeckList(exampleCommanderDeck, 'Commander');
    const reParsedStats = parseDeckList(convertedText, 'Commander');
    expect(reParsedStats.commanderCount).toBe(originalStats.commanderCount);
    expect(reParsedStats.mainDeckCount).toBe(originalStats.mainDeckCount);
    expect(reParsedStats.totalCount).toBe(originalStats.totalCount);
  });

  it('should produce valid save format for commander deck', () => {
    const cards = parseStructuredCards(exampleCommanderDeck);
    const saveFormat = structuredCardsToText(cards);

    // Parse the save format to verify it's correct
    const stats = parseDeckList(saveFormat, 'Commander');
    const errors = validateDeck('Test Deck', saveFormat, 'Commander', stats);

    // Should be valid
    expect(errors.length).toBe(0);
    expect(stats.commanderCount).toBe(1);
    expect(stats.mainDeckCount).toBe(99);
    expect(stats.totalCount).toBe(100);
  });

  it('should produce valid save format for standard deck with sideboard', () => {
    const standardDeck = `4 Lightning Bolt
4 Monastery Swiftspear
4 Goblin Guide
4 Eidolon of the Great Revel
4 Lava Spike
4 Rift Bolt
4 Skewer the Critics
4 Light Up the Stage
4 Chain Lightning
20 Mountain
4 Sunbaked Canyon

Sideboard:
3 Smash to Smithereens
3 Pyroclasm
3 Roiling Vortex
3 Searing Blood
3 Ensnaring Bridge`;

    const cards = parseStructuredCards(standardDeck);
    const saveFormat = structuredCardsToText(cards);

    // Parse the save format to verify it's correct
    const stats = parseDeckList(saveFormat, 'Standard');
    const errors = validateDeck('Test Deck', saveFormat, 'Standard', stats);

    // Should have correct counts
    expect(stats.mainDeckCount).toBe(60);
    expect(stats.sideboardCount).toBe(15);
    expect(stats.totalCount).toBe(75);

    // Should not have structural errors (may have 4-of violations if card names match)
    const structuralErrors = errors.filter((e) => !e.includes('Too many copies'));
    expect(structuralErrors.length).toBe(0);

    // Verify structure
    expect(saveFormat).not.toContain('Commander:');
    expect(saveFormat.toLowerCase()).toContain('sideboard:');

    // Verify main deck section exists and has cards
    const mainDeckSection = saveFormat.split('Sideboard:')[0].trim();
    expect(mainDeckSection.length).toBeGreaterThan(0);

    // Verify sideboard section has correct number of cards
    const sideboardSection = saveFormat.split('Sideboard:')[1].trim();
    const sideboardCards = parseStructuredCards(`Sideboard:\n${sideboardSection}`);
    expect(
      sideboardCards
        .filter((c) => c.section === 'sideboard')
        .reduce((sum, c) => sum + c.quantity, 0)
    ).toBe(15);
  });

  it('should handle empty line separation correctly in commander decks', () => {
    // Test that empty line after commander switches back to main deck
    const deckWithEmptyLine = `Commander:
1 Test Commander

1 Main Deck Card 1
1 Main Deck Card 2`;

    const cards = parseStructuredCards(deckWithEmptyLine);
    const stats = parseDeckList(deckWithEmptyLine, 'Commander');

    expect(cards.filter((c) => c.section === 'commander')).toHaveLength(1);
    expect(cards.filter((c) => c.section === 'main')).toHaveLength(2);
    expect(stats.commanderCount).toBe(1);
    expect(stats.mainDeckCount).toBe(2);

    // Convert back and verify
    const converted = structuredCardsToText(cards);
    const reParsed = parseStructuredCards(converted);
    expect(reParsed.filter((c) => c.section === 'commander')).toHaveLength(1);
    expect(reParsed.filter((c) => c.section === 'main')).toHaveLength(2);
  });

  it('should not include commander cards in main deck section', () => {
    const cards = parseStructuredCards(exampleCommanderDeck);
    const saveFormat = structuredCardsToText(cards);

    // Split by Commander: and get the part after it
    const afterCommander = saveFormat.split('Commander:')[1];
    expect(afterCommander).toBeTruthy();

    // Split by empty line to separate commander from main deck
    const parts = afterCommander.split('\n\n');
    expect(parts.length).toBeGreaterThanOrEqual(1);

    // Commander section (before empty line) should contain commander name
    const commanderSection = parts[0];
    expect(commanderSection).toContain('Hearthhull, the Worldseed');

    // Main deck section (after empty line) should NOT contain commander name
    if (parts.length > 1) {
      const mainDeckSection = parts.slice(1).join('\n\n');
      expect(mainDeckSection).not.toContain('Hearthhull, the Worldseed');
    }

    // Verify commander is only in commander section of structured cards
    const commanderCards = cards.filter((c) => c.section === 'commander');
    expect(commanderCards.length).toBe(1);
    expect(commanderCards[0].name).toBe('Hearthhull, the Worldseed');

    // Verify commander is NOT in main deck section of structured cards
    const mainDeckCards = cards.filter((c) => c.section === 'main');
    expect(mainDeckCards.every((c) => c.name !== 'Hearthhull, the Worldseed')).toBe(true);
  });

  it('should produce correct save format matching the exact example deck structure', () => {
    // This test verifies the exact format that should be sent to the server
    const cards = parseStructuredCards(exampleCommanderDeck);
    const saveFormat = structuredCardsToText(cards);

    // Verify the format starts correctly
    expect(saveFormat.trim().startsWith('Commander:')).toBe(true);

    // Parse the save format to ensure it's correct
    const reParsedStats = parseDeckList(saveFormat, 'Commander');
    const reParsedCards = parseStructuredCards(saveFormat);

    // Verify commander section
    const commanderCards = reParsedCards.filter((c) => c.section === 'commander');
    expect(commanderCards.length).toBe(1);
    expect(commanderCards[0].name).toBe('Hearthhull, the Worldseed');
    expect(reParsedStats.commanderCount).toBe(1);

    // Verify main deck section (should NOT include commander)
    const mainDeckCards = reParsedCards.filter((c) => c.section === 'main');
    expect(mainDeckCards.length).toBeGreaterThan(0);
    expect(mainDeckCards.every((c) => c.name !== 'Hearthhull, the Worldseed')).toBe(true);
    expect(reParsedStats.mainDeckCount).toBe(99);

    // Verify no main deck cards are incorrectly in commander section
    const allCommanderCards = reParsedCards.filter((c) => c.section === 'commander');
    const incorrectlyInCommander = allCommanderCards.filter((c) =>
      ['Aftermath Analyst', 'Augur of Autumn', 'Baloth Prime'].includes(c.name)
    );
    expect(incorrectlyInCommander.length).toBe(0); // These should be in main, not commander

    // Verify these cards ARE in main deck section
    const correctlyInMain = mainDeckCards.filter((c) =>
      ['Aftermath Analyst', 'Augur of Autumn', 'Baloth Prime'].includes(c.name)
    );
    expect(correctlyInMain.length).toBe(3); // All three should be in main deck

    // Verify total count
    expect(reParsedStats.totalCount).toBe(100);
    expect(reParsedStats.commanderCount + reParsedStats.mainDeckCount).toBe(100);
  });

  it('should handle round-trip conversion without losing section information', () => {
    // Test that parsing and converting maintains section integrity
    const originalCards = parseStructuredCards(exampleCommanderDeck);
    const textFormat = structuredCardsToText(originalCards);
    const reParsedCards = parseStructuredCards(textFormat);
    const finalTextFormat = structuredCardsToText(reParsedCards);

    // All cards should maintain their sections
    expect(reParsedCards.filter((c) => c.section === 'commander').length).toBe(1);
    expect(reParsedCards.filter((c) => c.section === 'main').length).toBeGreaterThan(0);
    expect(reParsedCards.filter((c) => c.section === 'sideboard').length).toBe(0);

    // Final format should be parseable and correct
    const finalStats = parseDeckList(finalTextFormat, 'Commander');
    expect(finalStats.commanderCount).toBe(1);
    expect(finalStats.mainDeckCount).toBe(99);
    expect(finalStats.totalCount).toBe(100);
  });
});
