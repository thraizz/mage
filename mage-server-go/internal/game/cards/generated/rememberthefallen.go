package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Remember The Fallen", NewRememberTheFallen)
}

// NewRememberTheFallen creates a Remember The Fallen
// {2}{W} - SORCERY
func NewRememberTheFallen(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Remember The Fallen")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ReturnFromGraveyardToHandTargetEffect with complex parameters
		// TODO: ReturnFromGraveyardToHandTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
