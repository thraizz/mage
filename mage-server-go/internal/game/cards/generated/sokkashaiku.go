package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Sokkas Haiku", NewSokkasHaiku)
}

// NewSokkasHaiku creates a Sokkas Haiku
// {3}{U}{U} - INSTANT
func NewSokkasHaiku(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sokkas Haiku")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCounterSpellEffect()).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewMillCardsControllerEffect(1)).
		AddEffect(abilities.NewUntapEffect("untap target land")).
		AddTarget(abilities.NewSpellTargetFilter()).
		AddTarget(abilities.NewLandTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
