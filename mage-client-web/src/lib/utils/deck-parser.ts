export interface DeckStats {
	mainDeckCount: number;
	sideboardCount: number;
	commanderCount: number;
	totalCount: number;
	uniqueCards: number;
	errors: string[];
}

export interface CardEntry {
	name: string;
	quantity: number;
	section: 'commander' | 'main' | 'sideboard';
}

export type DeckFormat =
	| 'Standard'
	| 'Modern'
	| 'Commander'
	| 'Legacy'
	| 'Vintage'
	| 'Pioneer'
	| 'Pauper'
	| 'Historic';

const BASIC_LANDS = [
	'plains',
	'island',
	'swamp',
	'mountain',
	'forest',
	'wastes',
	'snow-covered plains',
	'snow-covered island',
	'snow-covered swamp',
	'snow-covered mountain',
	'snow-covered forest'
];

/**
 * Normalize common deck export formats into a canonical card name.
 *
 * Supports MTGO-style suffixes like:
 * - "Beast Within (EOC) 93"
 * - "Beast Within (eoc) 93"
 * - "Beast Within (EOC)"
 *
 * This strips the trailing "(SET) [collector]" portion so backend name validation works.
 */
export function normalizeImportedCardName(rawName: string): string {
	const name = rawName.trim();
	if (!name) return name;

	// Strip trailing: " (SET) 123" or " (SET)" (case-insensitive)
	// Keep it conservative: SET is 2-6 alnum chars, collector is 1+ digits (optional trailing letter).
	return name
		.replace(/\s+\(([a-z0-9]{2,6})\)(?:\s+\d+[a-z]?)?\s*$/i, '')
		.trim();
}

/**
 * Parse a deck list text into structured card entries
 */
export function parseStructuredCards(text: string): CardEntry[] {
	const cards: CardEntry[] = [];
	const lines = text.split('\n');
	let inSideboard = false;
	let inCommander = false;
	let hadCommanderCards = false;
	let commanderCardsProcessed = 0;

	for (let i = 0; i < lines.length; i++) {
		const line = lines[i].trim();

		// Empty line: if we were in commander section and had cards, switch back to main
		if (!line) {
			if (inCommander && hadCommanderCards) {
				inCommander = false;
				commanderCardsProcessed = 0;
			}
			continue;
		}

		// Skip comments
		if (line.startsWith('#') || line.startsWith('//')) {
			continue;
		}

		// Check for section markers
		const lowerLine = line.toLowerCase();

		// Skip Arena metadata lines (About, Name, etc.)
		if (lowerLine === 'about' || lowerLine.startsWith('name ')) {
			continue;
		}

		// Arena format: "Commander" (without colon) on its own line
		if (lowerLine === 'commander' || lowerLine === 'command zone') {
			inCommander = true;
			inSideboard = false;
			hadCommanderCards = false;
			commanderCardsProcessed = 0;
			continue;
		}

		// Arena format: "Deck" (without colon) on its own line
		if (lowerLine === 'deck' || lowerLine === 'main' || lowerLine === 'main deck') {
			inSideboard = false;
			inCommander = false;
			commanderCardsProcessed = 0;
			continue;
		}

		// Arena format: "Sideboard" (without colon) on its own line
		if (lowerLine === 'sideboard') {
			inSideboard = true;
			inCommander = false;
			commanderCardsProcessed = 0;
			continue;
		}

		// Traditional format with colons (broader match)
		if (
			(lowerLine.includes('commander') || lowerLine.includes('command zone')) &&
			lowerLine.includes(':')
		) {
			inCommander = true;
			inSideboard = false;
			hadCommanderCards = false;
			commanderCardsProcessed = 0;
			continue;
		}
		if (lowerLine.includes('sideboard') && lowerLine.includes(':')) {
			inSideboard = true;
			inCommander = false;
			commanderCardsProcessed = 0;
			continue;
		}
		if (
			(lowerLine.includes('main') || lowerLine.includes('deck')) &&
			lowerLine.includes(':') &&
			!lowerLine.includes('commander') &&
			!lowerLine.includes('sideboard')
		) {
			inSideboard = false;
			inCommander = false;
			commanderCardsProcessed = 0;
			continue;
		}

		// Parse card line
		const match = line.match(/^(\d+)x?\s+(.+)$/i);
		if (match) {
			const quantity = parseInt(match[1]);
			const cardName = normalizeImportedCardName(match[2].trim());
			if (cardName && quantity > 0) {
				// In Commander format: allow up to 2 commander cards total (for partner commanders)
				// After processing the first commander card, allow a second ONLY if it would make total exactly 2
				// Otherwise, switch to main deck
				if (inCommander && commanderCardsProcessed >= 2) {
					inCommander = false;
				} else if (inCommander && commanderCardsProcessed >= 1) {
					// We've processed at least 1 commander card
					// Allow second commander only if it would make total exactly 2 (for partners)
					if (commanderCardsProcessed + quantity !== 2) {
						inCommander = false;
					}
				}

				const section: 'commander' | 'main' | 'sideboard' = inCommander
					? 'commander'
					: inSideboard
						? 'sideboard'
						: 'main';
				cards.push({ name: cardName, quantity, section });
				if (inCommander) {
					hadCommanderCards = true;
					commanderCardsProcessed += quantity;
					// After processing commander cards, switch to main deck for next card if we've reached the limit
					if (commanderCardsProcessed >= 2) {
						inCommander = false;
					}
				}
			}
		} else if (line) {
			// Single card without quantity
			const cardName = normalizeImportedCardName(line.trim());
			if (cardName) {
				// In Commander format: allow up to 2 commander cards total (for partner commanders)
				// After processing the first commander card, allow a second ONLY if it would make total exactly 2
				// Otherwise, switch to main deck
				if (inCommander && commanderCardsProcessed >= 2) {
					inCommander = false;
				} else if (inCommander && commanderCardsProcessed >= 1) {
					// We've processed at least 1 commander card
					// Allow second commander only if it would make total exactly 2 (for partners)
					if (commanderCardsProcessed + 1 !== 2) {
						inCommander = false;
					}
				}

				const section: 'commander' | 'main' | 'sideboard' = inCommander
					? 'commander'
					: inSideboard
						? 'sideboard'
						: 'main';
				cards.push({ name: cardName, quantity: 1, section });
				if (inCommander) {
					hadCommanderCards = true;
					commanderCardsProcessed += 1;
					// After processing commander cards, switch to main deck for next card if we've reached the limit
					if (commanderCardsProcessed >= 2) {
						inCommander = false;
					}
				}
			}
		}
	}

	return cards;
}

