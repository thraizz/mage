package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Boros Guildmage", NewBorosGuildmage)
}

// NewBorosGuildmage creates a Boros Guildmage
// {R/W}{R/W} - CREATURE
func NewBorosGuildmage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Boros Guildmage")
	card.ManaCost = "{R/W}{R/W}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewGrantAbilityEffect("HasteAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewGrantAbilityEffect("FirstStrikeAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
