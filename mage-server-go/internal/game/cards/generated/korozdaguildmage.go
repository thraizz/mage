package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Korozda Guildmage", NewKorozdaGuildmage)
}

// NewKorozdaGuildmage creates a Korozda Guildmage
// {B}{G} - CREATURE
func NewKorozdaGuildmage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Korozda Guildmage")
	card.ManaCost = "{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "SHAMAN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1)).
		AddEffect(abilities.NewGrantAbilityEffect("IntimidateAbility", effects.DurationEndOfTurn)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("SaprolingToken")
	if err != nil {
		return nil, err
	}
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewCreateTokenEffect(token1_0, SacrificeCostCreaturesToughness.instance)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
