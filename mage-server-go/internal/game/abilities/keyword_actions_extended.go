package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

// This file implements additional keyword actions beyond Scry, Mill, and Proliferate
// Rule 701: Keyword Actions

// ===== Explore (Rule 701.46) =====

// ExploreEffect represents the Explore keyword action
// Rule 701.46a: To explore, a player reveals the top card of their library. If a land
// card is revealed this way, that player puts that card into their hand. Otherwise,
// that player puts a +1/+1 counter on the exploring creature, then chooses to either
// put the revealed card back on top of their library or into their graveyard.
type ExploreEffect struct {
	description       string
	exploringCreature uuid.UUID // The creature that is exploring
}

// NewExploreEffect creates an Explore effect
func NewExploreEffect(exploringCreature uuid.UUID) *ExploreEffect {
	return &ExploreEffect{
		description:       "Explore",
		exploringCreature: exploringCreature,
	}
}

// Apply implements the Effect interface
func (e *ExploreEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Reveal top card of controller's library
	// 2. If land card:
	//    a. Put into hand
	// 3. If not land:
	//    a. Put +1/+1 counter on exploring creature
	//    b. Controller chooses: top of library or graveyard
	//
	// Example: Jadelight Ranger "{1}{G}{G}, Creature - Merfolk Scout 2/1, When ~ enters,
	// it explores, then it explores again"

	return fmt.Errorf("explore effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *ExploreEffect) GetDescription() string {
	return e.description + ". (Reveal the top card of your library. Put that card into your hand if it's a land. Otherwise, put a +1/+1 counter on this creature, then put the card back or put it into your graveyard.)"
}

// ===== Connive (Rule 701.47) =====

// ConniveEffect represents the Connive keyword action
// Rule 701.47a: To connive, a player draws a card, then discards a card. If a nonland
// card was discarded this way, put a +1/+1 counter on the conniving creature.
type ConniveEffect struct {
	description       string
	connivingCreature uuid.UUID // The creature that is conniving
	conniveCount      int       // Number of times to connive (usually 1)
}

// NewConniveEffect creates a Connive effect
func NewConniveEffect(connivingCreature uuid.UUID, count int) *ConniveEffect {
	return &ConniveEffect{
		description:       "Connive",
		connivingCreature: connivingCreature,
		conniveCount:      count,
	}
}

// Apply implements the Effect interface
func (e *ConniveEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// For each connive:
	// 1. Draw a card
	// 2. Discard a card
	// 3. If discarded card was nonland, put +1/+1 counter on conniving creature
	//
	// Example: Raffine's Informant "{1}{W}, Creature - Human Wizard 1/1, When ~ enters,
	// it connives"

	return fmt.Errorf("connive effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *ConniveEffect) GetDescription() string {
	if e.conniveCount == 1 {
		return e.description + ". (Draw a card, then discard a card. If you discarded a nonland card, put a +1/+1 counter on this creature.)"
	}
	return fmt.Sprintf("%s %d times", e.description, e.conniveCount)
}

// ===== Amass (Rule 701.44) =====

// AmassEffect represents the Amass keyword action
// Rule 701.44a: To amass N, if you don't control an Army creature, create a 0/0 black
// Zombie Army creature token. Then you choose an Army creature you control and put N
// +1/+1 counters on it.
type AmassEffect struct {
	description string
	count       int
	armySubtype string // "Zombies" or "Orcs" (War of the Spark vs. Lord of the Rings)
}

// NewAmassEffect creates an Amass effect
func NewAmassEffect(count int, armySubtype string) *AmassEffect {
	if armySubtype == "" {
		armySubtype = "Zombies" // Default to Zombie Army
	}
	return &AmassEffect{
		description: "Amass",
		count:       count,
		armySubtype: armySubtype,
	}
}

// Apply implements the Effect interface
func (e *AmassEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Check if controller controls an Army creature
	// 2. If not, create 0/0 black [subtype] Army token
	// 3. Choose an Army creature
	// 4. Put N +1/+1 counters on it
	//
	// Example: Dreadhorde Invasion "{1}{B}, Enchantment, At the beginning of your upkeep,
	// you lose 1 life and amass Zombies 1"

	return fmt.Errorf("amass effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *AmassEffect) GetDescription() string {
	return fmt.Sprintf("%s %s %d. (Put %d +1/+1 counters on an Army you control. It's also a %s. If you don't control an Army, create a 0/0 black %s Army creature token first.)",
		e.description, e.armySubtype, e.count, e.count, e.armySubtype, e.armySubtype)
}

// ===== Adapt (Rule 702.142) =====

// AdaptEffect represents the Adapt keyword action
// Rule 702.142a: Adapt represents an activated ability. "Adapt N" means "If this
// permanent has no +1/+1 counters on it, put N +1/+1 counters on it."
type AdaptEffect struct {
	description string
	count       int
	permanent   uuid.UUID
}

// NewAdaptEffect creates an Adapt effect
func NewAdaptEffect(permanent uuid.UUID, count int) *AdaptEffect {
	return &AdaptEffect{
		description: "Adapt",
		count:       count,
		permanent:   permanent,
	}
}

// Apply implements the Effect interface
func (e *AdaptEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Check if permanent has any +1/+1 counters
	// 2. If not, put N +1/+1 counters on it
	// 3. Otherwise, do nothing
	//
	// Example: Aeromunculus "{1}{G}{U}, Creature - Homunculus Mutant 2/3, {2}{G}{U}: Adapt 1"

	// Get counter game context
	counterGame, ok := game.(CounterGameContext)
	if !ok {
		return fmt.Errorf("game context does not support counter operations")
	}

	// Get permanent
	permanent, err := counterGame.GetPermanent(e.permanent)
	if err != nil {
		return fmt.Errorf("failed to get permanent: %w", err)
	}

	// TODO: Check if permanent has +1/+1 counters
	// For now, just add counters (needs permanent counter tracking)
	counter := counters.NewCounter("+1/+1", e.count)
	if err := counterGame.AddCountersToPermanent(permanent, counter); err != nil {
		return fmt.Errorf("failed to add counters: %w", err)
	}

	return nil
}

// GetDescription returns the effect description
func (e *AdaptEffect) GetDescription() string {
	return fmt.Sprintf("%s %d. (If this creature has no +1/+1 counters on it, put %d +1/+1 counters on it.)", e.description, e.count, e.count)
}

// ===== Monstrosity (Rule 701.31) =====

// MonstrosityEffect represents the Monstrosity keyword action
// Rule 701.31a: Monstrosity represents two linked abilities. "Monstrosity N" means
// "If this permanent isn't monstrous, put N +1/+1 counters on it and it becomes monstrous."
type MonstrosityEffect struct {
	description string
	count       int
	permanent   uuid.UUID
}

// NewMonstrosityEffect creates a Monstrosity effect
func NewMonstrosityEffect(permanent uuid.UUID, count int) *MonstrosityEffect {
	return &MonstrosityEffect{
		description: "Monstrosity",
		count:       count,
		permanent:   permanent,
	}
}

// Apply implements the Effect interface
func (e *MonstrosityEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Check if permanent is monstrous
	// 2. If not:
	//    a. Put N +1/+1 counters on it
	//    b. Mark it as monstrous (permanent status)
	// 3. Otherwise, do nothing
	//
	// Example: Arbor Colossus "{2}{G}{G}{G}, Creature - Giant 6/6, {3}{G}{G}{G}: Monstrosity 3.
	// When ~ becomes monstrous, destroy target creature with flying an opponent controls"

	return fmt.Errorf("monstrosity effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *MonstrosityEffect) GetDescription() string {
	return fmt.Sprintf("%s %d. (If this creature isn't monstrous, put %d +1/+1 counters on it and it becomes monstrous.)", e.description, e.count, e.count)
}

// ===== Bolster (Rule 701.33) =====

// BolsterEffect represents the Bolster keyword action
// Rule 701.33a: To bolster N, choose a creature you control with the least toughness or
// tied for least toughness among creatures you control. Put N +1/+1 counters on it.
type BolsterEffect struct {
	description string
	count       int
}

// NewBolsterEffect creates a Bolster effect
func NewBolsterEffect(count int) *BolsterEffect {
	return &BolsterEffect{
		description: "Bolster",
		count:       count,
	}
}

// Apply implements the Effect interface
func (e *BolsterEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Find all creatures controller controls
	// 2. Find creature(s) with lowest toughness
	// 3. If multiple tied, controller chooses one
	// 4. Put N +1/+1 counters on chosen creature
	//
	// Example: Dragon Bell Monk "{3}{W}, Creature - Human Monk 2/2, Whenever ~ attacks,
	// bolster 1"

	return fmt.Errorf("bolster effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *BolsterEffect) GetDescription() string {
	return fmt.Sprintf("%s %d. (Choose a creature with the least toughness among creatures you control and put %d +1/+1 counters on it.)", e.description, e.count, e.count)
}

// ===== Support (Rule 701.34) =====

// SupportEffect represents the Support keyword action
// Rule 701.34a: To support N, put a +1/+1 counter on each of up to N target creatures.
type SupportEffect struct {
	description string
	count       int
}

// NewSupportEffect creates a Support effect
func NewSupportEffect(count int) *SupportEffect {
	return &SupportEffect{
		description: "Support",
		count:       count,
	}
}

// Apply implements the Effect interface
func (e *SupportEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. For each target (up to N targets)
	// 2. Put one +1/+1 counter on it
	//
	// Example: Relief Captain "{2}{W}, Creature - Kor Knight Ally 3/2, When ~ enters,
	// support 3"

	counterGame, ok := game.(CounterGameContext)
	if !ok {
		return fmt.Errorf("game context does not support counter operations")
	}

	// Add one +1/+1 counter to each target
	counter := counters.NewCounter("+1/+1", 1)
	for _, targetID := range targets {
		permanent, err := counterGame.GetPermanent(targetID)
		if err != nil {
			continue // Skip invalid targets
		}
		if err := counterGame.AddCountersToPermanent(permanent, counter); err != nil {
			return fmt.Errorf("failed to add counter to target: %w", err)
		}
	}

	return nil
}

// GetDescription returns the effect description
func (e *SupportEffect) GetDescription() string {
	return fmt.Sprintf("%s %d. (Put a +1/+1 counter on each of up to %d other target creatures.)", e.description, e.count, e.count)
}

// ===== Goad (Rule 701.38) =====

// GoadEffect represents the Goad keyword action
// Rule 701.38a: To goad a creature means to make it goaded until a player's next turn.
// A goaded creature attacks each combat if able and attacks a player other than the
// player who goaded it if able.
type GoadEffect struct {
	description    string
	goadingPlayer  uuid.UUID
	goadedCreature uuid.UUID
}

// NewGoadEffect creates a Goad effect
func NewGoadEffect(goadingPlayer, goadedCreature uuid.UUID) *GoadEffect {
	return &GoadEffect{
		description:    "Goad",
		goadingPlayer:  goadingPlayer,
		goadedCreature: goadedCreature,
	}
}

// Apply implements the Effect interface
func (e *GoadEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Mark creature as goaded
	// 2. Track who goaded it
	// 3. Create continuous effect:
	//    a. Must attack each combat if able
	//    b. Must attack player other than goading player if able
	// 4. Effect lasts until goading player's next turn
	//
	// Example: Karazikar, the Eye Tyrant "{3}{B}{R}, Legendary Creature - Beholder 5/5,
	// Whenever you attack a player, gain control of target creature that player controls
	// until end of turn. Untap it. It gains haste until end of turn. Goad it."

	return fmt.Errorf("goad effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *GoadEffect) GetDescription() string {
	return e.description + ". (Until your next turn, that creature attacks each combat if able and attacks a player other than you if able.)"
}

// ===== Learn (Rule 701.45) =====

// LearnEffect represents the Learn keyword action
// Rule 701.45a: To learn, a player may discard a card. If they do, they draw a card.
// If they don't discard a card, they may reveal a Lesson card from outside the game
// and put it into their hand.
type LearnEffect struct {
	description string
}

// NewLearnEffect creates a Learn effect
func NewLearnEffect() *LearnEffect {
	return &LearnEffect{
		description: "Learn",
	}
}

// Apply implements the Effect interface
func (e *LearnEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Player chooses:
	//    Option A: Discard a card, then draw a card
	//    Option B: Reveal a Lesson card from sideboard and put it into hand
	// 2. Execute chosen option
	//
	// Example: Illuminate History "{2}{R}, Instant - Lesson, Discard any number of cards,
	// then draw that many cards. Then if seven or more cards are in your graveyard,
	// create a 3/2 red and white Spirit creature token"

	return fmt.Errorf("learn effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *LearnEffect) GetDescription() string {
	return e.description + ". (You may discard a card. If you do, draw a card. If you don't, you may reveal a Lesson card you own from outside the game and put it into your hand.)"
}

// ===== Fateseal (Like Scry for opponent) =====

// FatesealEffect represents the Fateseal action (like Scry but for opponent)
// Not an official keyword action, but appears on Jace, the Mind Sculptor
type FatesealEffect struct {
	description string
	count       int
	opponent    uuid.UUID
}

// NewFatesealEffect creates a Fateseal effect
func NewFatesealEffect(count int, opponent uuid.UUID) *FatesealEffect {
	return &FatesealEffect{
		description: "Fateseal",
		count:       count,
		opponent:    opponent,
	}
}

// Apply implements the Effect interface
func (e *FatesealEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Look at top N cards of opponent's library
	// 2. Controller (not opponent) chooses:
	//    - Any number go to bottom of library in any order
	//    - Rest stay on top in any order
	//
	// Like Scry but controller manipulates opponent's library
	// Example: Jace, the Mind Sculptor "+0: Fateseal 1"

	return fmt.Errorf("fateseal effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *FatesealEffect) GetDescription() string {
	return fmt.Sprintf("%s %d. (Look at the top %d cards of an opponent's library, then put any number of them on the bottom of that library and the rest on top in any order.)", e.description, e.count, e.count)
}

// ===== Clash (Rule 702.63) =====

// ClashEffect represents the Clash keyword action
// Rule 702.63a: To clash, a player reveals the top card of their library. That player
// may then put that card on the bottom of their library. A player wins a clash if the
// card they revealed had a higher mana value than all other cards revealed in that clash.
type ClashEffect struct {
	description string
	player      uuid.UUID
}

// NewClashEffect creates a Clash effect
func NewClashEffect(player uuid.UUID) *ClashEffect {
	return &ClashEffect{
		description: "Clash",
		player:      player,
	}
}

// Apply implements the Effect interface
func (e *ClashEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	// Implementation steps:
	// 1. Reveal top card of each clashing player's library
	// 2. Compare mana values
	// 3. Player with highest mana value wins
	// 4. Each player may put their revealed card on bottom of library
	//
	// Example: Coordinated Barrage "{W}, Instant, Choose up to three target creatures
	// you control. Clash with an opponent. If you win, tap those creatures"

	return fmt.Errorf("clash effect not yet fully implemented in game context")
}

// GetDescription returns the effect description
func (e *ClashEffect) GetDescription() string {
	return e.description + ". (Each clashing player reveals the top card of their library, then puts that card on the top or bottom. A player wins if their card had a higher mana value.)"
}

// ===== Common Keyword Action Examples =====

// Example cards using these actions:
// - Explore: Jadelight Ranger, Path of Discovery, Journey to Eternity
// - Connive: Raffine's Informant, Ledger Shredder, Obscura Interceptor
// - Surveil: Enhanced Surveillance, Dimir Informant, Disinformation Campaign
// - Amass: Dreadhorde Invasion, Lazotep Reaver, Gleaming Overseer
// - Adapt: Aeromunculus, Sauroform Hybrid, Skatewing Spy
// - Monstrosity: Arbor Colossus, Colossus of Akros, Fleecemane Lion
// - Bolster: Dragon Bell Monk, Cached Defenses, Scale Blessing
// - Support: Relief Captain, Spawnbinder Mage, Expedition Raptor
// - Goad: Karazikar, Kardur's Vicious Return, Geode Rager
// - Learn: Lesson cards, Study Break, Eager First-Year
// - Fateseal: Jace, the Mind Sculptor
// - Clash: Coordinated Barrage, Bog-Strider Ash, Shields of Velis Vel
