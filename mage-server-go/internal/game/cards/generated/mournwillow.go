package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mournwillow", NewMournwillow)
}

// NewMournwillow creates a Mournwillow
// {1}{B}{G} - CREATURE
// Haste
func NewMournwillow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mournwillow")
	card.ManaCost = "{1}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PLANT", "SKELETON"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
