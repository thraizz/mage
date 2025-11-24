package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cateran Overlord", NewCateranOverlord)
}

// NewCateranOverlord creates a Cateran Overlord
// {4}{B}{B}{B} - CREATURE
func NewCateranOverlord(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cateran Overlord")
	card.ManaCost = "{4}{B}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HORROR", "MERCENARY"}
	card.Power = "7"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddManaCost("{6}").
		// TODO: SearchLibraryPutInPlayEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
