package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Soramaro First To Dream", NewSoramaroFirstToDream)
}

// NewSoramaroFirstToDream creates a Soramaro First To Dream
// {4}{U}{U} - CREATURE
// Flying
func NewSoramaroFirstToDream(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Soramaro First To Dream")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
