package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("River Herald Guide", NewRiverHeraldGuide)
}

// NewRiverHeraldGuide creates a River Herald Guide
// {2}{G} - CREATURE
// Vigilance
func NewRiverHeraldGuide(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "River Herald Guide")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "SCOUT"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
