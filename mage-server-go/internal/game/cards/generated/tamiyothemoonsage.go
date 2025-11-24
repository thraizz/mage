package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tamiyo The Moon Sage", NewTamiyoTheMoonSage)
}

// NewTamiyoTheMoonSage creates a Tamiyo The Moon Sage
// {3}{U}{U} - PLANESWALKER
func NewTamiyoTheMoonSage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tamiyo The Moon Sage")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TAMIYO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: TapTargetEffect with complex parameters
		AddEffect(abilities.NewDrawCardsEffect(TappedCreaturesControlledByTargetCount.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
