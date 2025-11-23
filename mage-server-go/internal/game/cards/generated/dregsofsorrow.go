package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dregs Of Sorrow", NewDregsOfSorrow)
}

// NewDregsOfSorrow creates a Dregs Of Sorrow
// {X}{4}{B} - SORCERY
func NewDregsOfSorrow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dregs Of Sorrow")
	card.ManaCost = "{X}{4}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect("Destroy X target nonblack creatures")).
		AddEffect(abilities.NewDrawCardsEffect(GetXValue.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
