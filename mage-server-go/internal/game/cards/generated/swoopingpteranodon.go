package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Swooping Pteranodon", NewSwoopingPteranodon)
}

// NewSwoopingPteranodon creates a Swooping Pteranodon
// {3}{R}{W} - CREATURE
// Flying, Haste
func NewSwoopingPteranodon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Swooping Pteranodon")
	card.ManaCost = "{3}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	return card, nil
}
