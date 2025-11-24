package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wakening Suns Avatar", NewWakeningSunsAvatar)
}

// NewWakeningSunsAvatar creates a Wakening Suns Avatar
// {5}{W}{W}{W} - CREATURE
func NewWakeningSunsAvatar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wakening Suns Avatar")
	card.ManaCost = "{5}{W}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "AVATAR"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: DestroyAllEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
