package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Notorious Throng", NewNotoriousThrong)
}

// NewNotoriousThrong creates a Notorious Throng
// {3}{U} - KINDRED SORCERY
func NewNotoriousThrong(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Notorious Throng")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"ROGUE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ConditionalOneShotEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
