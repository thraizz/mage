package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wrens Run Hydra", NewWrensRunHydra)
}

// NewWrensRunHydra creates a Wrens Run Hydra
// {X}{G} - CREATURE
// Reach
func NewWrensRunHydra(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wrens Run Hydra")
	card.ManaCost = "{X}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HYDRA"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	return card, nil
}
