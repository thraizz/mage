package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Black Chocobo", NewBlackChocobo)
}

// NewBlackChocobo creates a Black Chocobo
//   - CREATURE
func NewBlackChocobo(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Black Chocobo")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSearchLibraryPutInPlayEffect(abilities.NewTargetRequirement(0, 1, abilities.NewLandTargetFilter()), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 0, filter, false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
