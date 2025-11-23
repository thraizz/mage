package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nightveil Specter", NewNightveilSpecter)
}

// NewNightveilSpecter creates a Nightveil Specter
// {U/B}{U/B}{U/B} - CREATURE
// Flying
func NewNightveilSpecter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nightveil Specter")
	card.ManaCost = "{U/B}{U/B}{U/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPECTER"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
