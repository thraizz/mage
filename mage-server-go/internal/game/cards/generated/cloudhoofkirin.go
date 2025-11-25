package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cloudhoof Kirin", NewCloudhoofKirin)
}

// NewCloudhoofKirin creates a Cloudhoof Kirin
// {3}{U}{U} - CREATURE
// Flying
func NewCloudhoofKirin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cloudhoof Kirin")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KIRIN", "SPIRIT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SpellCastControllerTriggeredAbility
	//   - Effect: CloudhoofKirinEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	// card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability1)
	return card, nil
}
