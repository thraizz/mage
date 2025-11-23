package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unstable Obelisk", NewUnstableObelisk)
}

// NewUnstableObelisk creates a Unstable Obelisk
// {3} - ARTIFACT
func NewUnstableObelisk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unstable Obelisk")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{7}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
