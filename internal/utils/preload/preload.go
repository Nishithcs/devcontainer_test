package preload

import "gorm.io/gorm"

func ApplyPreloads(query *gorm.DB, with []string) *gorm.DB {
	for _, relation := range with {
		query = query.Preload(relation)
	}
	return query
}
