package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Morningtide", NewMorningtide)
}

// NewMorningtide creates a Morningtide
// {1}{W} - SORCERY
func NewMorningtide(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Morningtide")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ExileGraveyardAllPlayersEffect()
	// card.AddAbility(ability0)
	return card, nil
}
