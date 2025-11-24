package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Merciless Eviction", NewMercilessEviction)
}

// NewMercilessEviction creates a Merciless Eviction
// {4}{W}{B} - SORCERY
func NewMercilessEviction(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Merciless Eviction")
	card.ManaCost = "{4}{W}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewExileAllEffect(abilities.NewCreatureTargetFilter())).
		AddEffect(abilities.NewExileAllEffect(abilities.NewAnyTargetFilter())).
		AddEffect(abilities.NewExileAllEffect(abilities.NewCreatureTargetFilter())).
		AddEffect(abilities.NewExileAllEffect(abilities.NewAnyTargetFilter())).
		AddEffect(abilities.NewExileAllEffect(abilities.NewAnyTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}