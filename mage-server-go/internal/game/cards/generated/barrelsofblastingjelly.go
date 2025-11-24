package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Barrels Of Blasting Jelly", NewBarrelsOfBlastingJelly)
}

// NewBarrelsOfBlastingJelly creates a Barrels Of Blasting Jelly
// {1} - ARTIFACT
func NewBarrelsOfBlastingJelly(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Barrels Of Blasting Jelly")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{5}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDamageEffect(5)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}