package abilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/token"
)

// ========================================
// Create Token Effect
// ========================================

// CreateTokenEffect creates token(s) on the battlefield.
// Mirrors Java CreateTokenEffect.
type CreateTokenEffect struct {
	token     *token.Token
	amount    int
	tapped    bool
	attacking bool

	// Optional: counters to add when tokens enter
	counterType   *counters.CounterType
	counterAmount int

	lastAddedIDs []uuid.UUID
}

// NewCreateTokenEffect creates an effect that creates one token.
func NewCreateTokenEffect(tok *token.Token) *CreateTokenEffect {
	return &CreateTokenEffect{
		token:        tok.Copy(),
		amount:       1,
		tapped:       false,
		attacking:    false,
		lastAddedIDs: make([]uuid.UUID, 0),
	}
}

// NewCreateTokenEffectAmount creates an effect that creates the specified number of tokens.
func NewCreateTokenEffectAmount(tok *token.Token, amount int) *CreateTokenEffect {
	return &CreateTokenEffect{
		token:        tok.Copy(),
		amount:       amount,
		tapped:       false,
		attacking:    false,
		lastAddedIDs: make([]uuid.UUID, 0),
	}
}

// NewCreateTokenEffectTapped creates an effect that creates token(s) tapped.
func NewCreateTokenEffectTapped(tok *token.Token, amount int, tapped bool) *CreateTokenEffect {
	return &CreateTokenEffect{
		token:        tok.Copy(),
		amount:       amount,
		tapped:       tapped,
		attacking:    false,
		lastAddedIDs: make([]uuid.UUID, 0),
	}
}

// NewCreateTokenEffectAttacking creates an effect that creates token(s) tapped and attacking.
func NewCreateTokenEffectAttacking(tok *token.Token, amount int, tapped, attacking bool) *CreateTokenEffect {
	return &CreateTokenEffect{
		token:        tok.Copy(),
		amount:       amount,
		tapped:       tapped,
		attacking:    attacking,
		lastAddedIDs: make([]uuid.UUID, 0),
	}
}

// WithCounters adds counters to tokens as they enter the battlefield.
func (e *CreateTokenEffect) WithCounters(counterType counters.CounterType, amount int) *CreateTokenEffect {
	e.counterType = &counterType
	e.counterAmount = amount
	return e
}

func (e *CreateTokenEffect) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if e.token == nil || e.amount <= 0 {
		return nil
	}

	// Cast to TokenGameContext for token operations
	tokenGame, ok := game.(TokenGameContext)
	if !ok {
		return fmt.Errorf("game context does not support token operations")
	}

	// Create the tokens
	createdIDs, err := tokenGame.CreateTokens(e.token, e.amount, source, e.tapped, e.attacking)
	if err != nil {
		return fmt.Errorf("failed to create tokens: %w", err)
	}

	e.lastAddedIDs = createdIDs

	// Add counters to newly created tokens if specified
	if e.counterType != nil && e.counterAmount > 0 {
		counter := e.counterType.CreateInstance(e.counterAmount)
		for _, tokenID := range createdIDs {
			permanent, err := tokenGame.GetPermanent(tokenID)
			if err != nil {
				continue // Token may have been destroyed/removed
			}
			if err := tokenGame.AddCountersToPermanent(permanent, counter.Copy()); err != nil {
				return fmt.Errorf("failed to add counters to token: %w", err)
			}
		}
	}

	// Inform players
	if len(createdIDs) > 0 {
		if e.amount == 1 {
			tokenGame.InformPlayers(fmt.Sprintf("Created %s", e.token.Description))
		} else {
			tokenGame.InformPlayers(fmt.Sprintf("Created %d %s", e.amount, e.token.Description))
		}
	}

	return nil
}

func (e *CreateTokenEffect) GetDescription() string {
	if e.token == nil {
		return ""
	}

	desc := "create "

	if e.amount == 1 {
		desc += "a " + e.token.Description
	} else {
		desc += fmt.Sprintf("%d %s", e.amount, e.token.Description)
	}

	if e.counterType != nil && e.counterAmount > 0 {
		if e.counterAmount == 1 {
			desc += fmt.Sprintf(" with a %s counter on it", e.counterType.String())
		} else {
			desc += fmt.Sprintf(" with %d %s counters on it", e.counterAmount, e.counterType.String())
		}
	}

	if e.tapped {
		desc += " tapped"
	}

	if e.attacking {
		desc += " and attacking"
	}

	return desc
}

