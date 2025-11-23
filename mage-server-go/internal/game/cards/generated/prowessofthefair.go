package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Prowess Of The Fair", NewProwessOfTheFair)
}

// NewProwessOfTheFair creates a Prowess Of The Fair
// {1}{B} - KINDRED ENCHANTMENT
func NewProwessOfTheFair(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prowess Of The Fair")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"ELF"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("ElfWarriorToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
