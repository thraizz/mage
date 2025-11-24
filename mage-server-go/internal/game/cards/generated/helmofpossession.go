package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Helm Of Possession", NewHelmOfPossession)
}

// NewHelmOfPossession creates a Helm Of Possession
// {4} - ARTIFACT
func NewHelmOfPossession(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Helm Of Possession")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddTapCost().
		AddEffect(abilities.NewGainControlTargetEffect(abilities.DurationWhileControlled)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}