/**
 * Convert structured card entries back to deck list text format
 */
export function structuredCardsToText(cards: CardEntry[]): string {
	const sections: { commander: string[]; main: string[]; sideboard: string[] } = {
		commander: [],
		main: [],
		sideboard: []
	};

	for (const card of cards) {
		const line = card.quantity === 1 ? card.name : `${card.quantity} ${card.name}`;
		sections[card.section].push(line);
	}

	const parts: string[] = [];
	if (sections.commander.length > 0) {
		parts.push('Commander:');
		parts.push(...sections.commander);
		parts.push('');
	}
	if (sections.main.length > 0) {
		parts.push(...sections.main);
	}
	if (sections.sideboard.length > 0) {
		parts.push('');
		parts.push('Sideboard:');
		parts.push(...sections.sideboard);
	}

	return parts.join('\n');
}

/**
 * Parse a deck list and return statistics
 */
export function parseDeckList(text: string, format: DeckFormat = 'Standard'): DeckStats {
	const lines = text.split('\n');
	let mainDeckCount = 0;
	let sideboardCount = 0;
	let commanderCount = 0;
	const uniqueMainCards = new Set<string>();
	const uniqueSideCards = new Set<string>();
	const uniqueCommanderCards = new Set<string>();
	const cardQuantities = new Map<string, number>();
	const parseErrors: string[] = [];
	let inSideboard = false;
	let inCommander = false;
	let hadCommanderCards = false;

	let commanderCardsProcessed = 0;

	for (let i = 0; i < lines.length; i++) {
		const line = lines[i].trim();

		// Empty line: if we were in commander section and had cards, switch back to main
		if (!line) {
			if (inCommander && hadCommanderCards) {
				inCommander = false;
				commanderCardsProcessed = 0;
			}
			continue;
		}

		// Skip comments
		if (line.startsWith('#') || line.startsWith('//')) {
			continue;
		}

		// Check for section markers - support multiple formats
		const lowerLine = line.toLowerCase();

		// Skip Arena metadata lines (About, Name, etc.)
		if (lowerLine === 'about' || lowerLine.startsWith('name ')) {
			continue;
		}

		// Arena format: "Commander" (without colon) on its own line
		if (lowerLine === 'commander' || lowerLine === 'command zone') {
			inCommander = true;
			inSideboard = false;
			hadCommanderCards = false;
			commanderCardsProcessed = 0;
			continue;
		}

		// Arena format: "Deck" (without colon) on its own line
		if (lowerLine === 'deck' || lowerLine === 'main' || lowerLine === 'main deck') {
			inSideboard = false;
			inCommander = false;
			commanderCardsProcessed = 0;
			continue;
		}

		// Arena format: "Sideboard" (without colon) on its own line
		if (lowerLine === 'sideboard') {
			inSideboard = true;
			inCommander = false;
			commanderCardsProcessed = 0;
			continue;
		}

		// Traditional format with colons (broader match)
		if (
			(lowerLine.includes('commander') || lowerLine.includes('command zone')) &&
			lowerLine.includes(':')
		) {
			inCommander = true;
			inSideboard = false;
			hadCommanderCards = false;
			commanderCardsProcessed = 0;
			continue;
		}
		if (lowerLine.includes('sideboard') && lowerLine.includes(':')) {
			inSideboard = true;
			inCommander = false;
			commanderCardsProcessed = 0;
			continue;
		}
		// Support "Main:", "Deck:", "Main Deck:" as explicit main deck markers
		if (
			(lowerLine.includes('main') || lowerLine.includes('deck')) &&
			lowerLine.includes(':') &&
			!lowerLine.includes('commander') &&
			!lowerLine.includes('sideboard')
		) {
			inSideboard = false;
			inCommander = false;
			commanderCardsProcessed = 0;
			continue;
		}

		// Parse card line: "4 Lightning Bolt" or "Lightning Bolt" or "4x Lightning Bolt"
		const match = line.match(/^(\d+)x?\s+(.+)$/i);
		if (match) {
			const quantity = parseInt(match[1]);
			const cardName = normalizeImportedCardName(match[2].trim());

			if (quantity <= 0 || quantity > 100) {
				parseErrors.push(`Line ${i + 1}: Invalid quantity (${quantity})`);
				continue;
			}

			if (!cardName) {
				parseErrors.push(`Line ${i + 1}: Missing card name`);
				continue;
			}

			// In Commander format: allow up to 2 commander cards total (for partner commanders)
			// After processing the first commander card, allow a second ONLY if it would make total exactly 2
			// Otherwise, switch to main deck
			if (inCommander && commanderCardsProcessed >= 2) {
				inCommander = false;
			} else if (inCommander && commanderCardsProcessed >= 1) {
				// We've processed at least 1 commander card
				// Allow second commander only if it would make total exactly 2 (for partners)
				if (commanderCardsProcessed + quantity !== 2) {
					inCommander = false;
				}
			}

			// Track quantities for 4-of validation (skip basic lands and commander zone)
			if (!inCommander && !inSideboard) {
				const lowerName = cardName.toLowerCase();
				const isBasicLand = BASIC_LANDS.includes(lowerName);

				if (!isBasicLand) {
					const currentQty = cardQuantities.get(lowerName) || 0;
					cardQuantities.set(lowerName, currentQty + quantity);
				}
			}

			if (inSideboard) {
				sideboardCount += quantity;
				uniqueSideCards.add(cardName.toLowerCase());
			} else if (inCommander) {
				commanderCount += quantity;
				uniqueCommanderCards.add(cardName.toLowerCase());
				hadCommanderCards = true;
				commanderCardsProcessed += quantity;
				// After processing commander cards, switch to main deck for next card
				if (commanderCardsProcessed >= 2) {
					inCommander = false;
				}
			} else {
				mainDeckCount += quantity;
				uniqueMainCards.add(cardName.toLowerCase());
			}
		} else if (line) {
			// Single card without quantity prefix
			const cardName = normalizeImportedCardName(line.trim());
			const lowerName = cardName.toLowerCase();

			// In Commander format: allow up to 2 commander cards total (for partner commanders)
			// After processing the first commander card, allow a second ONLY if it would make total exactly 2
			// Otherwise, switch to main deck
			if (inCommander && commanderCardsProcessed >= 2) {
				inCommander = false;
			} else if (inCommander && commanderCardsProcessed >= 1) {
				// We've processed at least 1 commander card
				// Allow second commander only if it would make total exactly 2 (for partners)
				if (commanderCardsProcessed + 1 !== 2) {
					inCommander = false;
				}
			}

			// Track for 4-of validation
			if (!inCommander && !inSideboard) {
				const isBasicLand = BASIC_LANDS.includes(lowerName);

				if (!isBasicLand) {
					const currentQty = cardQuantities.get(lowerName) || 0;
					cardQuantities.set(lowerName, currentQty + 1);
				}
			}

			if (inSideboard) {
				sideboardCount += 1;
				uniqueSideCards.add(lowerName);
			} else if (inCommander) {
				commanderCount += 1;
				uniqueCommanderCards.add(lowerName);
				hadCommanderCards = true;
				commanderCardsProcessed += 1;
				// After processing commander cards, switch to main deck for next card
				if (commanderCardsProcessed >= 2) {
					inCommander = false;
				}
			} else {
				mainDeckCount += 1;
				uniqueMainCards.add(lowerName);
			}
		}
	}

	// Check for 4-of violations (except in Commander format)
	if (format !== 'Commander') {
		cardQuantities.forEach((qty, cardName) => {
			if (qty > 4) {
				parseErrors.push(`Too many copies of "${cardName}" (${qty}/4 max)`);
			}
		});
	}

	return {
		mainDeckCount,
		sideboardCount,
		commanderCount,
		totalCount: mainDeckCount + sideboardCount + commanderCount,
		uniqueCards: uniqueMainCards.size + uniqueSideCards.size + uniqueCommanderCards.size,
		errors: parseErrors
	};
}

