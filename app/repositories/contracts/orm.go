package contracts

import contractsorm "github.com/goravel/framework/contracts/database/orm"

type ORM interface {
	Query() contractsorm.Query
}
