package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unity Of The Droids", NewUnityOfTheDroids)
}

// NewUnityOfTheDroids creates a Unity Of The Droids
// {W}{U}{B} - INSTANT
func NewUnityOfTheDroids(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unity Of The Droids")
	card.ManaCost = "{W}{U}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(4, 1, PutCards.HAND, PutCards.GRAVEYARD)
	// card.AddAbility(ability0)
	return card, nil
}
