package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jumbo Cactuar", NewJumboCactuar)
}

// NewJumboCactuar creates a Jumbo Cactuar
// {5}{G}{G} - CREATURE
func NewJumboCactuar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jumbo Cactuar")
	card.ManaCost = "{5}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PLANT"}
	card.Power = "1"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(9999, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
