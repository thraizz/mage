package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Clamorous Ironclad", NewClamorousIronclad)
}

// NewClamorousIronclad creates a Clamorous Ironclad
// {3}{R} - ARTIFACT
// Menace
func NewClamorousIronclad(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Clamorous Ironclad")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Power = "6"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordMenace)
	card.AddAbility(ability0)
	return card, nil
}
