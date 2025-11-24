package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Huatli The Suns Heart", NewHuatliTheSunsHeart)
}

// NewHuatliTheSunsHeart creates a Huatli The Suns Heart
// {2}{G/W} - PLANESWALKER
func NewHuatliTheSunsHeart(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Huatli The Suns Heart")
	card.ManaCost = "{2}{G/W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"HUATLI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(GreatestAmongPermanentsValue.TOUGHNESS_CONTROLLED_CREATURES)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
