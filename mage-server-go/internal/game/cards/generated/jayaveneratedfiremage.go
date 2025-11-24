package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jaya Venerated Firemage", NewJayaVeneratedFiremage)
}

// NewJayaVeneratedFiremage creates a Jaya Venerated Firemage
// {4}{R} - PLANESWALKER
func NewJayaVeneratedFiremage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jaya Venerated Firemage")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"JAYA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}