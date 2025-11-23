package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Akki Scrapchomper", NewAkkiScrapchomper)
}

// NewAkkiScrapchomper creates a Akki Scrapchomper
// {R} - CREATURE
func NewAkkiScrapchomper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Akki Scrapchomper")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "GOBLIN"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
