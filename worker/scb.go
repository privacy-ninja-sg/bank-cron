package worker

import (
	"bank-crons/drivers"
	"bank-crons/model"
	wallet "bank-crons/pkg"
	"bank-crons/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func GetSCBBalance() {
	log.Println("==== Get transaction ====")
	log.Printf("debug ENV => %s", os.Getenv("BANK_API_URL"))
	url := os.Getenv("BANK_API_URL") + "/api/scb"
	method := "POST"
	data := fmt.Sprintf(`{
	  "username": "%s",
	  "password": "%s"
	}`, os.Getenv("SCB_USR"), os.Getenv("SCB_PSW"))
	payload := strings.NewReader(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(body))
	var resp model.SCBResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Printf("error to unmarshal %s", err)
	}

	log.Printf("%+v", resp)
	db := drivers.ConnectPostgres()
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error with db %s", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Printf("error when ping %s", err)
	}
	for _, v := range resp.Transaction {
		var scbLog model.SCBLog
		log.Printf("%+v", v)
		amount, _ := strconv.ParseFloat(v.Credit, 64)
		filter := &model.SCBLog{
			BankNumber: v.BankNumber,
			Amount:     amount,
			Datetime:   v.Datetime,
			BankType:   v.Name,
		}
		if isInsert := db.Table("scb_log").Where(filter).First(&scbLog); isInsert.Error != nil {
			log.Printf("error when query => %s", err)
			if isInsert.Error == gorm.ErrRecordNotFound {
				data := &model.SCBLog{
					Title:      "SCB ENET",
					BankType:   v.Name,
					BankNumber: v.BankNumber,
					Amount:     amount,
					Datetime:   v.Datetime,
					Status:     "pending",
				}
				if result := db.Table("scb_log").Create(data); result.Error != nil {
					log.Println(result.Error)
				}
			}
		}
	}
	defer sqlDB.Close()
}

func IncreaseSCBBalance() {
	log.Println("==== SCB Balance crons ====")
	db := drivers.ConnectPostgres()
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error with db %s", err)
	}
	var scbLog []model.SCBLog
	filter := model.SCBLog{
		Status: "pending",
	}
	if result := db.Table("scb_log").Where(filter).Find(&scbLog); result.Error != nil {
		log.Printf("error with database %s", result.Error)
	}
	log.Printf("%+v", scbLog)
	for _, v := range scbLog {

		url := os.Getenv("WYNCLUB888_URL") + "/be/api/wynn/internal/bank/check"
		method := "POST"
		// bankNumber := strings.Split(v.BankNumber, "x")
		bankNumber := v.BankNumber[len(v.BankNumber)-4:]
		log.Printf("bank number => %+v", bankNumber)
		data := fmt.Sprintf(`{"last_bank_id": "%s"}`, bankNumber)
		resp, err := utils.DoReq(method, url, data, true)
		if err != nil {
			log.Printf("error to call api %s", err)
		}
		var user model.CheckBankResponse
		if err := json.Unmarshal(resp, &user); err != nil {
			log.Printf("error to unmarshal %s", err)
		}
		log.Printf("%+v", user)
		if user.Code == 200 && user.Data.ID != 0 {
			log.Printf("increase credit to [%d] start", user.Data.ID)
			increase, err := wallet.CallIncreaseCredit(user.Data.Edges.Owner.ID, v.Amount)
			if err != nil {
				log.Printf("have some error to increase amount [%f] => uid [%d]", v.Amount, user.Data.Edges.Owner.ID)
			}
			log.Printf("%+v", increase)
			log.Printf("increase credit to [%d] success", user.Data.ID)
			if increase.Code == 200 {
				updateFilter := model.SCBLog{
					ID: v.ID,
				}
				if update := db.Table("scb_log").Where(updateFilter).Update("status", "success"); update.Error != nil {
					log.Printf("update error %s", update.Error)
				}
				logDeposit := &model.DepositLog{
					Title:       "Deposit",
					UserID:      user.Data.ID,
					ReceiveBank: "SCB",
					BankNumber:  user.Data.BankAccountID,
					Amount:      v.Amount,
					Datetime:    v.Datetime,
					CreatedBy:   "BOT",
				}
				if insertLog := db.Table("deposit_log").Create(logDeposit); insertLog.Error != nil {
					log.Printf("error to insert log %s", insertLog.Error)
				}
			}
		}
	}
	defer sqlDB.Close()
}
