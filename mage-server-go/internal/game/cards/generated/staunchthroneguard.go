package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Staunch Throneguard", NewStaunchThroneguard)
}

// NewStaunchThroneguard creates a Staunch Throneguard
// {5} - ARTIFACT CREATURE
// Vigilance
func NewStaunchThroneguard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Staunch Throneguard")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
