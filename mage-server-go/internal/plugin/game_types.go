package plugin

// TwoPlayerDuel represents a standard 1v1 game
type TwoPlayerDuel struct{}

func init() {
	RegisterGameType(&TwoPlayerDuel{})
}

func (g *TwoPlayerDuel) Name() string {
	return "Two Player Duel"
}

func (g *TwoPlayerDuel) MinPlayers() int {
	return 2
}

func (g *TwoPlayerDuel) MaxPlayers() int {
	return 2
}

func (g *TwoPlayerDuel) Description() string {
	return "Standard 1v1 Magic game with 20 life"
}

func (g *TwoPlayerDuel) Rules() GameRules {
	return GameRules{
		StartingLife:            20,
		MinimumDeckSize:         60,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

func (g *TwoPlayerDuel) Behaviors() []GameBehavior {
	return nil // No special behaviors for standard games
}

// FreeForAll represents a multiplayer free-for-all game
type FreeForAll struct{}

func init() {
	RegisterGameType(&FreeForAll{})
}

func (g *FreeForAll) Name() string {
	return "Free For All"
}

func (g *FreeForAll) MinPlayers() int {
	return 2
}

func (g *FreeForAll) MaxPlayers() int {
	return 10
}

func (g *FreeForAll) Description() string {
	return "Multiplayer free-for-all game"
}

func (g *FreeForAll) Rules() GameRules {
	return GameRules{
		StartingLife:            20,
		MinimumDeckSize:         60,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

func (g *FreeForAll) Behaviors() []GameBehavior {
	return nil // No special behaviors
}

// CommanderFreeForAll represents a Commander format multiplayer game
type CommanderFreeForAll struct{}

func init() {
	RegisterGameType(&CommanderFreeForAll{})
}

func (g *CommanderFreeForAll) Name() string {
	return "Commander Free For All"
}

func (g *CommanderFreeForAll) MinPlayers() int {
	return 2
}

func (g *CommanderFreeForAll) MaxPlayers() int {
	return 10
}

func (g *CommanderFreeForAll) Description() string {
	return "Commander format multiplayer game with 40 life"
}

func (g *CommanderFreeForAll) Rules() GameRules {
	return GameRules{
		StartingLife:            40,
		MinimumDeckSize:         100,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

func (g *CommanderFreeForAll) Behaviors() []GameBehavior {
	return []GameBehavior{NewCommanderBehavior()}
}

// CommanderDuel represents a 1v1 Commander game
type CommanderDuel struct{}

func init() {
	RegisterGameType(&CommanderDuel{})
}

func (g *CommanderDuel) Name() string {
	return "Commander Duel"
}

func (g *CommanderDuel) MinPlayers() int {
	return 2
}

func (g *CommanderDuel) MaxPlayers() int {
	return 2
}

func (g *CommanderDuel) Description() string {
	return "1v1 Commander game with 40 life"
}

func (g *CommanderDuel) Rules() GameRules {
	return GameRules{
		StartingLife:            40,
		MinimumDeckSize:         100,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

func (g *CommanderDuel) Behaviors() []GameBehavior {
	return []GameBehavior{NewCommanderBehavior()}
}

// Brawl represents Brawl format
type Brawl struct{}

func init() {
	RegisterGameType(&Brawl{})
}

func (g *Brawl) Name() string {
	return "Brawl"
}

func (g *Brawl) MinPlayers() int {
	return 2
}

func (g *Brawl) MaxPlayers() int {
	return 10
}

func (g *Brawl) Description() string {
	return "Brawl format game with 25 life"
}

func (g *Brawl) Rules() GameRules {
	return GameRules{
		StartingLife:            25,
		MinimumDeckSize:         60,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

func (g *Brawl) Behaviors() []GameBehavior {
	// Brawl uses commander rules but without commander damage
	cb := NewCommanderBehavior()
	cb.CheckDamage = false
	return []GameBehavior{cb}
}

// CanadianHighlander represents Canadian Highlander format
type CanadianHighlander struct{}

func init() {
	RegisterGameType(&CanadianHighlander{})
}

func (g *CanadianHighlander) Name() string {
	return "Canadian Highlander"
}

func (g *CanadianHighlander) MinPlayers() int {
	return 2
}

func (g *CanadianHighlander) MaxPlayers() int {
	return 10
}

func (g *CanadianHighlander) Description() string {
	return "Canadian Highlander singleton format"
}

func (g *CanadianHighlander) Rules() GameRules {
	return GameRules{
		StartingLife:            20,
		MinimumDeckSize:         100,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

func (g *CanadianHighlander) Behaviors() []GameBehavior {
	return nil // No commander behavior - singleton but no commanders
}

// Commander is a convenience alias that maps to Commander Free For All
// This name matches what the web client sends
type Commander struct {
	CommanderFreeForAll
}

func init() {
	RegisterGameType(&Commander{})
}

func (g *Commander) Name() string {
	return "Commander"
}

func (g *Commander) Description() string {
	return "Commander format game with 40 life (alias for Commander Free For All)"
}

// Standard is a convenience alias for Two Player Duel
// This name matches what the web client sends
type Standard struct {
	TwoPlayerDuel
}

func init() {
	RegisterGameType(&Standard{})
}

func (g *Standard) Name() string {
	return "Standard"
}

func (g *Standard) Description() string {
	return "Standard format game with 20 life"
}

// Modern is a convenience alias for Two Player Duel
type Modern struct {
	TwoPlayerDuel
}

func init() {
	RegisterGameType(&Modern{})
}

func (g *Modern) Name() string {
	return "Modern"
}

func (g *Modern) Description() string {
	return "Modern format game with 20 life"
}

// Legacy is a convenience alias for Two Player Duel
type Legacy struct {
	TwoPlayerDuel
}

func init() {
	RegisterGameType(&Legacy{})
}

func (g *Legacy) Name() string {
	return "Legacy"
}

func (g *Legacy) Description() string {
	return "Legacy format game with 20 life"
}

// Vintage is a convenience alias for Two Player Duel
type Vintage struct {
	TwoPlayerDuel
}

func init() {
	RegisterGameType(&Vintage{})
}

func (g *Vintage) Name() string {
	return "Vintage"
}

func (g *Vintage) Description() string {
	return "Vintage format game with 20 life"
}

// Pauper is a convenience alias for Two Player Duel
type Pauper struct {
	TwoPlayerDuel
}

func init() {
	RegisterGameType(&Pauper{})
}

func (g *Pauper) Name() string {
	return "Pauper"
}

func (g *Pauper) Description() string {
	return "Pauper format game with 20 life"
}

// Pioneer is a convenience alias for Two Player Duel
type Pioneer struct {
	TwoPlayerDuel
}

func init() {
	RegisterGameType(&Pioneer{})
}

func (g *Pioneer) Name() string {
	return "Pioneer"
}

func (g *Pioneer) Description() string {
	return "Pioneer format game with 20 life"
}

// Historic is a convenience alias for Two Player Duel
type Historic struct {
	TwoPlayerDuel
}

func init() {
	RegisterGameType(&Historic{})
}

func (g *Historic) Name() string {
	return "Historic"
}

func (g *Historic) Description() string {
	return "Historic format game with 20 life"
}

// Alchemy is a convenience alias for Two Player Duel
type Alchemy struct {
	TwoPlayerDuel
}

func init() {
	RegisterGameType(&Alchemy{})
}

func (g *Alchemy) Name() string {
	return "Alchemy"
}

func (g *Alchemy) Description() string {
	return "Alchemy format game with 20 life"
}

// LimitedGameType is a convenience alias for Two Player Duel (for sealed/draft)
// Named differently from tournament type to avoid conflict
type LimitedGameType struct{}

func init() {
	RegisterGameType(&LimitedGameType{})
}

func (g *LimitedGameType) Name() string {
	return "Limited"
}

func (g *LimitedGameType) MinPlayers() int {
	return 2
}

func (g *LimitedGameType) MaxPlayers() int {
	return 2
}

func (g *LimitedGameType) Description() string {
	return "Limited format game with 20 life (40 card minimum)"
}

func (g *LimitedGameType) Rules() GameRules {
	return GameRules{
		StartingLife:            20,
		MinimumDeckSize:         40,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

func (g *LimitedGameType) Behaviors() []GameBehavior {
	return nil
}

// DraftGameType is a convenience alias for Limited
// Named differently from tournament type to avoid conflict
type DraftGameType struct {
	LimitedGameType
}

func init() {
	RegisterGameType(&DraftGameType{})
}

func (g *DraftGameType) Name() string {
	return "Draft"
}

func (g *DraftGameType) Description() string {
	return "Draft format game with 20 life (40 card minimum)"
}

// SealedGameType is a convenience alias for Limited
// Named differently from tournament type to avoid conflict
type SealedGameType struct {
	LimitedGameType
}

func init() {
	RegisterGameType(&SealedGameType{})
}

func (g *SealedGameType) Name() string {
	return "Sealed"
}

func (g *SealedGameType) Description() string {
	return "Sealed format game with 20 life (40 card minimum)"
}

// DuelCommander represents Duel Commander variant (different from regular commander duel)
type DuelCommander struct{}

func init() {
	RegisterGameType(&DuelCommander{})
}

func (g *DuelCommander) Name() string {
	return "Duel Commander"
}

func (g *DuelCommander) MinPlayers() int {
	return 2
}

func (g *DuelCommander) MaxPlayers() int {
	return 2
}

func (g *DuelCommander) Description() string {
	return "Duel Commander format with 20 life and no commander damage"
}

func (g *DuelCommander) Rules() GameRules {
	return GameRules{
		StartingLife:            20,
		MinimumDeckSize:         100,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

func (g *DuelCommander) Behaviors() []GameBehavior {
	// Duel Commander since Nov 2016 uses no commander damage rule
	cb := NewCommanderBehavior()
	cb.CheckDamage = false
	return []GameBehavior{cb}
}

// CenturionCommander represents Centurion Commander variant
type CenturionCommander struct{}

func init() {
	RegisterGameType(&CenturionCommander{})
}

func (g *CenturionCommander) Name() string {
	return "Centurion Commander"
}

func (g *CenturionCommander) MinPlayers() int {
	return 2
}

func (g *CenturionCommander) MaxPlayers() int {
	return 2
}

func (g *CenturionCommander) Description() string {
	return "Centurion Commander format with 25 life"
}

func (g *CenturionCommander) Rules() GameRules {
	return GameRules{
		StartingLife:            25,
		MinimumDeckSize:         100,
		StartingHandSize:        7,
		StartingPlayerSkipsDraw: true,
	}
}

func (g *CenturionCommander) Behaviors() []GameBehavior {
	return []GameBehavior{NewCommanderBehavior()}
}
