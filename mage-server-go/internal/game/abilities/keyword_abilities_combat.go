package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// This file implements combat-related keyword abilities
// Rule 702: Keyword Abilities

// ===== Flying (Rule 702.9) =====

// FlyingAbility represents the Flying keyword
// Rule 702.9a: Flying is an evasion ability
// Rule 702.9b: A creature with flying can't be blocked except by creatures with
// flying and/or reach
type FlyingAbility struct {
	baseAbility
}

// NewFlyingAbility creates a Flying keyword ability
func NewFlyingAbility(source uuid.UUID) *FlyingAbility {
	return &FlyingAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *FlyingAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *FlyingAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true // Always active
}

func (a *FlyingAbility) Resolve(ctx context.Context, game GameContext) error {
	// Flying doesn't resolve - it's checked during combat (declare blockers step)
	return nil
}

func (a *FlyingAbility) String() string {
	return "Flying"
}

// ===== First Strike (Rule 702.7) =====

// FirstStrikeAbility represents the First Strike keyword
// Rule 702.7a: First strike is a static ability that modifies combat damage timing
// Rule 702.7b: Creatures with first strike deal combat damage before creatures without
type FirstStrikeAbility struct {
	baseAbility
}

// NewFirstStrikeAbility creates a First Strike keyword ability
func NewFirstStrikeAbility(source uuid.UUID) *FirstStrikeAbility {
	return &FirstStrikeAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *FirstStrikeAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *FirstStrikeAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *FirstStrikeAbility) Resolve(ctx context.Context, game GameContext) error {
	// First strike doesn't resolve - it modifies combat damage step timing
	// Creates an additional first strike damage step if any creature has first strike
	return nil
}

func (a *FirstStrikeAbility) String() string {
	return "First strike"
}

// ===== Double Strike (Rule 702.4) =====

// DoubleStrikeAbility represents the Double Strike keyword
// Rule 702.4a: Double strike is a static ability that modifies combat damage timing
// Rule 702.4b: Creatures with double strike deal combat damage in both the first
// strike and regular combat damage steps
type DoubleStrikeAbility struct {
	baseAbility
}

// NewDoubleStrikeAbility creates a Double Strike keyword ability
func NewDoubleStrikeAbility(source uuid.UUID) *DoubleStrikeAbility {
	return &DoubleStrikeAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *DoubleStrikeAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *DoubleStrikeAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *DoubleStrikeAbility) Resolve(ctx context.Context, game GameContext) error {
	// Double strike doesn't resolve - creature deals damage twice
	return nil
}

func (a *DoubleStrikeAbility) String() string {
	return "Double strike"
}

// ===== Deathtouch (Rule 702.2) =====

// DeathtouchAbility represents the Deathtouch keyword
// Rule 702.2a: Deathtouch is a static ability
// Rule 702.2b: Any nonzero amount of combat damage assigned to a creature by a source
// with deathtouch is considered to be lethal damage
type DeathtouchAbility struct {
	baseAbility
}

// NewDeathtouchAbility creates a Deathtouch keyword ability
func NewDeathtouchAbility(source uuid.UUID) *DeathtouchAbility {
	return &DeathtouchAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *DeathtouchAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *DeathtouchAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *DeathtouchAbility) Resolve(ctx context.Context, game GameContext) error {
	// Deathtouch doesn't resolve - it's checked during damage assignment/resolution
	// Creates a state-based action that destroys creatures dealt damage by deathtouch source
	return nil
}

func (a *DeathtouchAbility) String() string {
	return "Deathtouch"
}

// ===== Lifelink (Rule 702.15) =====

// LifelinkAbility represents the Lifelink keyword
// Rule 702.15a: Lifelink is a static ability
// Rule 702.15b: Damage dealt by a source with lifelink causes that source's controller
// to gain that much life (in addition to any other results that damage causes)
type LifelinkAbility struct {
	baseAbility
}

// NewLifelinkAbility creates a Lifelink keyword ability
func NewLifelinkAbility(source uuid.UUID) *LifelinkAbility {
	return &LifelinkAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *LifelinkAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *LifelinkAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *LifelinkAbility) Resolve(ctx context.Context, game GameContext) error {
	// Lifelink doesn't resolve - it's triggered when source deals damage
	// Controller gains life equal to damage dealt
	return nil
}

func (a *LifelinkAbility) String() string {
	return "Lifelink"
}

// ===== Vigilance (Rule 702.20) =====

// VigilanceAbility represents the Vigilance keyword
// Rule 702.20a: Vigilance is a static ability that modifies how a creature attacks
// Rule 702.20b: Attacking doesn't cause a creature with vigilance to tap
type VigilanceAbility struct {
	baseAbility
}

