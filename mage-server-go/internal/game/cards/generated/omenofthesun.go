package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Omen Of The Sun", NewOmenOfTheSun)
}

// NewOmenOfTheSun creates a Omen Of The Sun
// {2}{W} - ENCHANTMENT
// Flash
func NewOmenOfTheSun(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Omen Of The Sun")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewScryEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
