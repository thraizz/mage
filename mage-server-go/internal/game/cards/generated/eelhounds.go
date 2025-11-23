package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eel Hounds", NewEelHounds)
}

// NewEelHounds creates a Eel Hounds
// {3}{G} - CREATURE
// Trample
func NewEelHounds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eel Hounds")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FISH", "DOG"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(2, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
