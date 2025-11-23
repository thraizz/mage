package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Freejam Regent", NewFreejamRegent)
}

// NewFreejamRegent creates a Freejam Regent
// {4}{R}{R} - CREATURE
// Flying
func NewFreejamRegent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Freejam Regent")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
