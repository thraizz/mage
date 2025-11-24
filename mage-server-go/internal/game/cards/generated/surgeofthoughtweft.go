package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Surge Of Thoughtweft", NewSurgeOfThoughtweft)
}

// NewSurgeOfThoughtweft creates a Surge Of Thoughtweft
// {1}{W} - KINDRED INSTANT
func NewSurgeOfThoughtweft(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Surge Of Thoughtweft")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"KITHKIN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1)).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}