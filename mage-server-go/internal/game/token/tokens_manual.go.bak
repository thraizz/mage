package token

// This file contains predefined token types.
// Mirrors Java token implementations from mage.game.permanent.token package.

// ========================================
// Treasure Token
// ========================================

// NewTreasureToken creates a Treasure artifact token.
// "{T}, Sacrifice this artifact: Add one mana of any color."
func NewTreasureToken() *Token {
	return NewToken("Treasure Token", "Treasure token").
		AddCardType(CardTypeArtifact).
		AddSubtype("Treasure").
		AddAbility("treasure") // Treasure ability: tap and sacrifice for mana
}

// ========================================
// Saproling Token
// ========================================

// NewSaprolingToken creates a 1/1 green Saproling creature token.
func NewSaprolingToken() *Token {
	return NewToken("Saproling Token", "1/1 green Saproling creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Saproling").
		SetColor(Color{Green: true}).
		SetPowerToughness(1, 1)
}

// ========================================
// Squirrel Token
// ========================================

// NewSquirrelToken creates a 1/1 green Squirrel creature token.
func NewSquirrelToken() *Token {
	return NewToken("Squirrel Token", "1/1 green Squirrel creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Squirrel").
		SetColor(Color{Green: true}).
		SetPowerToughness(1, 1)
}

// ========================================
// Soldier Token
// ========================================

// NewSoldierToken creates a 1/1 white Soldier creature token.
func NewSoldierToken() *Token {
	return NewToken("Soldier Token", "1/1 white Soldier creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Soldier").
		SetColor(Color{White: true}).
		SetPowerToughness(1, 1)
}

// NewSoldierTokenVigilance creates a 1/1 white Soldier creature token with vigilance.
func NewSoldierTokenVigilance() *Token {
	return NewToken("Soldier Token", "1/1 white Soldier creature token with vigilance").
		AddCardType(CardTypeCreature).
		AddSubtype("Soldier").
		SetColor(Color{White: true}).
		SetPowerToughness(1, 1).
		AddAbility("vigilance")
}

// NewSoldierToken22 creates a 2/2 white Soldier creature token.
func NewSoldierToken22() *Token {
	return NewToken("Soldier Token", "2/2 white Soldier creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Soldier").
		SetColor(Color{White: true}).
		SetPowerToughness(2, 2)
}

// ========================================
// Goblin Token
// ========================================

// NewGoblinToken creates a 1/1 red Goblin creature token.
func NewGoblinToken() *Token {
	return NewToken("Goblin Token", "1/1 red Goblin creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Goblin").
		SetColor(Color{Red: true}).
		SetPowerToughness(1, 1)
}

// ========================================
// Elf Token
// ========================================

// NewElfWarriorToken creates a 1/1 green Elf Warrior creature token.
func NewElfWarriorToken() *Token {
	return NewToken("Elf Warrior Token", "1/1 green Elf Warrior creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Elf").
		AddSubtype("Warrior").
		SetColor(Color{Green: true}).
		SetPowerToughness(1, 1)
}

// NewElfDruidToken creates a 1/1 green Elf Druid creature token.
func NewElfDruidToken() *Token {
	return NewToken("Elf Druid Token", "1/1 green Elf Druid creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Elf").
		AddSubtype("Druid").
		SetColor(Color{Green: true}).
		SetPowerToughness(1, 1)
}

// ========================================
// Zombie Token
// ========================================

// NewZombieToken creates a 2/2 black Zombie creature token.
func NewZombieToken() *Token {
	return NewToken("Zombie Token", "2/2 black Zombie creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Zombie").
		SetColor(Color{Black: true}).
		SetPowerToughness(2, 2)
}

// NewZombieTokenDecayed creates a 2/2 black Zombie creature token with decayed.
func NewZombieTokenDecayed() *Token {
	return NewToken("Zombie Token", "2/2 black Zombie creature token with decayed").
		AddCardType(CardTypeCreature).
		AddSubtype("Zombie").
		SetColor(Color{Black: true}).
		SetPowerToughness(2, 2).
		AddAbility("decayed")
}

// ========================================
// Angel Token
// ========================================

// NewAngelToken creates a 4/4 white Angel creature token with flying.
func NewAngelToken() *Token {
	return NewToken("Angel Token", "4/4 white Angel creature token with flying").
		AddCardType(CardTypeCreature).
		AddSubtype("Angel").
		SetColor(Color{White: true}).
		SetPowerToughness(4, 4).
		AddAbility("flying")
}

// NewAngelTokenVigilance creates a 4/4 white Angel creature token with flying and vigilance.
func NewAngelTokenVigilance() *Token {
	return NewToken("Angel Token", "4/4 white Angel creature token with flying and vigilance").
		AddCardType(CardTypeCreature).
		AddSubtype("Angel").
		SetColor(Color{White: true}).
		SetPowerToughness(4, 4).
		AddAbility("flying").
		AddAbility("vigilance")
}

// ========================================
// Merfolk Token
// ========================================

// NewMerfolkToken creates a 1/1 blue Merfolk creature token.
func NewMerfolkToken() *Token {
	return NewToken("Merfolk Token", "1/1 blue Merfolk creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Merfolk").
		SetColor(Color{Blue: true}).
		SetPowerToughness(1, 1)
}

// NewMerfolkTokenHexproof creates a 1/1 blue Merfolk creature token with hexproof.
func NewMerfolkTokenHexproof() *Token {
	return NewToken("Merfolk Token", "1/1 blue Merfolk creature token with hexproof").
		AddCardType(CardTypeCreature).
		AddSubtype("Merfolk").
		SetColor(Color{Blue: true}).
		SetPowerToughness(1, 1).
		AddAbility("hexproof")
}

