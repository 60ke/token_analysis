package main

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestDBInit(t *testing.T) {
	var eth FXHEthHolder
	var eths []FXHEthHolder
	dsn := "root:root@tcp(127.0.0.1:3306)/trias_analysis?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.Order("version desc").Find(&eths)
	fmt.Println(eths[0])

	db.First(&eths, "id = ?", "201")
	fmt.Println(eths)
	fmt.Println(len(eths))

	db.Order("version desc").Find(&eth)
	fmt.Println(eth.Version)
}

func TestDBCreate(t *testing.T) {
	dsn := "root:root@tcp(127.0.0.1:3306)/trias_analysis?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// db.AutoMigrate(&Transfer{})
	db.Table("transfer_test").AutoMigrate(&Transfer{})
}

func TestDBSave(t *testing.T) {

	var transfer Transfer
	transfer.Block = "0x1628ca7"
	transfer.Time = 1669087415
	transfer.Hash = "0x0490f5e2a3b6c832530f6a798119eaaa65077b5f5ce2d371613980d32e562043"
	dsn := "root:root@tcp(127.0.0.1:3306)/trias_analysis?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// db.Save(&transfer)
	db.Table("transfer_test").Create(&transfer)
}

func TestDBUpdate(t *testing.T) {

	var transfer Transfer
	transfer.ID = 141
	transfer.From = "1251"
	dsn := "root:root@tcp(127.0.0.1:3306)/trias_analysis?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// db.Table("transfer_test").FirstOrCreate(&transfer)
	// db.Table("transfer_test").FirstOrCreate(&transfer, Transfer{Hash: "0x0490f5e2a3b6c832530f6a798119eaaa65077b5f5ce2d371613980d32e562047"})
	db.Table("transfer_test").Save(&transfer)

	if db.Table("transfer_test").Where("id = ?", 141).Updates(&transfer).RowsAffected == 0 {
		db.Create(&transfer)
	}
	// time.Sleep(time.Second * 3)
	db.Table("transfer_test").Where(Transfer{Hash: "0x0490f5e2a3b6c832530f6a798119eaaa65077b5f5ce2d371613980d32e562049"}).Assign(Transfer{ID: 20}).FirstOrCreate(&transfer)
	transfer1 := transfer
	transfer1.From = "111"
	db.Table("transfer_test").FirstOrCreate(&transfer1)
}

func TestGetHolders(t *testing.T) {
	var holders []Holder
	dsn := "root:root@tcp(127.0.0.1:3306)/trias_analysis?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// table := getHolderTable("BSC")
	db.Table("bsc_holder").Find(&holders)
	t.Log(len(holders))
	// fmt.Println(holder)

}
