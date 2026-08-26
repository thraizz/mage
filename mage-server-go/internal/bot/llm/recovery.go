package llm

import (
	"context"
	"errors"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

// recovery.go is the failure matrix, ported from
// reference/pilot_recovery.py plus the constants block in §0.4 of the plan
// (mage-bench pilot.py:64-81 -- values arrived at over 60 harness epochs, not
// guessed).
//
// THE RULE IT ENCODES: a bot that cannot think must still take its turn. The
// table has three other players in it, and a seat that hangs waiting for an API
// that is not coming ruins the game for all of them. Every failure path here
// ends in a legal action -- normally a pass -- and the terminal state is
// mage-bench's in-character degradation: say something, then autopilot for the
// rest of the game. Nobody at the table has to know it was an outage; they just
// see a player who stopped doing much.
const (
	// MaxConsecutiveTimeouts is how many LLM timeouts in a row before the
	// transcript itself is treated as the suspect and reset.
	MaxConsecutiveTimeouts = 3
	// MaxConsecutiveTruncations is how many MaxTokens stops in a row before a
	// context reset. A model that keeps hitting the cap is usually reasoning
	// itself in circles over a transcript it cannot summarise.
	MaxConsecutiveTruncations = 3
	// MaxEmptyResponses is the lifetime budget for responses with no tool call
	// at all. Upstream's MAX_EMPTY_RESPONSES.
	MaxEmptyResponses = 10
	// MaxConsecutiveEmptyChoices bounds empty responses in a row, well below
	// the lifetime budget, so a stuck model degrades in one decision rather
	// than ten.
	MaxConsecutiveEmptyChoices = 5
	// MaxMalformedPerDecision is how many unusable tool calls one decision
	// tolerates before force-passing. One retry, then out -- the plan's
	// "malformed output retries once then falls back".
	MaxMalformedPerDecision = 2
	// MaxConsecutivePassErrors bounds forced passes in a row before autopilot.
	MaxConsecutivePassErrors = 3
)

// DegradationLine is mage-bench's in-character autopilot announcement
// (§0.4). It goes out through the ordinary ChatSource path, so it obeys the
// same per-turn chat cap as anything else the bot says.
const DegradationLine = "My brain is fried... going on autopilot for the rest of this game. GG!"

// StallLine is what the bot says when one decision fell over but the game is
// still worth playing. reference/pilot_recovery.py::_recover_from_stall.
const StallLine = "Brain freeze! Auto-passing this one..."

// FailureKind classifies what went wrong on one request.
type FailureKind int

const (
	// FailureNone means the request succeeded and produced a usable action.
	FailureNone FailureKind = iota
	// FailureTimeout is a deadline or cancellation on the request.
	FailureTimeout
	// FailureTruncated is stop_reason=max_tokens with no usable tool call.
	FailureTruncated
	// FailureEmpty is a well-formed response containing no tool call.
	FailureEmpty
	// FailureMalformed is a tool call that could not be resolved to a macro.
	FailureMalformed
	// FailureRetryable is a 429/500/529-class API error.
	FailureRetryable
	// FailurePermanent is 401/402/403/404 -- bad key, no credit, unknown model.
	// Retrying it forever is how a bad key turns into a bill-free game that
	// still takes three hours to not happen.
	FailurePermanent
)

// Recovery tracks failure counters across a game for one seat.
//
// Not safe for concurrent use; one per seat goroutine, like Conversation.
type Recovery struct {
	consecutiveTimeouts     int
	consecutiveTruncations  int
	consecutiveEmpty        int
	totalEmpty              int
	consecutivePassErrors   int
	autopilot               bool
	autopilotReason         string
	permanentFailureMessage string
}

// NewRecovery returns a zeroed Recovery.
func NewRecovery() *Recovery { return &Recovery{} }

// Autopilot reports whether the seat has degraded to auto-passing.
func (r *Recovery) Autopilot() bool { return r.autopilot }

// AutopilotReason reports why, for logs.
func (r *Recovery) AutopilotReason() string { return r.autopilotReason }

// PermanentFailure returns the message of a permanent API failure, if one has
// been seen. Non-empty means no further request will be made this game.
func (r *Recovery) PermanentFailure() string { return r.permanentFailureMessage }

// Success clears the consecutive counters. The lifetime ones stay.
func (r *Recovery) Success() {
	r.consecutiveTimeouts = 0
	r.consecutiveTruncations = 0
	r.consecutiveEmpty = 0
	r.consecutivePassErrors = 0
}

// Record folds one failure into the counters and reports what to do next.
type Action int

const (
	// ActionRetry means try the request again within this decision.
	ActionRetry Action = iota
	// ActionResetContext means reset the transcript, then retry.
	ActionResetContext
	// ActionForcePass means give up on this decision and take the fallback macro.
	ActionForcePass
	// ActionAutopilot means give up on the model for the rest of the game.
	ActionAutopilot
)

// Record updates the counters for kind and returns the recovery action.
//
// The ordering matters: a permanent failure short-circuits everything, a
// repeated failure escalates to a context reset before it escalates to
// autopilot, and everything that is not a clean success eventually reaches
// ActionForcePass rather than looping.
func (r *Recovery) Record(kind FailureKind) Action {
	switch kind {
	case FailureNone:
		r.Success()
		return ActionRetry

	case FailurePermanent:
		r.autopilot = true
		r.autopilotReason = "permanent API failure"
		return ActionAutopilot

	case FailureTimeout:
		r.consecutiveTimeouts++
		if r.consecutiveTimeouts >= MaxConsecutiveTimeouts {
			r.consecutiveTimeouts = 0
			r.autopilot = true
			r.autopilotReason = "repeated LLM timeouts"
			return ActionAutopilot
		}
		// One timeout is not a crisis; it is a slow request. Pass this
		// decision and try again on the next one.
		return ActionForcePass

	case FailureTruncated:
		r.consecutiveTruncations++
		if r.consecutiveTruncations >= MaxConsecutiveTruncations {
			r.consecutiveTruncations = 0
			return ActionResetContext
		}
		return ActionRetry

	case FailureEmpty:
		r.consecutiveEmpty++
		r.totalEmpty++
		if r.totalEmpty >= MaxEmptyResponses {
			r.autopilot = true
			r.autopilotReason = "empty response budget exhausted"
			return ActionAutopilot
		}
		if r.consecutiveEmpty >= MaxConsecutiveEmptyChoices {
			r.consecutiveEmpty = 0
			return ActionResetContext
		}
		return ActionRetry

	case FailureMalformed:
		// Malformed output is per-decision and the caller counts it: one
		// correction, then the fallback. Repeating the same bad call is not a
		// strategy the model recovers from by being asked a third time.
		return ActionRetry

	case FailureRetryable:
		return ActionRetry
	}
	return ActionForcePass
}

// ForcedPass records that a decision ended in a forced pass and reports whether
// the seat should degrade to autopilot.
func (r *Recovery) ForcedPass() bool {
	r.consecutivePassErrors++
	if r.consecutivePassErrors >= MaxConsecutivePassErrors {
		r.autopilot = true
		r.autopilotReason = "repeated forced passes"
		return true
	}
	return false
}

// Classify maps a transport error onto a FailureKind.
//
// Status buckets are §0.7's: 429/500/529 retryable (the SDK's own predicate
// also covers 408, 409 and everything >= 500), 400/401/403/404/413 not. 402 is
// added from reference/pilot_recovery.py::_classify_permanent_llm_failure,
// whose "credits exhausted" case is exactly the failure a long unattended
// simulation hits at 3am.
func Classify(err error) FailureKind {
	if err == nil {
		return FailureNone
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return FailureTimeout
	}
	var apierr *anthropic.Error
	if errors.As(err, &apierr) {
		switch apierr.StatusCode {
		case 401, 402, 403, 404:
			return FailurePermanent
		case 400, 413:
			// A rejected request will be rejected identically next time.
			return FailurePermanent
		case 408, 409, 429:
			return FailureRetryable
		}
		if apierr.StatusCode >= 500 {
			return FailureRetryable
		}
		return FailurePermanent
	}
	// A bare timeout from the HTTP layer does not always wrap
	// context.DeadlineExceeded.
	if s := err.Error(); strings.Contains(s, "context deadline exceeded") ||
		strings.Contains(s, "Timeout") || strings.Contains(s, "timeout") {
		return FailureTimeout
	}
	return FailureRetryable
}

// PermanentFailureReason mirrors reference/pilot_recovery.py:184-193: a 404
// without a 401 is a model that does not exist, anything else in the permanent
// bucket is a credential or credit problem. Worth distinguishing because one is
// a typo in a config file and the other is a bill.
func PermanentFailureReason(err error) string {
	if err == nil {
		return ""
	}
	var apierr *anthropic.Error
	if errors.As(err, &apierr) {
		switch apierr.StatusCode {
		case 404:
			return "model not found"
		case 401, 403:
			return "credentials rejected"
		case 402:
			return "credits exhausted"
		case 400, 413:
			return "request rejected"
		}
	}
	return err.Error()
}

// isTruncated reports whether a response stopped on the output cap.
func isTruncated(m *anthropic.BetaMessage) bool {
	return m != nil && m.StopReason == anthropic.BetaStopReasonMaxTokens
}
