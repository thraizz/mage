package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Fair Basilica", NewTheFairBasilica)
}

// NewTheFairBasilica creates a The Fair Basilica
//   - LAND
func NewTheFairBasilica(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Fair Basilica")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"SPHERE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
