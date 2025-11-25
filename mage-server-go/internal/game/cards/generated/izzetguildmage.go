package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Izzet Guildmage", NewIzzetGuildmage)
}

// NewIzzetGuildmage creates a Izzet Guildmage
// {U/R}{U/R} - CREATURE
func NewIzzetGuildmage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Izzet Guildmage")
	card.ManaCost = "{U/R}{U/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - CopyTargetStackObjectEffect()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - CopyTargetStackObjectEffect()
	// card.AddAbility(ability1)
	return card, nil
}
