package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kasmina Enigma Sage", NewKasminaEnigmaSage)
}

// NewKasminaEnigmaSage creates a Kasmina Enigma Sage
// {1}{G}{U} - PLANESWALKER
func NewKasminaEnigmaSage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kasmina Enigma Sage")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KASMINA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewScryEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}