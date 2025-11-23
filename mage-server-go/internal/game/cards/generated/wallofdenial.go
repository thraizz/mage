package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wall Of Denial", NewWallOfDenial)
}

// NewWallOfDenial creates a Wall Of Denial
// {1}{W}{U} - CREATURE
// Defender, Flying
func NewWallOfDenial(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wall Of Denial")
	card.ManaCost = "{1}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WALL"}
	card.Power = "0"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability1)
	return card, nil
}
