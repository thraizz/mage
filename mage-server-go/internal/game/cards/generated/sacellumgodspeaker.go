package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sacellum Godspeaker", NewSacellumGodspeaker)
}

// NewSacellumGodspeaker creates a Sacellum Godspeaker
// {2}{G} - CREATURE
func NewSacellumGodspeaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sacellum Godspeaker")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "DRUID"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SimpleManaAbility
	//   - Effect: SacellumGodspeakerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
