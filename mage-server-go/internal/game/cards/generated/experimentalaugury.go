package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Experimental Augury", NewExperimentalAugury)
}

// NewExperimentalAugury creates a Experimental Augury
// {1}{U} - INSTANT
func NewExperimentalAugury(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Experimental Augury")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 3, 1, PutCards.HAND, PutCards.BOT...)
	// card.AddAbility(ability0)
	return card, nil
}