// NewVigilanceAbility creates a Vigilance keyword ability
func NewVigilanceAbility(source uuid.UUID) *VigilanceAbility {
	return &VigilanceAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *VigilanceAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *VigilanceAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *VigilanceAbility) Resolve(ctx context.Context, game GameContext) error {
	// Vigilance doesn't resolve - it's checked during declare attackers step
	// Creature doesn't tap when attacking
	return nil
}

func (a *VigilanceAbility) String() string {
	return "Vigilance"
}

// ===== Trample (Rule 702.19) =====

// TrampleAbility represents the Trample keyword
// Rule 702.19a: Trample is a static ability that modifies damage assignment
// Rule 702.19b: If attacking creature with trample is blocked, excess damage may
// be assigned to the player or planeswalker it's attacking
type TrampleAbility struct {
	baseAbility
}

// NewTrampleAbility creates a Trample keyword ability
func NewTrampleAbility(source uuid.UUID) *TrampleAbility {
	return &TrampleAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *TrampleAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *TrampleAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *TrampleAbility) Resolve(ctx context.Context, game GameContext) error {
	// Trample doesn't resolve - it modifies combat damage assignment
	// Excess damage beyond lethal to blockers is assigned to defending player
	return nil
}

func (a *TrampleAbility) String() string {
	return "Trample"
}

// ===== Reach (Rule 702.17) =====

// ReachAbility represents the Reach keyword
// Rule 702.17a: Reach is a static ability
// Rule 702.17b: A creature with flying can be blocked by creatures with reach
type ReachAbility struct {
	baseAbility
}

// NewReachAbility creates a Reach keyword ability
func NewReachAbility(source uuid.UUID) *ReachAbility {
	return &ReachAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *ReachAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *ReachAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *ReachAbility) Resolve(ctx context.Context, game GameContext) error {
	// Reach doesn't resolve - it's checked during declare blockers step
	// Allows blocking creatures with flying
	return nil
}

func (a *ReachAbility) String() string {
	return "Reach"
}

// ===== Menace (Rule 702.111) =====

// MenaceAbility represents the Menace keyword
// Rule 702.111a: Menace is an evasion ability
// Rule 702.111b: A creature with menace can't be blocked except by two or more creatures
type MenaceAbility struct {
	baseAbility
}

// NewMenaceAbility creates a Menace keyword ability
func NewMenaceAbility(source uuid.UUID) *MenaceAbility {
	return &MenaceAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *MenaceAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *MenaceAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *MenaceAbility) Resolve(ctx context.Context, game GameContext) error {
	// Menace doesn't resolve - it's enforced during declare blockers step
	// Requires at least 2 creatures to block
	return nil
}

func (a *MenaceAbility) String() string {
	return "Menace"
}

// ===== Defender (Rule 702.3) =====

// DefenderAbility represents the Defender keyword
// Rule 702.3a: Defender is a static ability
// Rule 702.3b: A creature with defender can't attack
type DefenderAbility struct {
	baseAbility
}

// NewDefenderAbility creates a Defender keyword ability
func NewDefenderAbility(source uuid.UUID) *DefenderAbility {
	return &DefenderAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *DefenderAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *DefenderAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *DefenderAbility) Resolve(ctx context.Context, game GameContext) error {
	// Defender doesn't resolve - it's checked during declare attackers step
	// Prevents creature from being declared as an attacker
	return nil
}

func (a *DefenderAbility) String() string {
	return "Defender"
}

// ===== Hexproof (Rule 702.11) =====

// HexproofAbility represents the Hexproof keyword
// Rule 702.11a: Hexproof is a static ability
// Rule 702.11b: A permanent or player with hexproof can't be the target of spells
// or abilities opponents control
type HexproofAbility struct {
	baseAbility
}

// NewHexproofAbility creates a Hexproof keyword ability
func NewHexproofAbility(source uuid.UUID) *HexproofAbility {
	return &HexproofAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *HexproofAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *HexproofAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *HexproofAbility) Resolve(ctx context.Context, game GameContext) error {
	// Hexproof doesn't resolve - it's checked during targeting
	// Prevents opponents from targeting this permanent
	return nil
}

func (a *HexproofAbility) String() string {
	return "Hexproof"
}

// ===== Shroud (Rule 702.18) =====

// ShroudAbility represents the Shroud keyword
// Rule 702.18a: Shroud is a static ability
// Rule 702.18b: A permanent with shroud can't be the target of spells or abilities
type ShroudAbility struct {
	baseAbility
}

// NewShroudAbility creates a Shroud keyword ability
func NewShroudAbility(source uuid.UUID) *ShroudAbility {
	return &ShroudAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
	}
}

func (a *ShroudAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *ShroudAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *ShroudAbility) Resolve(ctx context.Context, game GameContext) error {
	// Shroud doesn't resolve - it's checked during targeting
	// Prevents anyone (including controller) from targeting this permanent
	return nil
}

