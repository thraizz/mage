package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hellcarver Demon", NewHellcarverDemon)
}

// NewHellcarverDemon creates a Hellcarver Demon
// {3}{B}{B}{B} - CREATURE
// Flying
func NewHellcarverDemon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hellcarver Demon")
	card.ManaCost = "{3}{B}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
