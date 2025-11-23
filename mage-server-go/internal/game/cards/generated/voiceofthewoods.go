package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Voice Of The Woods", NewVoiceOfTheWoods)
}

// NewVoiceOfTheWoods creates a Voice Of The Woods
// {3}{G}{G} - CREATURE
func NewVoiceOfTheWoods(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Voice Of The Woods")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("VoiceOfTheWoodsElementalToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
