package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sygg River Cutthroat", NewSyggRiverCutthroat)
}

// NewSyggRiverCutthroat creates a Sygg River Cutthroat
// {U/B}{U/B} - CREATURE
func NewSyggRiverCutthroat(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sygg River Cutthroat")
	card.ManaCost = "{U/B}{U/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
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
