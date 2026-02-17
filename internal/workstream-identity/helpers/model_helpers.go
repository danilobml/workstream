package helpers

import "github.com/danilobml/workstream/internal/platform/models"

func ParseRoles(names []string) ([]models.Role, error) {
	roles := make([]models.Role, 0, len(names))
	for _, name := range names {
		role, err := models.ParseRole(name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func GetRoleNames(roles []models.Role) []string {
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		roleName := role.GetName()
		names = append(names, roleName)
	}

	return names
}
