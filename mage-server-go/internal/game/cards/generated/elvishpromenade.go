package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Elvish Promenade", NewElvishPromenade)
}

// NewElvishPromenade creates a Elvish Promenade
// {3}{G} - KINDRED SORCERY
func NewElvishPromenade(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elvish Promenade")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"ELF"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("ElfWarriorToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, elfCount)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}