package utils

func SwapBankToCode(code string) string {
	bankCode := "0" + code
	switch bankCode {
	case "001":
		return "BOT"
	case "002":
		return "BBL"
	case "004":
		return "KBANK"
	case "006":
		return "KTB"
	case "011":
		return "TMB"
	case "014":
		return "SCB"
	case "024":
		return "UOBT"
	case "025":
		return "BAY"
	case "030":
		return "GOV"
	case "031":
		return "HSBC"
	case "032":
		return "DEUTSCHE"
	case "033":
		return "GHB"
	case "039":
		return "MHCB"
	default:
		return "OTHER BANK"
	}
}
