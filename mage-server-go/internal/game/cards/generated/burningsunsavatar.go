package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Burning Suns Avatar", NewBurningSunsAvatar)
}

// NewBurningSunsAvatar creates a Burning Suns Avatar
// {3}{R}{R}{R} - CREATURE
func NewBurningSunsAvatar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Burning Suns Avatar")
	card.ManaCost = "{3}{R}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "AVATAR"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(3)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
