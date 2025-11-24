package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vineshaper Prodigy", NewVineshaperProdigy)
}

// NewVineshaperProdigy creates a Vineshaper Prodigy
// {1}{G} - CREATURE
func NewVineshaperProdigy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vineshaper Prodigy")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 3, 1, PutCards.HAND, PutCards.BOT...)
	// card.AddAbility(ability0)
	return card, nil
}