// GetLastAddedTokenIDs returns the IDs of tokens created by this effect.
func (e *CreateTokenEffect) GetLastAddedTokenIDs() []uuid.UUID {
	result := make([]uuid.UUID, len(e.lastAddedIDs))
	copy(result, e.lastAddedIDs)
	return result
}

// ========================================
// Create Token Effect with Dynamic Amount
// ========================================

// CreateTokenEffectDynamic creates token(s) with a dynamically calculated amount.
// Mirrors Java CreateTokenEffect with DynamicValue parameter.
type CreateTokenEffectDynamic struct {
	token        *token.Token
	amountValue  DynamicValue
	tapped       bool
	attacking    bool
	lastAddedIDs []uuid.UUID
}

// NewCreateTokenEffectDynamic creates an effect that creates tokens based on a dynamic value.
func NewCreateTokenEffectDynamic(tok *token.Token, amount DynamicValue) *CreateTokenEffectDynamic {
	return &CreateTokenEffectDynamic{
		token:        tok.Copy(),
		amountValue:  amount,
		tapped:       false,
		attacking:    false,
		lastAddedIDs: make([]uuid.UUID, 0),
	}
}

// NewCreateTokenEffectDynamicTapped creates an effect that creates tokens tapped based on a dynamic value.
func NewCreateTokenEffectDynamicTapped(tok *token.Token, amount DynamicValue, tapped bool) *CreateTokenEffectDynamic {
	return &CreateTokenEffectDynamic{
		token:        tok.Copy(),
		amountValue:  amount,
		tapped:       tapped,
		attacking:    false,
		lastAddedIDs: make([]uuid.UUID, 0),
	}
}

func (e *CreateTokenEffectDynamic) Apply(ctx context.Context, game GameContext, source uuid.UUID, targets []uuid.UUID) error {
	if e.token == nil {
		return nil
	}

	// Calculate the dynamic amount
	amount := e.amountValue.Calculate(ctx, game, source)
	if amount <= 0 {
		return nil
	}

	// Cast to TokenGameContext for token operations
	tokenGame, ok := game.(TokenGameContext)
	if !ok {
		return fmt.Errorf("game context does not support token operations")
	}

	// Create the tokens
	createdIDs, err := tokenGame.CreateTokens(e.token, amount, source, e.tapped, e.attacking)
	if err != nil {
		return fmt.Errorf("failed to create tokens: %w", err)
	}

	e.lastAddedIDs = createdIDs

	// Inform players
	if len(createdIDs) > 0 {
		if amount == 1 {
			tokenGame.InformPlayers(fmt.Sprintf("Created %s", e.token.Description))
		} else {
			tokenGame.InformPlayers(fmt.Sprintf("Created %d %s", amount, e.token.Description))
		}
	}

	return nil
}

func (e *CreateTokenEffectDynamic) GetDescription() string {
	if e.token == nil {
		return ""
	}

	desc := fmt.Sprintf("create X %s where X is %s", e.token.Description, e.amountValue.GetMessage())

	if e.tapped {
		desc += " tapped"
	}

	if e.attacking {
		desc += " and attacking"
	}

	return desc
}

// GetLastAddedTokenIDs returns the IDs of tokens created by this effect.
func (e *CreateTokenEffectDynamic) GetLastAddedTokenIDs() []uuid.UUID {
	result := make([]uuid.UUID, len(e.lastAddedIDs))
	copy(result, e.lastAddedIDs)
	return result
}

// ========================================
// Extended GameContext Interface for Tokens
// ========================================

// TokenGameContext extends CounterGameContext with token-specific methods.
// This should be implemented by the game engine.
type TokenGameContext interface {
	CounterGameContext

	// CreateTokens creates the specified number of tokens on the battlefield.
	// Returns the UUIDs of the created permanents.
	CreateTokens(token *token.Token, amount int, source uuid.UUID, tapped, attacking bool) ([]uuid.UUID, error)
}
