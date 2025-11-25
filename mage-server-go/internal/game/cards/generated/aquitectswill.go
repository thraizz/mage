package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Aquitects Will", NewAquitectsWill)
}

// NewAquitectsWill creates a Aquitects Will
// {U} - KINDRED SORCERY
func NewAquitectsWill(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aquitects Will")
	card.ManaCost = "{U}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"MERFOLK"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.NewCounter("flood", 1))).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddTarget(abilities.NewLandTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
