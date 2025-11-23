package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Warmonger Hellkite", NewWarmongerHellkite)
}

// NewWarmongerHellkite creates a Warmonger Hellkite
// {4}{R}{R} - CREATURE
// Flying
func NewWarmongerHellkite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Warmonger Hellkite")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
