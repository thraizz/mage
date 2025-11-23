package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aetherflux Conduit", NewAetherfluxConduit)
}

// NewAetherfluxConduit creates a Aetherflux Conduit
// {6} - ARTIFACT
func NewAetherfluxConduit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aetherflux Conduit")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(7)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
