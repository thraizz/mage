package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tower Of Fortunes", NewTowerOfFortunes)
}

// NewTowerOfFortunes creates a Tower Of Fortunes
// {4} - ARTIFACT
func NewTowerOfFortunes(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tower Of Fortunes")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{8}").
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(4)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
