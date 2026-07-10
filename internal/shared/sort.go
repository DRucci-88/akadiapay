package shared

import "strings"

func ResolveSortClause(
	sortBy *string,
	order *string,
	allowed map[string]string,
	defaultColumn string,
	defaultOrder string,
) string {
	column := defaultColumn
	if sortBy != nil {
		if mappedColumn, ok := allowed[strings.ToLower(strings.TrimSpace(*sortBy))]; ok {
			column = mappedColumn
		}
	}

	direction := strings.ToUpper(strings.TrimSpace(defaultOrder))
	if order != nil {
		switch strings.ToUpper(strings.TrimSpace(*order)) {
		case "ASC":
			direction = "ASC"
		case "DESC":
			direction = "DESC"
		}
	}

	return column + " " + direction
}
