package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"net/http"
	"os"
)

type SenderMail struct {
	Mail string `json:"mail"`
}

type AdminMail struct {
	Mail      string `json:"mail"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func handleSendMailToEmployee(c *gin.Context) {
	var u SenderMail
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	}
	sendMail(u)
	c.JSON(http.StatusOK, u)

}

func handleSendMailToAdmin(c *gin.Context) {
	var u AdminMail
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	}
	SendMailtoAdmin(u)
	c.JSON(http.StatusOK, u)

}

func SendMailtoAdmin(u AdminMail) {
	adminMail := os.Getenv("ADMIN_MAIL")
	adminName := os.Getenv("ADMIN_NAME")
	m := gomail.NewMessage()
	m.SetHeader("From", "karthick.pannerselvam@humetis.in")
	m.SetHeader("To", adminMail)
	//m.SetHeader("To", "karthick.pannerselvam@humetis.in")
	m.SetHeader("Subject", "Waiting for your approval")
	m.SetBody("text/html", fmt.Sprintf(`<html><p><strong>Dear %s,</strong></p>
<p>&nbsp;</p>
<p>Waiting for your approval. Candidate Name %s and Mail Id %s</p>
<p>&nbsp;</p>
<br />Any questions, please let us know!</p></html>`, adminName, u.FirstName+" "+u.LastName, u.Mail))

	d := gomail.NewPlainDialer("secure.emailsrvr.com", 587, "karthick.pannerselvam@humetis.in", "Humetisindia@2020")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func sendMail(u SenderMail) {

	emp, err := GetEmployee(u.Mail)
	if err != nil {
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "karthick.pannerselvam@humetis.in")
	m.SetHeader("To", u.Mail)
	//m.SetHeader("To", "karthick.pannerselvam@humetis.in")
	m.SetHeader("Subject", "Employee Approval")
	m.SetBody("text/html", fmt.Sprintf(`<html><p><strong>Dear %s,</strong></p>
<p>&nbsp;</p>
<p>Your offer will be $%s salary annual gross. This includes 20 days PTO (20 days PTO including holidays, vacation, sick time). In practical terms, this would mean an average of 160 hours a month (1920 hours a year = 2080 potential hours a year minus 20 days/160 hours PTO)</p>
<p>&nbsp;</p>
<p>In line with our percentage of billing spend/budget commitment on our employee&rsquo;s billing for billable employees, our commitment will be Buget1% for Project 1 (On your project with %s client) and Budget 2% on project 2 (On your project with %s client). What this means to you is that the agreed percentage of billing as our budget commitment is tracked and spent on your salary, benefits, and expenses we incur on your behalf and expenses you can get reimbursed.</p>
<p>&nbsp;</p>
<p>Please note, our salary offer is our recommendation based on your budgets. If you want to change the offer within parameters of your budgets, LCA pay, etc., you can discuss it with us.</p>
<p>&nbsp;</p>
<p>How does it work and what can you expect</p>
<p>&nbsp;</p>
<p>You will get paid a semi-monthly salary for 80 hours every semi-monthly period and on an ongoing basis, you will receive details on exact budgets, spend, and $%s left in the budget from us. We will send this to you in the short term, but you will have it available online very soon to access at any point in time.</p>
<p>&nbsp;</p>
<p>There are certain benefits, support, and programs we cover with our margins and these are on us<br />Our team and our support &ndash; HR, Immigration, Recruiting, marketing, sales, admin and compliance Group life insurance Various training programs we will provide access with no cost from time to time</p>
<ul>
<li>401k profit sharing or matching we intend to provide in the future (most likely from January 2021)</li>
<li>Compliance</li>
<li>Payroll processing, tax filings</li>
<li>Non-immigrant employees - Visa and immigration compliance except premium fees and I-1485 stage for GC processing</li>
<li>Any agreed referral commissions and commissions for any coordination or support you provide to our corporate</li>
<li>Approved expenses you can get reimbursed from your budgets</li>
<li>Travel and incidental exp-Payroll processing, tax filing senses</li>
<li>Phone, internet, computer, etc. expenses</li>
<li>Any training/certifications/education (Tuition) expenses based on what you pursue</li>
<li>Any medical expenses not covered by our medical, dental, vision insurance</li>
<li>Public transit, mileage for work-related travel, parking expenses</li>
<li>Expenses we spend on your behalf from your budgets</li>
<li>Any relocation &ndash; moving, settling, flight charges, lodging, incidental</li>
<li>Non-immigrant visa holders - Premium fees and I-485 expenses</li>
<li>Any paid training you want us to pay for</li>
<li>Any company spend on your medical (company contribution)</li>
<li>Salary advances</li>
<li>Payroll during non-billable hours or insufficient hours</li>
<li>Any other legitimate expense we incur on your behalf</li>
<li>Additional advantages that do not meet your eyes explicitly</li>
<li>Based on your billing if the client approves, there is no limitation on your PTO</li>
<li>Stability between projects and we do provide access to learning and internal projects while</li>
<li>between billing projects, so you can call your employment with Zero bench period</li>
<li>You have an opportunity to maximize your pre-tax income</li>
<li>Ability to work and get compensated on more than 1 billing project</li>
</ul>
<p>Please note &ndash; in a loose sense most large corporations operate this way &ndash; they set budgets - They typically pay a low pay rate, bill clients high and spend a portion of their profits on benefits,<br />travel, bench pay, etc, &ndash; but their budgets are much lower percentage and they will miss flexibility and transparency. This model of Ours is an effort to bring stability, flexibility, remove the bench, maximize pre-tax benefits and provide opportunities for learning and provide a progressive career to work with us as a team instead of just a billing and payroll relationship.<br /> <br /><strong>Welcome to the family!</strong><br /> <br />Any questions, please let us know!</p></html>`, emp.FirstName+" "+emp.LastName, emp.ChooseYourOwnOffer, emp.Client1Name, emp.Client2Name, emp.ExpBonusBudget))
	m.Attach("template/report.pdf")

	d := gomail.NewPlainDialer("secure.emailsrvr.com", 587, "karthick.pannerselvam@humetis.in", "Humetisindia@2020")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