func (a *ShroudAbility) String() string {
	return "Shroud"
}

// ===== Protection (Rule 702.16) =====

// ProtectionAbility represents the Protection keyword
// Rule 702.16a: Protection is a static ability with a quality
// Rule 702.16b: Protection from [quality] means:
// - Damage: All damage that would be dealt by sources with [quality] is prevented
// - Enchant/Equip: This can't be enchanted/equipped by permanents with [quality]
// - Block: This can't be blocked by creatures with [quality]
// - Target: This can't be targeted by spells or abilities with [quality]
type ProtectionAbility struct {
	baseAbility
	fromQuality string // e.g., "red", "artifacts", "creatures", "everything"
}

// NewProtectionAbility creates a Protection keyword ability
func NewProtectionAbility(source uuid.UUID, fromQuality string) *ProtectionAbility {
	return &ProtectionAbility{
		baseAbility: baseAbility{
			id:       uuid.New(),
			sourceID: source,
		},
		fromQuality: fromQuality,
	}
}

func (a *ProtectionAbility) GetType() AbilityType {
	return AbilityTypeStatic
}

func (a *ProtectionAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

func (a *ProtectionAbility) Resolve(ctx context.Context, game GameContext) error {
	// Protection doesn't resolve - it creates multiple prevention/restriction effects
	// DEBT: Damage, Enchant/Equip, Block, Target
	return nil
}

func (a *ProtectionAbility) GetFromQuality() string {
	return a.fromQuality
}

func (a *ProtectionAbility) String() string {
	return fmt.Sprintf("Protection from %s", a.fromQuality)
}

// ===== Helper Functions =====

// HasFlying checks if a permanent has Flying
func HasFlying(permanentID uuid.UUID, game GameContext) bool {
	// This would check the permanent's abilities for Flying
	return false // Placeholder
}

// HasFirstStrike checks if a permanent has First Strike
func HasFirstStrike(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasDoubleStrike checks if a permanent has Double Strike
func HasDoubleStrike(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasDeathtouch checks if a permanent has Deathtouch
func HasDeathtouch(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasLifelink checks if a permanent has Lifelink
func HasLifelink(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasVigilance checks if a permanent has Vigilance
func HasVigilance(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasTrample checks if a permanent has Trample
func HasTrample(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasReach checks if a permanent has Reach
func HasReach(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasMenace checks if a permanent has Menace
func HasMenace(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasDefender checks if a permanent has Defender
func HasDefender(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasHexproof checks if a permanent has Hexproof
func HasHexproof(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasShroud checks if a permanent has Shroud
func HasShroud(permanentID uuid.UUID, game GameContext) bool {
	return false // Placeholder
}

// HasProtectionFrom checks if a permanent has Protection from a specific quality
func HasProtectionFrom(permanentID uuid.UUID, quality string, game GameContext) bool {
	return false // Placeholder
}

// CanBlock checks if a creature can block an attacker
func CanBlock(blockerID, attackerID uuid.UUID, game GameContext) bool {
	// Check:
	// 1. Blocker doesn't have defender or is tapped
	// 2. Attacker doesn't have flying (unless blocker has flying/reach)
	// 3. Attacker doesn't have protection from blocker's qualities
	// 4. Blocker doesn't have protection that prevents blocking
	return false // Placeholder
}

// CanBeBlocked checks if an attacker can be legally blocked
func CanBeBlocked(attackerID uuid.UUID, blockers []uuid.UUID, game GameContext) bool {
	// Check:
	// 1. Attacker doesn't have "can't be blocked"
	// 2. If menace, at least 2 blockers
	// 3. If flying, blockers have flying/reach
	// 4. No protection preventing blocks
	return false // Placeholder
}

// ===== Common Keyword Ability Examples =====

// Example cards with these keywords:
// - Flying: Storm Crow, Serra Angel, Faerie Miscreant
// - First Strike: Elite Vanguard, First Strike creatures
// - Double Strike: Boros Reckoner, Double Strike creatures
// - Deathtouch: Typhoid Rats, Gifted Aetherborn
// - Lifelink: Vampire Nighthawk, Ajani's Pridemate triggers
// - Vigilance: Serra Angel, Sunblast Angel
// - Trample: Colossal Dreadmaw, Ghalta Primal Hunger
// - Reach: Giant Spider, Plummet targets
// - Menace: Goblin Piker variants
// - Defender: Wall of Omens, Overgrown Battlement
// - Hexproof: Invisible Stalker, Slippery Bogle
// - Shroud: Troll Ascetic (old), Progenitus has shroud-like ability
// - Protection: Circle of Protection series, Mother of Runes
