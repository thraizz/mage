package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cyclopean Snare", NewCyclopeanSnare)
}

// NewCyclopeanSnare creates a Cyclopean Snare
// {2} - ARTIFACT
func NewCyclopeanSnare(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cyclopean Snare")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddTapCost().
		AddEffect(abilities.NewTapEffect()).
		AddEffect(abilities.NewReturnToHandSourceEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
