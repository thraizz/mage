package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sylvan Caryatid", NewSylvanCaryatid)
}

// NewSylvanCaryatid creates a Sylvan Caryatid
// {1}{G} - CREATURE
// Defender, Hexproof
func NewSylvanCaryatid(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sylvan Caryatid")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PLANT"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability1)
	return card, nil
}
