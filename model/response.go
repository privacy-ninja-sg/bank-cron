package model

import "time"

type SCBResponse struct {
	TotalBalance string        `json:"total_balance"`
	Transaction  []Transaction `json:"transaction"`
}

type Transaction struct {
	Datetime   string `json:"datetime"`
	Credit     string `json:"credit"`
	Name       string `json:"name"`
	BankNumber string `json:"bankNumber"`
}

type SCBLog struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title"`
	BankType   string    `json:"bank_type"`
	BankNumber string    `json:"bank_number"`
	Amount     float64   `json:"amount"`
	Datetime   string    `json:"datetime"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CheckBankResponse struct {
	S    string `json:"s"`
	Code int    `json:"code"`
	Data struct {
		ID                int       `json:"id"`
		UUID              string    `json:"uuid"`
		BankAccountID     string    `json:"bank_account_id"`
		BankAccountIDLast string    `json:"bank_account_id_last"`
		BankAccountName   string    `json:"bank_account_name"`
		Status            string    `json:"status"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		Edges             struct {
			Owner struct {
				ID        int       `json:"id"`
				UUID      string    `json:"uuid"`
				Tel       string    `json:"tel"`
				Username  string    `json:"username"`
				Status    string    `json:"status"`
				Bonus     string    `json:"bonus"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
				Edges     struct {
				} `json:"edges"`
			} `json:"owner"`
			Bank struct {
				ID        int       `json:"id"`
				UUID      string    `json:"uuid"`
				Name      string    `json:"name"`
				ShortName string    `json:"short_name"`
				Logo      string    `json:"logo"`
				Status    string    `json:"status"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
				Edges     struct {
				} `json:"edges"`
			} `json:"bank"`
		} `json:"edges"`
	} `json:"data,omitempty"`
}

type DepositLog struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	BankNumber  string    `json:"bank_number"`
	Amount      float64   `json:"amount"`
	Datetime    string    `json:"datetime"`
	ReceiveBank string    `json:"receive_bank"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type IncreaseResponse struct {
	S    string `json:"s"`
	Code int    `json:"code"`
	Data struct {
		ID        int       `json:"id"`
		UUID      string    `json:"uuid"`
		Debit     int       `json:"debit"`
		Balance   int       `json:"balance"`
		Remark    string    `json:"remark"`
		TxnType   string    `json:"txn_type"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Edges     struct {
		} `json:"edges"`
	} `json:"data"`
}

type BAYResponse struct {
	Message string `json:"message"`
	Data    []struct {
		CreateDate    string `json:"createDate"`
		Transaction   string `json:"transaction"`
		DepositCredit string `json:"depositCredit"`
		LatestCredit  int    `json:"latestCredit"`
		Refer         string `json:"refer"`
	} `json:"data"`
}

type BAYLog struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title"`
	BankType     string    `json:"bank_type"`
	BankNumber   string    `json:"bank_number"`
	Amount       float64   `json:"amount"`
	LatestAmount float64   `json:"latest_amount"`
	Datetime     string    `json:"datetime"`
	Status       string    `json:"status"`
	Refer        string    `json:"refer"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
