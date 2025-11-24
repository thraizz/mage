package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rakdos The Muscle", NewRakdosTheMuscle)
}

// NewRakdosTheMuscle creates a Rakdos The Muscle
// {2}{B}{B}{R} - CREATURE
// Flying, Trample
func NewRakdosTheMuscle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rakdos The Muscle")
	card.ManaCost = "{2}{B}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON", "MERCENARY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	return card, nil
}