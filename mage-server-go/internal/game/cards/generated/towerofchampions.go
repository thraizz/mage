package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tower Of Champions", NewTowerOfChampions)
}

// NewTowerOfChampions creates a Tower Of Champions
// {4} - ARTIFACT
func NewTowerOfChampions(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tower Of Champions")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{8}").
		AddTapCost().
		AddEffect(abilities.NewBoostEffect(6, 6)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
