package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abundant Maw", NewAbundantMaw)
}

// NewAbundantMaw creates a Abundant Maw
// {8} - CREATURE
func NewAbundantMaw(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abundant Maw")
	card.ManaCost = "{8}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI", "LEECH"}
	card.Power = "6"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewLoseLifeEffect(3)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
