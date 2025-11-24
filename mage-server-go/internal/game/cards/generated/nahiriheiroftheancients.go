package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nahiri Heir Of The Ancients", NewNahiriHeirOfTheAncients)
}

// NewNahiriHeirOfTheAncients creates a Nahiri Heir Of The Ancients
// {2}{R}{W} - PLANESWALKER
func NewNahiriHeirOfTheAncients(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nahiri Heir Of The Ancients")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NAHIRI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(xValue)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(6, 1, filter, PutCards.HAND, PutCards.BOTTOM_RANDO...)
	// card.AddAbility(ability1)
	return card, nil
}