package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

type Employee struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	DOB        string `json:"dob"`
	Email      string `json:"email"`
	VisaStatus string `json:"visaStatus"`
	VisaType   string `json:"visaType"`
	Status     string `json:"status"`
}

type EmployeeFullDetail struct {
	FirstName                      string  `json:"firstName"`
	LastName                       string  `json:"lastName"`
	DOB                            string  `json:"dob"`
	Email                          string  `json:"email"`
	Status                         string  `json:"status"`
	VisaStatus                     string  `json:"visaStatus"`
	VisaType                       string  `json:"visaType"`
	VisaExpiryDate                 string  `json:"visaExpiryDate"`
	Client1Name                    string  `json:"client1Name"`
	Project1                       string  `json:"project1"`
	StartDate1                     string  `json:"startDate1"`
	HrsPerMonth1                   float64 `json:"hrsPerMonth1"`
	Location1                      string  `json:"location1"`
	Lca1                           float64 `json:"lca1"`
	ClientBillingRate1             float64 `json:"clientBillingRate1"`
	Vms1                           float64 `json:"vms1"`
	Dor1                           float64 `json:"dor1"`
	Load1                          float64 `json:"load1"`
	EffectiveBillingRateHumetis1   float64 `json:"effectiveBillingRateHumetis1"`
	PayConfirmed1                  float64 `json:"payConfirmed1"`
	Client2Name                    string  `json:"client2Name"`
	Project2                       string  `json:"project2"`
	StartDate2                     string  `json:"startDate2"`
	HrsPerMonth2                   float64 `json:"hrsPerMonth2"`
	Location2                      string  `json:"location2"`
	Lca2                           float64 `json:"lca2"`
	ClientBillingRate2             float64 `json:"clientBillingRate2"`
	Vms2                           float64 `json:"vms2"`
	Dor2                           float64 `json:"dor2"`
	Load2                          float64 `json:"load2"`
	EffectiveBillingRateHumetis2   float64 `json:"effectiveBillingRateHumetis2"`
	PayConfirmed2                  float64 `json:"payConfirmed2"`
	YearlyEstimatedBudget          float64 `json:"yearlyEstimatedBudget"`
	GroupDiscountedMedicalRequired bool    `json:"groupDiscountedMedicalRequired"`
	GcPw                           float64 `json:"gcPw"`
	LcaConsidered                  float64 `json:"lcaConsidered"`
	ExpBonusPay                    float64 `json:"expBonusPay"`
	ExpBonusBudget                 float64 `json:"expBonusBudget"`
	ExpBonusLca                    float64 `json:"expBonusLca"`
	ChooseYourOwnOffer             float64 `json:"chooseYourOwnOffer"`
	Exception                      bool    `json:"exception"`
	DateOfJoiningHumetis           string  `json:"dateOfJoiningHumetis"`
	Flexible                       bool    `json:"flexible"`
}

type RequestEmployee struct {
	Email string `json:"email"`
}

type EmployeeChangeStatus struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}

//Handler /api/employees
func handleGetEmployees(c *gin.Context) {
	var loadedTasks, err = GetAllEmployee()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"employees": loadedTasks})
}

func handleGetEmployee(c *gin.Context) {
	var u RequestEmployee
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	var employee, err = GetEmployee(u.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func handleUpdateEmployee(c *gin.Context) {
	var u EmployeeFullDetail
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	var err = UpdateEmployee(&u)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "updated"})
}

func handleAddEmployee(c *gin.Context) {
	var u EmployeeFullDetail
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	var err = AddEmployee(&u)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Added"})
}

func handleChangeEmployeeStatus(c *gin.Context) {
	var u EmployeeChangeStatus
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	emp, err := ChangeEmployeeStatus(&u)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, emp)
}

// GetAllTasks Retrieves all Employees from the db
func GetAllEmployee() ([]*Employee, error) {
	var employees []*Employee

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("humetis")
	collection := db.Collection("employees")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &employees)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return employees, nil
}

func GetEmployee(email string) (*EmployeeFullDetail, error) {
	var employee *EmployeeFullDetail
	fmt.Print(email)
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("humetis")
	collection := db.Collection("employees")
	result := collection.FindOne(ctx, bson.M{"email": email})
	if result == nil {
		return nil, errors.New("Could not find a user")
	}
	err := result.Decode(&employee)

	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	log.Printf("Users: %v", employee)
	return employee, nil
}

func DeleteEmployee(employee *EmployeeFullDetail) error {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	_, err := client.Database("humetis").Collection("employees").DeleteOne(ctx, bson.M{"email": employee.Email})
	if err != nil {
		log.Printf("Could not delete employee: %v", err)
		return err
	}
	return nil
}

func InsertEmployee(employee *EmployeeFullDetail) error {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	_, err := client.Database("humetis").Collection("employees").InsertOne(ctx, employee)
	if err != nil {
		log.Printf("Could not insert employee: %v", err)
		return err
	}
	return nil
}

func UpdateEmployee(employee *EmployeeFullDetail) error {
	err := DeleteEmployee(employee)
	if err != nil {
		return err
	}
	err = InsertEmployee(employee)
	if err != nil {
		return err
	}
	return nil
}

func AddEmployee(employee *EmployeeFullDetail) error {
	_, err := GetEmployee(employee.Email)
	if err == nil {
		return errors.New("employee: already exist")
	}
	err = InsertEmployee(employee)
	if err != nil {
		return err
	}
	return nil
}

func ChangeEmployeeStatus(employee *EmployeeChangeStatus) (*EmployeeFullDetail, error) {
	emp, err := GetEmployee(employee.Email)
	if err != nil {
		return nil, errors.New("employee: Not exist")
	}
	err = DeleteEmployee(emp)
	if err != nil {
		return nil, err
	}
	emp.Status = employee.Status
	err = InsertEmployee(emp)
	if err != nil {
		return nil, err
	}
	return emp, nil
}
