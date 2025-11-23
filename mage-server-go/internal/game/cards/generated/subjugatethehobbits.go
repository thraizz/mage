package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Subjugate The Hobbits", NewSubjugateTheHobbits)
}

// NewSubjugateTheHobbits creates a Subjugate The Hobbits
// {5}{U}{U} - SORCERY
func NewSubjugateTheHobbits(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Subjugate The Hobbits")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainControlAllEffect(abilities.DurationCustom, abilities.NewAnyTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
