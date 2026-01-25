/**
 * Deck parser tests - API format conversion
 */

import { describe, it, expect } from 'vitest';
import { parseStructuredCards, structuredCardsToText } from '../deck-parser';

describe('Deck Parser - API Format', () => {
  it('should convert parsed cards to normalized text format', () => {
    const input = `Commander
1 Korvold, Fae-Cursed King

Deck
4 Lightning Bolt
4 Counterspell
10 Island`;

    const cards = parseStructuredCards(input);

    // Verify cards were parsed
    expect(cards).toHaveLength(4);

    // Convert to normalized text
    const normalized = structuredCardsToText(cards);

    // Check format
    expect(normalized).toContain('Commander:');
    expect(normalized).toContain('Korvold, Fae-Cursed King'); // Single card without quantity
    expect(normalized).toContain('4 Lightning Bolt');
    expect(normalized).toContain('4 Counterspell');
    expect(normalized).toContain('10 Island');
  });

  it('should handle round-trip conversion', () => {
    const input = `Commander:
1 Korvold, Fae-Cursed King

4 Lightning Bolt
1 Sol Ring

Sideboard:
2 Rest in Peace`;

    const cards = parseStructuredCards(input);
    const output = structuredCardsToText(cards);
    const reparsed = parseStructuredCards(output);

    // Should maintain same structure
    expect(reparsed).toEqual(cards);
  });

  it('should normalize Arena format to standard format', () => {
    const input = `About
Name My Deck

Commander
1 Korvold, Fae-Cursed King

Deck
4 Lightning Bolt
5 Forest`;

    const cards = parseStructuredCards(input);
    const normalized = structuredCardsToText(cards);

    // Should output standard format with colons
    expect(normalized).toContain('Commander:');
    expect(normalized).not.toContain('About');
    expect(normalized).not.toContain('Name');

    // Cards should be preserved
    expect(cards.find((c) => c.name === 'Korvold, Fae-Cursed King')).toBeDefined();
    expect(cards.find((c) => c.name === 'Lightning Bolt')).toBeDefined();
    expect(cards.find((c) => c.name === 'Forest')).toBeDefined();
  });

  it('should handle simple list with no sections', () => {
    const input = `1 Beast Within
4 Forest
1 Sol Ring`;

    const cards = parseStructuredCards(input);
    const normalized = structuredCardsToText(cards);

    // Should just be the cards (no section headers)
    expect(normalized).not.toContain('Commander:');
    expect(normalized).not.toContain('Deck:');
    expect(normalized).toContain('Beast Within');
    expect(normalized).toContain('4 Forest');
    expect(normalized).toContain('Sol Ring');
  });

  it('should preserve card names with special characters', () => {
    const input = `1 Jace, the Mind Sculptor
1 Walk-In Closet/Forgotten Cellar
1 Fire // Ice
1 "Ach! Hans, Run!"`;

    const cards = parseStructuredCards(input);

    expect(cards[0].name).toBe('Jace, the Mind Sculptor');
    expect(cards[1].name).toBe('Walk-In Closet/Forgotten Cellar');
    expect(cards[2].name).toBe('Fire // Ice');
    expect(cards[3].name).toBe('"Ach! Hans, Run!"');

    const normalized = structuredCardsToText(cards);

    // Should preserve exact names
    expect(normalized).toContain('Jace, the Mind Sculptor');
    expect(normalized).toContain('Walk-In Closet/Forgotten Cellar');
    expect(normalized).toContain('Fire // Ice');
    expect(normalized).toContain('"Ach! Hans, Run!"');
  });

  it('should handle multiple basic lands on separate lines', () => {
    const input = `4 Forest
4 Forest
2 Mountain
1 Mountain`;

    const cards = parseStructuredCards(input);

    // Should create separate entries
    expect(cards).toHaveLength(4);

    // Convert back to text
    const normalized = structuredCardsToText(cards);

    // Should preserve separate entries
    const lines = normalized.split('\n').filter((l) => l.trim());
    expect(lines).toHaveLength(4);
  });

  it('should correctly format cards for API submission', () => {
    const input = `Commander:
Korvold, Fae-Cursed King

Lightning Bolt
Sol Ring
Forest`;

    const cards = parseStructuredCards(input);

    // Cards should have all required fields for API
    cards.forEach((card) => {
      expect(card).toHaveProperty('name');
      expect(card).toHaveProperty('quantity');
      expect(card).toHaveProperty('section');
      expect(card.name).toBeTruthy();
      expect(card.quantity).toBeGreaterThan(0);
      expect(['commander', 'main', 'sideboard']).toContain(card.section);
    });

    // Commander section
    expect(cards.filter((c) => c.section === 'commander')).toHaveLength(1);

    // Main deck
    const mainDeck = cards.filter((c) => c.section === 'main');
    expect(mainDeck).toHaveLength(3);
    expect(mainDeck.every((c) => c.quantity === 1)).toBe(true);
  });
});
