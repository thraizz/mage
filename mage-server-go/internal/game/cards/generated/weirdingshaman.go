package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Weirding Shaman", NewWeirdingShaman)
}

// NewWeirdingShaman creates a Weirding Shaman
// {1}{B} - CREATURE
func NewWeirdingShaman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Weirding Shaman")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "SHAMAN"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("GoblinRogueToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
