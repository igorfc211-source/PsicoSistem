package permissions

import (
	"strings"

	sharederrors "api-on/internal/shared/errors"
)

type Scope string

const (
	ScopeNone Scope = "none"
	ScopeOwn  Scope = "own"
	ScopeAll  Scope = "all"
)

type AccountPermissions struct {
	UserDirectory Scope `json:"user_directory"`
	Patients      Scope `json:"patients"`
	Services      Scope `json:"services"`
	Calendar      Scope `json:"calendar"`
	Finance       Scope `json:"finance"`
	AIHistory     Scope `json:"ai_history"`
	Plans         Scope `json:"plans"`
}

// DefaultForRole define o conjunto-base de permissões de cada papel interno.
func DefaultForRole(role string) AccountPermissions {
	switch strings.TrimSpace(role) {
	case "owner", "admin":
		return allScopes()
	case "coordinator":
		return AccountPermissions{
			UserDirectory: ScopeOwn,
			Patients:      ScopeAll,
			Services:      ScopeAll,
			Calendar:      ScopeAll,
			Finance:       ScopeNone,
			AIHistory:     ScopeAll,
			Plans:         ScopeAll,
		}
	case "financial":
		return AccountPermissions{
			UserDirectory: ScopeOwn,
			Patients:      ScopeNone,
			Services:      ScopeNone,
			Calendar:      ScopeNone,
			Finance:       ScopeAll,
			AIHistory:     ScopeNone,
			Plans:         ScopeNone,
		}
	case "professional":
		fallthrough
	default:
		return ownScopes()
	}
}

func Validate(value AccountPermissions) error {
	for _, scope := range []Scope{
		value.UserDirectory,
		value.Patients,
		value.Services,
		value.Calendar,
		value.Finance,
		value.AIHistory,
		value.Plans,
	} {
		if !isValidScope(scope) {
			return sharederrors.Invalid("INVALID_PERMISSION_SCOPE", "invalid account permission scope")
		}
	}

	return nil
}

func Normalize(role string, input *AccountPermissions) (AccountPermissions, error) {
	if input == nil {
		return DefaultForRole(role), nil
	}

	normalized := AccountPermissions{
		UserDirectory: Scope(strings.TrimSpace(string(input.UserDirectory))),
		Patients:      Scope(strings.TrimSpace(string(input.Patients))),
		Services:      Scope(strings.TrimSpace(string(input.Services))),
		Calendar:      Scope(strings.TrimSpace(string(input.Calendar))),
		Finance:       Scope(strings.TrimSpace(string(input.Finance))),
		AIHistory:     Scope(strings.TrimSpace(string(input.AIHistory))),
		Plans:         Scope(strings.TrimSpace(string(input.Plans))),
	}

	if err := Validate(normalized); err != nil {
		return AccountPermissions{}, err
	}

	return normalized, nil
}

func (p AccountPermissions) CanViewAllUsers() bool {
	return p.UserDirectory == ScopeAll
}

func (p AccountPermissions) CanAccessOwnClinicalData() bool {
	return p.Patients == ScopeOwn ||
		p.Services == ScopeOwn ||
		p.Calendar == ScopeOwn ||
		p.Finance == ScopeOwn ||
		p.AIHistory == ScopeOwn ||
		p.Plans == ScopeOwn
}

func allScopes() AccountPermissions {
	return AccountPermissions{
		UserDirectory: ScopeAll,
		Patients:      ScopeAll,
		Services:      ScopeAll,
		Calendar:      ScopeAll,
		Finance:       ScopeAll,
		AIHistory:     ScopeAll,
		Plans:         ScopeAll,
	}
}

func ownScopes() AccountPermissions {
	return AccountPermissions{
		UserDirectory: ScopeOwn,
		Patients:      ScopeOwn,
		Services:      ScopeOwn,
		Calendar:      ScopeOwn,
		Finance:       ScopeOwn,
		AIHistory:     ScopeOwn,
		Plans:         ScopeOwn,
	}
}

func isValidScope(scope Scope) bool {
	switch scope {
	case ScopeNone, ScopeOwn, ScopeAll:
		return true
	default:
		return false
	}
}
