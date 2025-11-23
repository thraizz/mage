package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Akki Blizzard Herder", NewAkkiBlizzardHerder)
}

// NewAkkiBlizzardHerder creates a Akki Blizzard Herder
// {1}{R} - CREATURE
func NewAkkiBlizzardHerder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Akki Blizzard Herder")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "SHAMAN"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeAllEffect(filter)
	// card.AddAbility(ability0)
	return card, nil
}
