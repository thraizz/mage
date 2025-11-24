package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tamiyo Collector Of Tales", NewTamiyoCollectorOfTales)
}

// NewTamiyoCollectorOfTales creates a Tamiyo Collector Of Tales
// {2}{G}{U} - PLANESWALKER
func NewTamiyoCollectorOfTales(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tamiyo Collector Of Tales")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TAMIYO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewReturnFromGraveyardToHandTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}