// ========================================
// Dragon Token
// ========================================

// NewDragonToken creates a 5/5 red Dragon creature token with flying.
func NewDragonToken() *Token {
	return NewToken("Dragon Token", "5/5 red Dragon creature token with flying").
		AddCardType(CardTypeCreature).
		AddSubtype("Dragon").
		SetColor(Color{Red: true}).
		SetPowerToughness(5, 5).
		AddAbility("flying")
}

// ========================================
// Beast Token
// ========================================

// NewBeastToken creates a 3/3 green Beast creature token.
func NewBeastToken() *Token {
	return NewToken("Beast Token", "3/3 green Beast creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Beast").
		SetColor(Color{Green: true}).
		SetPowerToughness(3, 3)
}

// NewBeastToken44 creates a 4/4 green Beast creature token.
func NewBeastToken44() *Token {
	return NewToken("Beast Token", "4/4 green Beast creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Beast").
		SetColor(Color{Green: true}).
		SetPowerToughness(4, 4)
}

// ========================================
// Spirit Token
// ========================================

// NewSpiritToken creates a 1/1 white Spirit creature token with flying.
func NewSpiritToken() *Token {
	return NewToken("Spirit Token", "1/1 white Spirit creature token with flying").
		AddCardType(CardTypeCreature).
		AddSubtype("Spirit").
		SetColor(Color{White: true}).
		SetPowerToughness(1, 1).
		AddAbility("flying")
}

// NewSpiritTokenBlue creates a 1/1 blue Spirit creature token with flying.
func NewSpiritTokenBlue() *Token {
	return NewToken("Spirit Token", "1/1 blue Spirit creature token with flying").
		AddCardType(CardTypeCreature).
		AddSubtype("Spirit").
		SetColor(Color{Blue: true}).
		SetPowerToughness(1, 1).
		AddAbility("flying")
}

// ========================================
// Elemental Token
// ========================================

// NewElementalToken creates a 4/4 red Elemental creature token.
func NewElementalToken() *Token {
	return NewToken("Elemental Token", "4/4 red Elemental creature token").
		AddCardType(CardTypeCreature).
		AddSubtype("Elemental").
		SetColor(Color{Red: true}).
		SetPowerToughness(4, 4)
}

// ========================================
// Thopter Token
// ========================================

// NewThopterToken creates a 1/1 colorless Thopter artifact creature token with flying.
func NewThopterToken() *Token {
	return NewToken("Thopter Token", "1/1 colorless Thopter artifact creature token with flying").
		AddCardType(CardTypeArtifact).
		AddCardType(CardTypeCreature).
		AddSubtype("Thopter").
		SetColor(Color{Colorless: true}).
		SetPowerToughness(1, 1).
		AddAbility("flying")
}

// ========================================
// Clue Token
// ========================================

// NewClueToken creates a Clue artifact token.
// "{2}, Sacrifice this artifact: Draw a card."
func NewClueToken() *Token {
	return NewToken("Clue Token", "Clue token").
		AddCardType(CardTypeArtifact).
		AddSubtype("Clue").
		AddAbility("clue") // Clue ability: pay 2 and sacrifice to draw
}

// ========================================
// Food Token
// ========================================

// NewFoodToken creates a Food artifact token.
// "{2}, {T}, Sacrifice this artifact: You gain 3 life."
func NewFoodToken() *Token {
	return NewToken("Food Token", "Food token").
		AddCardType(CardTypeArtifact).
		AddSubtype("Food").
		AddAbility("food") // Food ability: pay 2, tap, and sacrifice to gain 3 life
}

// ========================================
// Blood Token
// ========================================

// NewBloodToken creates a Blood artifact token.
// "{1}, {T}, Discard a card, Sacrifice this artifact: Draw a card."
func NewBloodToken() *Token {
	return NewToken("Blood Token", "Blood token").
		AddCardType(CardTypeArtifact).
		AddSubtype("Blood").
		AddAbility("blood") // Blood ability: pay 1, tap, discard, and sacrifice to draw
}

// ========================================
// Incubator Token
// ========================================

// NewIncubatorToken creates a 0/0 colorless Incubator artifact token with N +1/+1 counters.
// This is a helper for Phyrexia: All Will Be One tokens.
func NewIncubatorToken() *Token {
	return NewToken("Incubator Token", "0/0 colorless Incubator artifact token").
		AddCardType(CardTypeArtifact).
		AddSubtype("Incubator").
		SetColor(Color{Colorless: true}).
		SetPowerToughness(0, 0)
}

// ========================================
// Powerstone Token
// ========================================

// NewPowerstoneToken creates a Powerstone artifact token.
// "{T}: Add {C}. This mana can't be spent to cast a nonartifact spell."
func NewPowerstoneToken() *Token {
	return NewToken("Powerstone Token", "Powerstone token").
		AddCardType(CardTypeArtifact).
		AddSubtype("Powerstone").
		AddAbility("powerstone") // Powerstone ability: tap for colorless mana
}

// ========================================
// Map Token
// ========================================

// NewMapToken creates a Map artifact token.
// "{1}, {T}, Sacrifice this artifact: Target player searches their library for a land card, reveals it, puts it into their hand, then shuffles."
func NewMapToken() *Token {
	return NewToken("Map Token", "Map token").
		AddCardType(CardTypeArtifact).
		AddSubtype("Map").
		AddAbility("map") // Map ability
}
