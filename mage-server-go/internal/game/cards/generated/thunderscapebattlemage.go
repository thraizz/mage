package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thunderscape Battlemage", NewThunderscapeBattlemage)
}

// NewThunderscapeBattlemage creates a Thunderscape Battlemage
// {2}{R} - CREATURE
func NewThunderscapeBattlemage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thunderscape Battlemage")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: DiscardTargetEffect(2)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewEnchantmentTargetFilter())
	// card.AddAbility(ability0)
	ability1 := abilities.NewKickerAbility(card.ID, "{1}{B}")
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(2)
	// card.AddAbility(ability2)
	return card, nil
}
