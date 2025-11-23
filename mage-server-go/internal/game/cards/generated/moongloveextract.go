package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Moonglove Extract", NewMoongloveExtract)
}

// NewMoongloveExtract creates a Moonglove Extract
// {3} - ARTIFACT
func NewMoongloveExtract(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Moonglove Extract")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDamageEffect(2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
