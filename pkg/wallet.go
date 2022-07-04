package wallet

import (
	"bank-crons/model"
	"bank-crons/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func CallIncreaseCredit(walletID int, amount float64) (*model.IncreaseResponse, error) {
	log.Printf("%d, %f", walletID, amount)
	url := os.Getenv("WYNCLUB888_URL") + "/be/api/wynn/internal/wallet/credit"
	method := "POST"
	data := fmt.Sprintf(`{"wallet_id": %d,"amount": %f,"remark":"scb-auto"}`, walletID, amount)
	resp, err := utils.DoReq(method, url, data, true)
	if err != nil {
		log.Printf("error to call api %s", err)
	}
	var increase model.IncreaseResponse
	if err := json.Unmarshal(resp, &increase); err != nil {
		log.Printf("error to unmarshal %s", err)
	}
	return &increase, err
}
