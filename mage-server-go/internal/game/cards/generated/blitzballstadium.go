package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Blitzball Stadium", NewBlitzballStadium)
}

// NewBlitzballStadium creates a Blitzball Stadium
// {X}{U} - ARTIFACT
func NewBlitzballStadium(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blitzball Stadium")
	card.ManaCost = "{X}{U}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(BlitzballStadiumValue.instance)).
		// TODO: GainAbilityTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
