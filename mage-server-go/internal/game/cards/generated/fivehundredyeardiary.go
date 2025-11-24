package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Five Hundred Year Diary", NewFiveHundredYearDiary)
}

// NewFiveHundredYearDiary creates a Five Hundred Year Diary
// {3}{U} - ARTIFACT
func NewFiveHundredYearDiary(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Five Hundred Year Diary")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"CLUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}