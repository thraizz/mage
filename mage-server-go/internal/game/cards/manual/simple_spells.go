package manual

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

// Register simple spells on package import
func init() {
	cards.Register("Lightning Bolt", NewLightningBolt)
	cards.Register("Shock", NewShock)
	cards.Register("Giant Growth", NewGiantGrowth)
	cards.Register("Divination", NewDivination)
	cards.Register("Murder", NewMurder)
	cards.Register("Counterspell", NewCounterspell)
}

// NewLightningBolt creates a Lightning Bolt
// {R} - Instant - Deal 3 damage to any target
func NewLightningBolt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lightning Bolt")
	card.ManaCost = "{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = "Lightning Bolt deals 3 damage to any target."

	// Build spell ability: Deal 3 damage to any target
	ability, err := abilities.NewSpellAbilityBuilder(card.ID, "{R}").
		AddEffect(abilities.NewDamageEffect(3)).
		AddTarget(abilities.NewAnyTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}

	card.AddAbility(ability)
	return card, nil
}

// NewShock creates a Shock
// {R} - Instant - Deal 2 damage to any target
func NewShock(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shock")
	card.ManaCost = "{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = "Shock deals 2 damage to any target."

	// Build spell ability: Deal 2 damage to any target
	ability, err := abilities.NewSpellAbilityBuilder(card.ID, "{R}").
		AddEffect(abilities.NewDamageEffect(2)).
		AddTarget(abilities.NewAnyTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}

	card.AddAbility(ability)
	return card, nil
}

// NewGiantGrowth creates a Giant Growth
// {G} - Instant - Target creature gets +3/+3 until end of turn
func NewGiantGrowth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Giant Growth")
	card.ManaCost = "{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = "Target creature gets +3/+3 until end of turn."

	// Build spell ability: Target creature gets +3/+3
	ability, err := abilities.NewSpellAbilityBuilder(card.ID, "{G}").
		AddEffect(abilities.NewBoostEffect(3, 3)).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}

	card.AddAbility(ability)
	return card, nil
}

// NewDivination creates a Divination
// {2}{U} - Sorcery - Draw two cards
func NewDivination(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Divination")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = "Draw two cards."

	// Build spell ability: Draw two cards
	ability, err := abilities.NewSpellAbilityBuilder(card.ID, "{2}{U}").
		AddEffect(abilities.NewDrawCardsEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}

	card.AddAbility(ability)
	return card, nil
}

// NewMurder creates a Murder
// {1}{B}{B} - Instant - Destroy target creature
func NewMurder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Murder")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = "Destroy target creature."

	// Build spell ability: Destroy target creature
	ability, err := abilities.NewSpellAbilityBuilder(card.ID, "{1}{B}{B}").
		AddEffect(abilities.NewDestroyEffect()).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}

	card.AddAbility(ability)
	return card, nil
}

// NewCounterspell creates a Counterspell
// {U}{U} - Instant - Counter target spell
func NewCounterspell(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Counterspell")
	card.ManaCost = "{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M10"
	card.Rarity = "common"
	card.RulesText = "Counter target spell."

	// Build spell ability: Counter target spell
	ability, err := abilities.NewSpellAbilityBuilder(card.ID, "{U}{U}").
		AddEffect(abilities.NewCounterSpellEffect()).
		AddTarget(abilities.NewSpellTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}

	card.AddAbility(ability)
	return card, nil
}
