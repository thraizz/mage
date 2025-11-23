package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ebony Rhino", NewEbonyRhino)
}

// NewEbonyRhino creates a Ebony Rhino
// {7} - ARTIFACT CREATURE
// Trample
func NewEbonyRhino(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ebony Rhino")
	card.ManaCost = "{7}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"RHINO"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
