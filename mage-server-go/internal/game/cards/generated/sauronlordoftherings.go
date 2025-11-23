package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sauron Lord Of The Rings", NewSauronLordOfTheRings)
}

// NewSauronLordOfTheRings creates a Sauron Lord Of The Rings
// {5}{U}{B}{R} - CREATURE
// Trample
func NewSauronLordOfTheRings(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sauron Lord Of The Rings")
	card.ManaCost = "{5}{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR", "HORROR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "9"
	card.Toughness = "9"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