/**
 * Validate a deck based on format rules
 */
export function validateDeck(
	deckName: string,
	deckList: string,
	format: DeckFormat,
	stats: DeckStats
): string[] {
	const validationErrors: string[] = [];

	// Check deck name
	if (!deckName.trim()) {
		validationErrors.push('❌ Deck name is required');
	}

	// Check deck list
	if (!deckList.trim()) {
		validationErrors.push('❌ Deck list is required');
	}

	// Format-specific validation
	if (format === 'Commander') {
		// Commander: 1 commander + 99 other cards = 100 total
		const totalCards = stats.mainDeckCount + stats.commanderCount;
		const hasCommanderSection = deckList.toLowerCase().includes('commander:');

		if (totalCards !== 100) {
			if (hasCommanderSection) {
				validationErrors.push(
					`⚠️ Commander deck must be exactly 100 cards total (1 commander + 99 deck) (currently ${totalCards}: ${stats.commanderCount} commander + ${stats.mainDeckCount} main)`
				);
			} else {
				validationErrors.push(
					`⚠️ Commander deck must be exactly 100 cards. Add a "Commander:" section for your commander card.`
				);
			}
		}

		if (stats.commanderCount === 0 && hasCommanderSection) {
			validationErrors.push('⚠️ Commander section is empty. Add your commander card.');
		}

		if (stats.commanderCount > 1) {
			validationErrors.push(
				`⚠️ Commander deck can only have 1 commander (currently ${stats.commanderCount})`
			);
		}

		if (stats.sideboardCount > 0) {
			validationErrors.push('⚠️ Commander decks cannot have a sideboard');
		}
	} else if (
		format === 'Standard' ||
		format === 'Modern' ||
		format === 'Pioneer' ||
		format === 'Legacy' ||
		format === 'Vintage' ||
		format === 'Pauper'
	) {
		if (stats.mainDeckCount < 60) {
			validationErrors.push(
				`⚠️ Main deck must be at least 60 cards (currently ${stats.mainDeckCount})`
			);
		}
		if (stats.sideboardCount > 15) {
			validationErrors.push(
				`⚠️ Sideboard cannot exceed 15 cards (currently ${stats.sideboardCount})`
			);
		}
	} else if (format === 'Historic') {
		if (stats.mainDeckCount < 60) {
			validationErrors.push(
				`⚠️ Main deck must be at least 60 cards (currently ${stats.mainDeckCount})`
			);
		}
		if (stats.sideboardCount > 15) {
			validationErrors.push(
				`⚠️ Sideboard cannot exceed 15 cards (currently ${stats.sideboardCount})`
			);
		}
	}

	// Add parsing errors
	if (stats.errors.length > 0) {
		validationErrors.push(...stats.errors.map((err) => `🔴 ${err}`));
	}

	return validationErrors;
}
