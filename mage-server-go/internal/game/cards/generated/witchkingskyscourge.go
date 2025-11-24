package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Witch King Sky Scourge", NewWitchKingSkyScourge)
}

// NewWitchKingSkyScourge creates a Witch King Sky Scourge
// {5}{B}{R} - CREATURE
// Flying
func NewWitchKingSkyScourge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Witch King Sky Scourge")
	card.ManaCost = "{5}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WRAITH", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}