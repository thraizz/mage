package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tapestry Warden", NewTapestryWarden)
}

// NewTapestryWarden creates a Tapestry Warden
// {3}{G} - ARTIFACT CREATURE
// Vigilance
func NewTapestryWarden(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tapestry Warden")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
