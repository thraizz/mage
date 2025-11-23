package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Glimmerpoint Stag", NewGlimmerpointStag)
}

// NewGlimmerpointStag creates a Glimmerpoint Stag
// {2}{W}{W} - CREATURE
// Vigilance
func NewGlimmerpointStag(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Glimmerpoint Stag")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELK"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
