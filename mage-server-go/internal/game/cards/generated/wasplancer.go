package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wasp Lancer", NewWaspLancer)
}

// NewWaspLancer creates a Wasp Lancer
// {U/B}{U/B}{U/B} - CREATURE
// Flying
func NewWaspLancer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wasp Lancer")
	card.ManaCost = "{U/B}{U/B}{U/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FAERIE", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}