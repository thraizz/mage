package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stream Of Unconsciousness", NewStreamOfUnconsciousness)
}

// NewStreamOfUnconsciousness creates a Stream Of Unconsciousness
// {U} - KINDRED INSTANT
func NewStreamOfUnconsciousness(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stream Of Unconsciousness")
	card.ManaCost = "{U}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"WIZARD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(-4, 0)).
		AddEffect(abilities.NewConditionalEffect(abilities.NewDrawCardsEffect(1), "unknown")).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
