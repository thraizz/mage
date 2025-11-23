package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tamiyo Seasoned Scholar", NewTamiyoSeasonedScholar)
}

// NewTamiyoSeasonedScholar creates a Tamiyo Seasoned Scholar
//   - PLANESWALKER
func NewTamiyoSeasonedScholar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tamiyo Seasoned Scholar")
	card.ManaCost = ""
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TAMIYO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewReturnToHandTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(-1, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
