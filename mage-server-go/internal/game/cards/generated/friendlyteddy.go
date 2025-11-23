package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Friendly Teddy", NewFriendlyTeddy)
}

// NewFriendlyTeddy creates a Friendly Teddy
// {2} - ARTIFACT CREATURE
func NewFriendlyTeddy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Friendly Teddy")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"BEAR", "TOY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
