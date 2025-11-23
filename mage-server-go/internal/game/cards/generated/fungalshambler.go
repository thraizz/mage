package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fungal Shambler", NewFungalShambler)
}

// NewFungalShambler creates a Fungal Shambler
// {4}{B}{G}{U} - CREATURE
// Trample
func NewFungalShambler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fungal Shambler")
	card.ManaCost = "{4}{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FUNGUS", "BEAST"}
	card.Power = "6"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	// card.AddAbility(ability1)
	return card, nil
}
