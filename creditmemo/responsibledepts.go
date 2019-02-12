package creditmemo
//
//import (
//	"crypto/tls"
//	"encoding/json"
//	"fmt"
//	"github.rackspace.com/SegmentSupport/raxss/internal/helpanto"
//	"io/ioutil"
//	"net/http"
//	"os"
//	"strings"
//	"time"
//)
//
//var baseUrl = "https://act-api.gscs.rackspace.com/v1/"
//
///*
//{
//   "responsibleDepartments": [
//       {
//           "responsibleDepartment": {
//               "id": 18,
//               "name": "Accounting",
//               "status": "active",
//               "workflowType": "Support",
//               "approverLdapGroup": "",
//               "category": "general"
//           }
//       },
//       {
//           "responsibleDepartment": {
//               "id": 1,
//               "name": "Account Management",
//               "status": "active",
//               "workflowType": "Support",
//               "approverLdapGroup": "",
//               "category": "general"
//           }
//       },
//*/
//
//type CreditMemoResponse struct {
//	Links struct {
//		Next struct {
//			Href string `json:"href"`
//		} `json:"next"`
//		Previous struct {
//			Href string `json:"href"`
//		} `json:"previous"`
//		CsvLink struct {
//			Href string `json:"href"`
//		} `json:"csvLink"`
//	} `json:"_links"`
//	Total    int `json:"total"`
//	Requests []struct {
//		Request struct {
//			ID                int         `json:"id"`
//			Amount            float64     `json:"amount"`
//			Type              string      `json:"type"`
//			Team              string      `json:"team"`
//			Calculation       string      `json:"calculation"`
//			OperatingUnit     interface{} `json:"operatingUnit"`
//			ContractingEntity struct {
//				ID          interface{} `json:"id"`
//				Code        interface{} `json:"code"`
//				Description interface{} `json:"description"`
//				DisplayName interface{} `json:"displayName"`
//			} `json:"contractingEntity"`
//			Currency struct {
//				Code string `json:"code"`
//			} `json:"currency"`
//			Initiator struct {
//				Sso         string `json:"sso"`
//				DisplayName string `json:"displayName"`
//			} `json:"initiator"`
//			AccountManager struct {
//				Sso         string `json:"sso"`
//				DisplayName string `json:"displayName"`
//			} `json:"accountManager"`
//			CustomerName           string      `json:"customerName"`
//			Description            string      `json:"description"`
//			DurationDays           int         `json:"durationDays"`
//			IncidentDate           string      `json:"incidentDate"`
//			IncidentEventID        interface{} `json:"incidentEventId"`
//			Prevention             string      `json:"prevention"`
//			RackspaceAccountNumber string      `json:"rackspaceAccountNumber"`
//			ReasonCode             struct {
//				ID              int    `json:"id"`
//				Name            string `json:"name"`
//				TransactionType string `json:"transactionType"`
//				WorkflowType    string `json:"workflowType"`
//			} `json:"reasonCode"`
//			AdditionalAmount      float64 `json:"additionalAmount"`
//			ResponsibleDepartment struct {
//				ID   int    `json:"id"`
//				Name string `json:"name"`
//			} `json:"responsibleDepartment"`
//			ServiceFailure struct {
//				ID       int    `json:"id"`
//				Category string `json:"category"`
//				Name     string `json:"name"`
//			} `json:"serviceFailure"`
//			Status          string    `json:"status"`
//			WorkflowState   string    `json:"workflowState"`
//			CreatedDatetime time.Time `json:"createdDatetime"`
//			UpdatedDatetime time.Time `json:"updatedDatetime"`
//			ClosedDate      string    `json:"closedDate"`
//			LastApprovedBy  struct {
//				Sso         string `json:"sso"`
//				DisplayName string `json:"displayName"`
//			} `json:"lastApprovedBy"`
//			QaApprovedBy struct {
//				Sso         string `json:"sso"`
//				DisplayName string `json:"displayName"`
//			} `json:"qaApprovedBy"`
//			Links struct {
//				Self struct {
//					Href string `json:"href"`
//				} `json:"self"`
//				History struct {
//					Href string `json:"href"`
//				} `json:"history"`
//				Notes struct {
//					Href string `json:"href"`
//				} `json:"notes"`
//				UsageEvents struct {
//					Href string `json:"href"`
//				} `json:"usageEvents"`
//			} `json:"_links"`
//			Tickets []struct {
//				ID    string `json:"id"`
//				Type  string `json:"type"`
//				Links struct {
//					Self struct {
//						Href string `json:"href"`
//					} `json:"self"`
//				} `json:"_links,omitempty"`
//			} `json:"tickets"`
//			IsTicketPrivate bool `json:"isTicketPrivate"`
//			AccountType     struct {
//				ID   int    `json:"id"`
//				Name string `json:"name"`
//			} `json:"accountType"`
//			Segment       string `json:"segment"`
//			SubSegment    string `json:"subSegment"`
//			BusinessUnit  string `json:"businessUnit"`
//			EarnedRevenue bool   `json:"earnedRevenue"`
//			Devices       []struct {
//				ID string `json:"id"`
//			} `json:"devices"`
//			BillingLocation struct {
//				ID           int    `json:"id"`
//				Name         string `json:"name"`
//				BillingOrgID int    `json:"billingOrgId"`
//			} `json:"billingLocation"`
//		} `json:"request"`
//	} `json:"requests"`
//}
//
//type ActRespDeptResponse struct {
//	ResponsibleDepartments []struct {
//		ResponsibleDepartment struct {
//			ID                int    `json:"id"`
//			Name              string `json:"name"`
//			Status            string `json:"status"`
//			WorkflowType      string `json:"workflowType"`
//			ApproverLdapGroup string `json:"approverLdapGroup"`
//			Category          string `json:"category"`
//		} `json:"responsibleDepartment"`
//	} `json:"responsibleDepartments"`
//}
//
////
////
////type ActRespDeptResponse struct {
////	ResponsibleDepartments ResponsibleDepartments
////	}
////
////
////type ResponsibleDepartments []struct{
////
////	ResponsibleDepartment ResponsibleDepartment `json:"responsible_departments"`
////}
////
////type ResponsibleDepartment struct {
////ID                int    `json:"id"`
////Name              string `json:"name"`
////Status            string `json:"status"`
////WorkflowType      string `json:"workflowType"`
////ApproverLdapGroup string `json:"approverLdapGroup"`
////Category          string `json:"category"`
////}
////
//
////type FilteredDepartments []ResponsibleDepartment
//func GetCreditMemoRequests(token string, searchValues map[string]string, month int, year int) (*CreditMemoResponse, error) {
//	/*
//		Get a list of requests, optionally filtered by ticketId, assignedTo, contracting entity, incidentDateStart, incidentDateEnd, amountMin, amountMax, rackspaceAccountNumber, initiatorSso, currency, reasonCode and/or team, operating unit
//	*/
//	ebi := GetEbiConfigFile("../config/ebiConfig.json")
//
//	connstring := ebi.GenerateConnectionString()
//	var baseUrl = "https://act-api.gscs.rackspace.com/v1/requests"
//	first := true
//	for k, v := range searchValues {
//		if first {
//			baseUrl += fmt.Sprintf("?%s=%s", k, v)
//			first = false
//		} else {
//			baseUrl += fmt.Sprintf("&%s=%s", k, v)
//
//		}
//	}
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client := &http.Client{
//
//		Transport: tr,
//	}
//
//	var resp *http.Response
//	request, err := http.NewRequest("GET", baseUrl, nil)
//
//	if err != nil {
//		return nil, err
//
//	}
//
//	request.Header.Add("Content-Type", "application/json")
//	request.Header.Add("X-Auth-Token", token)
//	resp, err = client.Do(request)
//
//	if err != nil {
//		return nil, err
//
//	}
//
//	if resp.StatusCode > 299 || resp.StatusCode < 200 {
//		return nil, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
//	}
//
//	body, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		return nil, err
//	}
//
//	//fmt.Println(string(body))
//	var creditMemoResponse CreditMemoResponse
//	err = json.Unmarshal(body, &creditMemoResponse)
//	if err != nil {
//		return nil, err
//
//	}
//
//	var filteredCreditMemoRequets CreditMemoResponse
//
//	for _, cm := range creditMemoResponse.Requests {
//
//		layout := "2006-01-02 15:04:05 -0700 MST"
//		closedTime, err := time.Parse(layout, cm.Request.UpdatedDatetime.String())
//		if err != nil {
//			fmt.Println("Error time stuff")
//			continue
//		}
//		closedMonth := int(closedTime.Month())
//		closedYear := closedTime.Year()
//
//		if closedMonth == month && closedYear == year {
//
//			if cm.Request.Currency.Code != helpanto.USD.ToString() {
//
//				newAmount, err := helpanto.ConvertToUSD(connstring, cm.Request.Currency.Code, closedMonth, closedYear, cm.Request.Amount)
//				if err != nil {
//					return nil, err
//				}
//				cm.Request.Amount = *newAmount
//				cm.Request.Currency.Code = helpanto.USD.ToString()
//
//				if cm.Request.AdditionalAmount > 0 {
//					newAddAmount, err := helpanto.ConvertToUSD(connstring, cm.Request.Currency.Code, closedMonth, closedYear, cm.Request.Amount)
//					if err != nil {
//						return nil, err
//					}
//					cm.Request.AdditionalAmount = *newAddAmount
//				}
//			}
//
//			filteredCreditMemoRequets.Requests = append(filteredCreditMemoRequets.Requests, cm)
//
//		}
//
//	}
//	if filteredCreditMemoRequets.Requests != nil {
//		creditMemoResponse.Requests = nil
//		creditMemoResponse.Requests = filteredCreditMemoRequets.Requests
//	}
//
//	return &creditMemoResponse, nil
//}
//func GetResponsibleDepartments(token string, searchValue *string) (map[int]string, error) {
//	deptMap := make(map[int]string)
//
//	baseUrl += "responsible-departments?status=active"
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client := &http.Client{
//
//		Transport: tr,
//	}
//
//	var resp *http.Response
//	request, err := http.NewRequest("GET", baseUrl, nil)
//
//	if err != nil {
//		fmt.Println(err.Error())
//
//		return nil, err
//
//	}
//
//	request.Header.Add("Content-Type", "application/json")
//	request.Header.Add("X-Auth-Token", token)
//	resp, err = client.Do(request)
//
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, err
//
//	}
//
//	if resp.StatusCode > 299 || resp.StatusCode < 200 {
//		return nil, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
//	}
//
//	body, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		fmt.Println(err.Error())
//
//		return nil, err
//	}
//
//	var actRespDeptResponse ActRespDeptResponse
//	err = json.Unmarshal(body, &actRespDeptResponse)
//	if err != nil {
//		return nil, err
//
//	}
//
//	if searchValue != nil {
//		searchValueLower := strings.ToLower(*searchValue)
//		filteredDeptMap := make(map[int]string)
//
//		for _, item := range actRespDeptResponse.ResponsibleDepartments {
//			deptMap[item.ResponsibleDepartment.ID] = item.ResponsibleDepartment.Name
//
//			lowerDeptName := strings.ToLower(item.ResponsibleDepartment.Name)
//			if strings.Contains(lowerDeptName, searchValueLower) {
//				filteredDeptMap[item.ResponsibleDepartment.ID] = item.ResponsibleDepartment.Name
//			}
//
//		}
//
//		return filteredDeptMap, nil
//	}
//	return deptMap, nil
//
//}
//
//type EbiDatabase struct {
//	Name             string `json:"name"`
//	Port             int    `json:"port"`
//	IPAddress        string `json:"ip_address"`
//	DatabaseName     string `json:"database_name"`
//	DatabaseUser     string `json:"database_user"`
//	DatabasePassword string `json:"database_password"`
//}
//
//func (s *EbiDatabase) GenerateConnectionString() string {
//	return fmt.Sprintf("odbc:server=%s; port=%d; user id=%s;password=%s; database=%s;log=3;encrypt=false;TrustServerCertificate=true", s.IPAddress, s.Port, s.DatabaseUser, s.DatabasePassword, s.DatabaseName)
//}
//
//func GetEbiConfigFile(filename string) *EbiDatabase {
//	//	log.Info("In Get Consoles")
//	if filename == "" {
//		filename = "ebiConfig.json"
//	}
//	raw, err := ioutil.ReadFile(filename)
//	if err != nil {
//		os.Exit(1)
//	}
//
//	var ebi *EbiDatabase
//	err = json.Unmarshal(raw, &ebi)
//	if err != nil {
//		os.Exit(1)
//	}
//
//	return ebi
//}
