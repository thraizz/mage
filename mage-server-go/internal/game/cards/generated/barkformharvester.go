package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Barkform Harvester", NewBarkformHarvester)
}

// NewBarkformHarvester creates a Barkform Harvester
// {3} - ARTIFACT CREATURE
// Reach
func NewBarkformHarvester(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Barkform Harvester")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - PutOnLibraryTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	// card.AddAbility(ability1)
	return card, nil
}
