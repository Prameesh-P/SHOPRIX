package main


func SyncDB(){
	DB.AutoMigrate(
		&User{},
	)
}