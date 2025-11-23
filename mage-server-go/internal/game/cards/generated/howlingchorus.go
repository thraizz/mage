package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Howling Chorus", NewHowlingChorus)
}

// NewHowlingChorus creates a Howling Chorus
//   - CREATURE
func NewHowlingChorus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Howling Chorus")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI", "WEREWOLF"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("EldraziHorrorToken")
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
