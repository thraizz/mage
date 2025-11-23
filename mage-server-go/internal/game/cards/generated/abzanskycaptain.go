package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abzan Skycaptain", NewAbzanSkycaptain)
}

// NewAbzanSkycaptain creates a Abzan Skycaptain
// {3}{W} - CREATURE
// Flying
func NewAbzanSkycaptain(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abzan Skycaptain")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD", "SOLDIER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
