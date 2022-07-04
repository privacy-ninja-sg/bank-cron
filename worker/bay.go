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

func GetBayBalance() {
	url := "https://bay-api-u4zgi.ondigitalocean.app/api/v0.2/bank/krungsribiz"
	method := "POST"
	data := fmt.Sprintf(`{
	  "username": "%s",
	  "password": "%s"
	}`, os.Getenv("BAY_USR"), os.Getenv("BAY_PSW"))
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
	var resp model.BAYResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Printf("error to unmarshal response bay => %s", err)
	}
	db := drivers.ConnectPostgres()
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error with db %s", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Printf("error when ping %s", err)
	}
	for _, v := range resp.Data {
		var bayLog model.BAYLog
		amount, _ := strconv.ParseFloat(v.DepositCredit, 64)
		filter := &model.BAYLog{
			Amount:   amount,
			Datetime: v.CreateDate,
			Refer:    v.Refer,
			Title:    v.Transaction,
		}
		if query := db.Table("bay_log").Where(filter).First(&bayLog); query.Error != nil {
			log.Printf("error %s", query.Error)
			n := strings.Split(v.Transaction, " ")
			number := n[2][0:2]
			bankNumber := n[2][2:]
			bankType := utils.SwapBankToCode(number)
			if query.Error == gorm.ErrRecordNotFound {
				data := &model.BAYLog{
					Title:        v.Transaction,
					BankType:     bankType,
					BankNumber:   bankNumber,
					Amount:       amount,
					LatestAmount: float64(v.LatestCredit),
					Status:       "pending",
					Datetime:     v.CreateDate,
					Refer:        v.Refer,
				}
				log.Println("==== insert data ====")
				if insert := db.Table("bay_log").Create(data); insert.Error != nil {
					log.Println(insert.Error)
				}
			}
		}
	}
	defer sqlDB.Close()
}

func IncreaseBAYBalance() {
	log.Println("==== BAY Balance crons ====")
	db := drivers.ConnectPostgres()
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error with db %s", err)
	}
	var bayLog []model.BAYLog
	filter := model.BAYLog{
		Status: "pending",
	}

	if result := db.Table("bay_log").Where(filter).Find(&bayLog); result.Error != nil {
		log.Printf("error with database bay_log %s", result.Error)
	}

	for _, v := range bayLog {
		url := os.Getenv("WYNCLUB888_URL") + "/be/api/wynn/internal/bank/check"
		method := "POST"
		// bankNumber := strings.Split(v.BankNumber, "x")
		data := fmt.Sprintf(`{"last_bank_id": "%s"}`, v.BankNumber[4:])
		resp, err := utils.DoReq(method, url, data, true)
		if err != nil {
			log.Printf("error to call api %s", err)
		}
		var user model.CheckBankResponse
		if err := json.Unmarshal(resp, &user); err != nil {
			log.Printf("error to unmarshal %s", err)
		}
		if user.Code == 200 && user.Data.ID != 0 {
			log.Printf("increase credit to [%d] start", user.Data.ID)
			increase, err := wallet.CallIncreaseCredit(user.Data.Edges.Owner.ID, v.Amount)
			if err != nil {
				log.Printf("have some error to increase amount [%f] => uid [%d]", v.Amount, user.Data.Edges.Owner.ID)
			}
			log.Printf("%+v", increase)
			log.Printf("increase credit to [%d] success", user.Data.ID)
			if increase.Code == 200 {
				updateFilter := model.BAYLog{
					ID: v.ID,
				}
				if update := db.Table("bay_log").Where(updateFilter).Update("status", "success"); update.Error != nil {
					log.Printf("update error %s", update.Error)
				}
				logDeposit := &model.DepositLog{
					Title:       "Deposit",
					UserID:      user.Data.ID,
					ReceiveBank: "BAY",
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
