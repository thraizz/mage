package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pyromancers Goggles", NewPyromancersGoggles)
}

// NewPyromancersGoggles creates a Pyromancers Goggles
// {5} - ARTIFACT
func NewPyromancersGoggles(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pyromancers Goggles")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability0)
	return card, nil
}
