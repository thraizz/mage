package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mr Foxglove", NewMrFoxglove)
}

// NewMrFoxglove creates a Mr Foxglove
// {2}{G}{W}{U} - CREATURE
// Lifelink
func NewMrFoxglove(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mr Foxglove")
	card.ManaCost = "{2}{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FOX", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	return card, nil
}
