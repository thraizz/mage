package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Progenitor Mimic", NewProgenitorMimic)
}

// NewProgenitorMimic creates a Progenitor Mimic
// {4}{G}{U} - CREATURE
func NewProgenitorMimic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Progenitor Mimic")
	card.ManaCost = "{4}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(new AbilityCopyApplier(                 new Beginn...)
	// card.AddAbility(ability0)
	return card, nil
}
