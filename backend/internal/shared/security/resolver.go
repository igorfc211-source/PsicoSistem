package security

import (
	"context"

	"github.com/google/uuid"
)

// IdentityResolver carrega a identidade atual do usuário a partir da persistência.
type IdentityResolver interface {
	ResolveInternalIdentity(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (Identity, error)
}
