package initializers

import model "engov/apps/api/models"

func SyncDb() {
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Account{})
}
