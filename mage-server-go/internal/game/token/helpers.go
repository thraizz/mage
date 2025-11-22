package token

// Helper functions for commonly used token variants that aren't in the Java codebase.
// These are convenience functions that build on the generated tokens.

// NewSoldierTokenVigilance creates a 1/1 white Soldier creature token with vigilance.
func NewSoldierTokenVigilance() *Token {
	tok := NewSoldierToken() // Use the generated base
	tok.AddAbility("vigilance")
	return tok
}

// NewMerfolkTokenHexproof creates a 1/1 blue Merfolk creature token with hexproof.
func NewMerfolkTokenHexproof() *Token {
	tok := NewMerfolkToken() // Use the generated base
	tok.AddAbility("hexproof")
	return tok
}

// NewZombieTokenDecayed creates a 2/2 black Zombie creature token with decayed.
func NewZombieTokenDecayed() *Token {
	tok := NewZombieToken() // Use the generated base
	tok.AddAbility("decayed")
	return tok
}

// NewClueToken creates a Clue artifact token.
// "{2}, Sacrifice this artifact: Draw a card."
func NewClueToken() *Token {
	// Java has ClueArtifactToken
	return NewClueArtifactToken()
}
