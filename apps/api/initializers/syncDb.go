package initializers

import model "envii/apps/api/models"

func SyncDb() {
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Account{})
}
