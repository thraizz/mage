package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kadena Slinking Sorcerer", NewKadenaSlinkingSorcerer)
}

// NewKadenaSlinkingSorcerer creates a Kadena Slinking Sorcerer
// {1}{B}{G}{U} - CREATURE
func NewKadenaSlinkingSorcerer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kadena Slinking Sorcerer")
	card.ManaCost = "{1}{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
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
