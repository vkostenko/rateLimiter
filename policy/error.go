package policy

import "fmt"

type unknownPolicyError struct {
	policy Policy
}

func newUnknownPolicyError(policy Policy) unknownPolicyError {
	return unknownPolicyError{
		policy: policy,
	}
}

func (e unknownPolicyError) Error() string {
	return fmt.Sprintf("unknown rate limiting policy %d", e.policy)
}
