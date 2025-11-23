package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bassara Tower Archer", NewBassaraTowerArcher)
}

// NewBassaraTowerArcher creates a Bassara Tower Archer
// {G}{G} - CREATURE
// Hexproof, Reach
func NewBassaraTowerArcher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bassara Tower Archer")
	card.ManaCost = "{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability1)
	return card, nil
}
