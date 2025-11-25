package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abstruse Archaic", NewAbstruseArchaic)
}

// NewAbstruseArchaic creates a Abstruse Archaic
// {4} - CREATURE
// Vigilance
func NewAbstruseArchaic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abstruse Archaic")
	card.ManaCost = "{4}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - CopyTargetStackObjectEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
