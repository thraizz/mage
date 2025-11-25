package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Insatiable Avarice", NewInsatiableAvarice)
}

// NewInsatiableAvarice creates a Insatiable Avarice
// {B} - SORCERY
func NewInsatiableAvarice(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Insatiable Avarice")
	card.ManaCost = "{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: SearchLibraryPutOnLibraryEffect with complex parameters
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewLoseLifeEffect(3)).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
