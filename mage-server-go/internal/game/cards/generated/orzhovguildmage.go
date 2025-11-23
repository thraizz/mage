package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Orzhov Guildmage", NewOrzhovGuildmage)
}

// NewOrzhovGuildmage creates a Orzhov Guildmage
// {W/B}{W/B} - CREATURE
func NewOrzhovGuildmage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Orzhov Guildmage")
	card.ManaCost = "{W/B}{W/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewGainLifeEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
