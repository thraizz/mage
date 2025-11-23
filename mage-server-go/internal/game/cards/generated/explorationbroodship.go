package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Exploration Broodship", NewExplorationBroodship)
}

// NewExplorationBroodship creates a Exploration Broodship
// {G} - ARTIFACT
// Flying
func NewExplorationBroodship(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Exploration Broodship")
	card.ManaCost = "{G}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"SPACECRAFT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
