package initializers

import model "envi/apps/api/models"

func SyncDb() {
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Account{})
}
