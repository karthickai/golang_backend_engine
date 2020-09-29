package main

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"math"
	"net/http"
)

func handleCalculate(c *gin.Context) {
	var u EmployeeFullDetail
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	var employee EmployeeFullDetail
	var err error
	if u.Flexible == true {
		employee, err = flexibleCalculate(u)
	} else {
		employee, err = normalCalculate(u)

	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func normalCalculate(emp EmployeeFullDetail) (EmployeeFullDetail, error) {
	emp.EffectiveBillingRateHumetis1 = floatRound((emp.ClientBillingRate1*(1-(emp.Vms1/100)))*(1-(emp.Dor1/100)) - emp.Load1)
	emp.EffectiveBillingRateHumetis2 = floatRound((emp.ClientBillingRate2*(1-(emp.Vms2/100)))*(1-(emp.Dor2/100)) - emp.Load2)

	if emp.Lca1 > emp.Lca2 {
		emp.LcaConsidered = emp.Lca1
	} else {
		emp.LcaConsidered = emp.Lca2

	}

	emp.YearlyEstimatedBudget = floatRound((emp.EffectiveBillingRateHumetis1 * (emp.PayConfirmed1 / 100) * emp.HrsPerMonth1 * 12) + (emp.EffectiveBillingRateHumetis2 * (emp.PayConfirmed2 / 100) * emp.HrsPerMonth2 * 12) - 4800)
	if emp.GroupDiscountedMedicalRequired {
		emp.ChooseYourOwnOffer = floatRound(0.9*emp.YearlyEstimatedBudget - 4800)
	} else {
		emp.ChooseYourOwnOffer = floatRound(0.9 * emp.YearlyEstimatedBudget)
	}

	emp.ExpBonusBudget = floatRound(((emp.YearlyEstimatedBudget - emp.ChooseYourOwnOffer) / emp.YearlyEstimatedBudget) * 100)
	emp.ExpBonusPay = floatRound(((emp.YearlyEstimatedBudget - emp.ChooseYourOwnOffer) / emp.ChooseYourOwnOffer) * 100)
	emp.ExpBonusLca = floatRound(((emp.YearlyEstimatedBudget - emp.ChooseYourOwnOffer) / emp.LcaConsidered) * 100)

	return emp, nil

}

func flexibleCalculate(emp EmployeeFullDetail) (EmployeeFullDetail, error) {
	if emp.GroupDiscountedMedicalRequired {
		emp.YearlyEstimatedBudget = floatRound((1/0.9)*emp.ChooseYourOwnOffer + 4800)
	} else {
		emp.YearlyEstimatedBudget = floatRound((1 / 0.9) * emp.ChooseYourOwnOffer)
	}

	if emp.ClientBillingRate2 > 0 {
		emp.EffectiveBillingRateHumetis2 = floatRound((emp.ClientBillingRate2*(1-(emp.Vms2/100)))*(1-(emp.Dor2/100)) - emp.Load2)
	} else {
		emp.EffectiveBillingRateHumetis2 = floatRound((emp.ClientBillingRate2 * (1 - (emp.Vms2 / 100))) * (1 - (emp.Dor2 / 100)))

	}

	emp.EffectiveBillingRateHumetis1 = floatRound((emp.YearlyEstimatedBudget - (emp.EffectiveBillingRateHumetis2 * (emp.PayConfirmed2 / 100) * emp.HrsPerMonth2 * 12)) / 160 / 12 / (emp.PayConfirmed1 / 100))
	emp.ClientBillingRate1 = floatRound((emp.EffectiveBillingRateHumetis1 * (1 / (1 - (emp.Vms1 / 100))) * (1 / (1 - (emp.Vms1 / 100)))) + emp.Load1)
	return emp, nil
}

type Report struct {
	From                           string `json:"from"`
	To                             string `json:"to"`
	HrsPayRollPeriod1              string `json:"hrs_pay_roll_period_1"`
	HrsPayRollPeriod2              string `json:"hrs_pay_roll_period_2"`
	MedicalCoverage                string `json:"medical_coverage"`
	MedicalCTC                     string `json:"medical_ctc"`
	BaseCorpRole                   string `json:"base_corp_role"`
	COMMBonus                      string `json:"comm_bonus"`
	TotalMoBudget                  string `json:"total_mo_budget"`
	Salary                         string `json:"salary"`
	AdvCompExp                     string `json:"adv_comp_exp"`
	ExpWD                          string `json:"exp_wd"`
	BonusWD                        string `json:"bonus_wd"`
	BudgetBalance                  string `json:"budget_balance"`
	NumberOfSelectedPayRollReserve string `json:"number_of_selected_pay_roll_reserve"`
	NumberOfLCAPayRoll             string `json:"number_of_lca_pay_roll"`
}

type Detail struct {
	EmployeeName string `json:"employee_name"`
	AnnualSalary string `json:"annual_salary"`
	LCASalary    string `json:"lca_salary"`
	P1Rate       string `json:"p_1_rate"`
	P1Percentage string `json:"p_1_percentage"`
	P2Rate       string `json:"p_2_rate"`
	P2Percentage string `json:"p_2_percentage"`
}

type Data struct {
	Reports []Report `json:"reports"`
	Detail  Detail   `json:"detail"`
}

func handleGenerateReport(c *gin.Context) {
	var u EmployeeFullDetail
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	}
	var err error
	data, err := generateReport(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, data)

}

func reportHelper(u EmployeeFullDetail, report Report, prevBudgetBalance float64) (Report, float64) {

	medicalCTC := 0.0
	if u.GroupDiscountedMedicalRequired {
		report.MedicalCoverage = "Yes"
		report.MedicalCTC = "200"
		medicalCTC = 200
	} else {
		report.MedicalCoverage = "No"
		report.MedicalCTC = "0"
	}
	tempTotalBudget := (u.ClientBillingRate1 * u.PayConfirmed1 * u.HrsPerMonth1) + (u.ClientBillingRate2 * u.PayConfirmed2 * u.HrsPerMonth2)
	report.TotalMoBudget = fmt.Sprintf("%.2f", tempTotalBudget)
	tempSalary := u.ChooseYourOwnOffer / 24
	report.Salary = fmt.Sprintf("%.2f", tempSalary)
	if prevBudgetBalance != 0 {
		prevBudgetBalance = tempTotalBudget - tempSalary - medicalCTC
		report.BudgetBalance = fmt.Sprintf("%.2f", prevBudgetBalance)
	} else {
		prevBudgetBalance = tempTotalBudget - tempSalary - medicalCTC
		report.BudgetBalance = fmt.Sprintf("%.2f", prevBudgetBalance)
	}
	report.NumberOfSelectedPayRollReserve = fmt.Sprintf("%.2f", prevBudgetBalance/tempSalary)
	report.NumberOfLCAPayRoll = fmt.Sprintf("%.2f", prevBudgetBalance/tempSalary)
	return report, prevBudgetBalance
}

func generateReport(u EmployeeFullDetail) (Data, error) {
	fmt.Println("Time parsing")
	dateString := "2020-09-10"
	totalMonths := 30
	var reports []Report
	var detail Detail
	t, err := dateparse.ParseLocal(dateString)
	if err != nil {
		fmt.Println("Error while parsing date :", err)
	}
	detail.EmployeeName = u.FirstName + " " + u.LastName
	detail.AnnualSalary = fmt.Sprintf("%.2f", u.ChooseYourOwnOffer)
	detail.LCASalary = fmt.Sprintf("%.2f", u.LcaConsidered)
	detail.P1Rate = fmt.Sprintf("%.2f", u.ClientBillingRate1)
	detail.P1Percentage = fmt.Sprintf("%.2f", u.PayConfirmed1)
	detail.P2Rate = fmt.Sprintf("%.2f", u.ClientBillingRate2)
	detail.P2Percentage = fmt.Sprintf("%.2f", u.PayConfirmed2)

	prevBudgetBalance := 0.0
	for i := 1; i <= totalMonths; i++ {
		var report Report
		//First Half Month
		from := now.With(t).BeginningOfMonth()
		report.From = from.Format("2006-01-02")
		to := from.AddDate(0, 0, 14)
		report.To = to.Format("2006-01-02")
		report.HrsPayRollPeriod1 = fmt.Sprintf("%.2f", u.HrsPerMonth1/2)
		report.HrsPayRollPeriod2 = fmt.Sprintf("%.2f", u.HrsPerMonth2/2)
		report, prevBudgetBalance = reportHelper(u, report, prevBudgetBalance)
		reports = append(reports, report)

		//Second Half Month
		report.From = to.AddDate(0, 0, 1).Format("2006-01-02")
		report.To = now.With(t).EndOfMonth().Format("2006-01-02")
		report.HrsPayRollPeriod1 = fmt.Sprintf("%.2f", u.HrsPerMonth1/2)
		report.HrsPayRollPeriod2 = fmt.Sprintf("%.2f", u.HrsPerMonth2/2)
		report, prevBudgetBalance = reportHelper(u, report, prevBudgetBalance)
		reports = append(reports, report)
		t = t.AddDate(0, 1, 0)
	}

	for _, rep := range reports {
		fmt.Printf("From : %s, To : %s\n", rep.From, rep.To)
	}
	//temp := template.Must(template.New("").Parse(`<table>{{range .}}<tr><td>{{.}}</td></tr>{{end}}</table>`))
	r := NewRequestPdf("")
	templatePath := "template/report.html"

	outputPath := "template/report.pdf"
	var data Data
	data.Detail = detail
	data.Reports = reports
	if err := r.ParseTemplate(templatePath, data); err == nil {
		ok, _ := r.GeneratePDF(outputPath)
		fmt.Println(ok, "pdf generated successfully")
	} else {
		fmt.Println(err)
		return data, err
	}

	return data, nil

}

func handleDownloadReport(c *gin.Context) {
	c.File("report.pdf")
}

func floatRound(x float64) float64 {
	return math.Round(x*100) / 100
}
