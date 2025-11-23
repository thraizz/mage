package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ghalta Primal Hunger", NewGhaltaPrimalHunger)
}

// NewGhaltaPrimalHunger creates a Ghalta Primal Hunger
// {10}{G}{G} - CREATURE
// Trample
func NewGhaltaPrimalHunger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ghalta Primal Hunger")
	card.ManaCost = "{10}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDER", "DINOSAUR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "12"
	card.Toughness = "12"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
