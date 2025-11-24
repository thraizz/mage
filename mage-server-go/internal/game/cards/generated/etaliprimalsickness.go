package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Etali Primal Sickness", NewEtaliPrimalSickness)
}

// NewEtaliPrimalSickness creates a Etali Primal Sickness
//   - CREATURE
//
// Trample, Indestructible
func NewEtaliPrimalSickness(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Etali Primal Sickness")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "ELDER", "DINOSAUR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "11"
	card.Toughness = "11"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability1)
	return card, nil
}
