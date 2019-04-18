package act

import (
	"encoding/json"
	"fmt"
	"github.rackspace.com/SegmentSupport/sspkgs/coreods"
	"github.rackspace.com/SegmentSupport/sspkgs/identity"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type CreditRequests struct {
	Links struct {
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
		Previous struct {
			Href string `json:"href"`
		} `json:"previous"`
		CsvLink struct {
			Href string `json:"href"`
		} `json:"csvLink"`
	} `json:"_links"`
	Total    int      `json:"total"`
	Requests Requests `json:"requests"`
}

type Requests []struct {
	Request CreditRequest `json:"request"`
}

type CreditRequest struct {
	ID                int         `json:"id"`
	Amount            float64     `json:"amount"`
	Type              string      `json:"type"`
	Team              string      `json:"team"`
	Calculation       string      `json:"calculation"`
	OperatingUnit     interface{} `json:"operatingUnit"`
	ContractingEntity struct {
		ID          interface{} `json:"id"`
		Code        interface{} `json:"code"`
		Description interface{} `json:"description"`
		DisplayName interface{} `json:"displayName"`
	} `json:"contractingEntity"`
	Currency struct {
		Code string `json:"code"`
	} `json:"currency"`
	Initiator struct {
		Sso         string `json:"sso"`
		DisplayName string `json:"displayName"`
	} `json:"initiator"`
	AccountManager struct {
		Sso         string `json:"sso"`
		DisplayName string `json:"displayName"`
	} `json:"accountManager"`
	CustomerName           string      `json:"customerName"`
	Description            string      `json:"description"`
	DurationDays           int         `json:"durationDays"`
	IncidentDate           string      `json:"incidentDate"`
	IncidentEventID        interface{} `json:"incidentEventId"`
	Prevention             string      `json:"prevention"`
	RackspaceAccountNumber string      `json:"rackspaceAccountNumber"`
	ReasonCode             struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		TransactionType string `json:"transactionType"`
		WorkflowType    string `json:"workflowType"`
	} `json:"reasonCode"`
	AdditionalAmount      float64 `json:"additionalAmount"`
	ResponsibleDepartment struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"responsibleDepartment"`
	ServiceFailure struct {
		ID       int    `json:"id"`
		Category string `json:"category"`
		Name     string `json:"name"`
	} `json:"serviceFailure"`
	Status          string    `json:"status"`
	WorkflowState   string    `json:"workflowState"`
	CreatedDatetime time.Time `json:"createdDatetime"`
	UpdatedDatetime time.Time `json:"updatedDatetime"`
	ClosedDate      string    `json:"closedDate"`
	TotalAmount     float64   `json:"total_amount"`
	LastApprovedBy  struct {
		Sso         string `json:"sso"`
		DisplayName string `json:"displayName"`
	} `json:"lastApprovedBy"`
	QaApprovedBy struct {
		Sso         string `json:"sso"`
		DisplayName string `json:"displayName"`
	} `json:"qaApprovedBy"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		History struct {
			Href string `json:"href"`
		} `json:"history"`
		Notes struct {
			Href string `json:"href"`
		} `json:"notes"`
		UsageEvents struct {
			Href string `json:"href"`
		} `json:"usageEvents"`
	} `json:"_links"`
	Tickets []struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links,omitempty"`
	} `json:"tickets"`
	IsTicketPrivate bool `json:"isTicketPrivate"`
	AccountType     struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"accountType"`
	Segment       string `json:"segment"`
	SubSegment    string `json:"subSegment"`
	BusinessUnit  string `json:"businessUnit"`
	EarnedRevenue bool   `json:"earnedRevenue"`
	Devices       []struct {
		ID string `json:"id"`
	} `json:"devices"`
	BillingLocation struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		BillingOrgID int    `json:"billingOrgId"`
	} `json:"billingLocation"`
}

type ActRespDeptResponse struct {
	ResponsibleDepartments []struct {
		ResponsibleDepartment struct {
			ID                int          `json:"id"`
			Name              string       `json:"name"`
			Status            string       `json:"status"`
			WorkflowType      WorkflowType `json:"workflowType"`
			ApproverLdapGroup string       `json:"approverLdapGroup"`
			Category          string       `json:"category"`
		} `json:"responsibleDepartment"`
	} `json:"responsibleDepartments"`
}

type WorkflowType struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ACTClient interface {
	GetDepartments(filter *string) (map[int]string, error)
	GetCreditRequests(pageLink *string, queryParams *string) (*CreditRequests, error)
	GetRequest(query string) (*[]byte, error)
}

type ACTService interface {
	ACTClient
}

var _ ACTService = &actService{}

type actService struct {
	ACTClient
}

var _ ACTClient = &actClient{}

type actClient struct {
	client   *http.Client
	baseURL  string
	username string
	password string
	domain   string

	coreodsService coreods.ODSService

	idService identity.IdentityService
	logger    *log.Logger
	token     *string
}

func NewACTService(client *http.Client, baseURL string, username string, password string, domain string, coreodsService coreods.ODSService, idService identity.IdentityService, logger *log.Logger, token *string) ACTService {

	logger.WithFields(log.Fields{
		"func":    "NewACTService",
		"baseURL": baseURL,
	}).Info("")
	authReq := identity.Request{}
	authReq.Auth.RAXAUTHDomain.Name = domain
	authReq.Auth.PasswordCredentials.Username = username
	authReq.Auth.PasswordCredentials.Password = password

	resp, err := idService.Authenticate(authReq)

	if err != nil {

	}

	tempToken := resp.Access.Token.ID
	internalToken := &tempToken

	act := &actClient{client, baseURL, username, password, domain, coreodsService, idService, logger, internalToken}

	return &actService{ACTClient: act}
}
func (ac *actClient) GetRequest(query string) (*[]byte, error) {
	var resp *http.Response

	ac.logger.Info("get req query: ", query)

	reqURL := fmt.Sprintf("%s%s", ac.baseURL, query)
	ac.logger.Info("reqURL: ", reqURL)

	request, err := http.NewRequest("GET", reqURL, nil)

	if err != nil {
		return nil, err

	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Auth-Token", *ac.token)
	resp, err = ac.client.Do(request)
	if err != nil {

		return nil, err

	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		return nil, err
	}

	return &body, nil
}
func (ac *actClient) GetDepartments(filter *string) (map[int]string, error) {
	query := "/responsible-departments?status=active"
	var respDept ActRespDeptResponse

	body, err := ac.GetRequest(query)
	if err != nil {

		return nil, err

	}

	err = json.Unmarshal(*body, &respDept)
	if err != nil {

		return nil, err

	}
	deptMap := make(map[int]string)

	if filter != nil {
		for _, dept := range respDept.ResponsibleDepartments {
			if strings.Contains(dept.ResponsibleDepartment.Name, *filter) {
				deptMap[dept.ResponsibleDepartment.ID] = dept.ResponsibleDepartment.Name

			}
		}

	} else {
		for _, dept := range respDept.ResponsibleDepartments {
			deptMap[dept.ResponsibleDepartment.ID] = dept.ResponsibleDepartment.Name

		}
	}
	return deptMap, err
}

// GetResponsibleDepartments gets a list of requests, optionally filtered by ticketId, assignedTo, contracting entity, incidentDateStart, incidentDateEnd, amountMin, amountMax, rackspaceAccountNumber, initiatorSso, currency, reasonCode and/or team, operating unit

func (ac *actClient) GetCreditRequests(pageLink *string, queryParams *string) (*CreditRequests, error) {

	if pageLink == nil {

		pl := fmt.Sprintf("/requests")

		if queryParams != nil {

			pl = fmt.Sprintf("/requests%s", *queryParams)
			log.Info("PL: ", pl)
		}
		pageLink = &pl
	}

	var creditRequests CreditRequests

	body, err := ac.GetRequest(*pageLink)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(*body, &creditRequests)
	if err != nil {
		return nil, err
	}

	//log.Infof("body: %s", body)
	for _, cr := range creditRequests.Requests {
		//2018-04-09 14:46:49.748 +0000 UTC
		//2006-01-02T15:04:05Z07:00
		timeFormat := "2006-01-02"
		closeDate, err := time.Parse(timeFormat, cr.Request.ClosedDate)
		if err != nil {
			return nil, err
		}
		closedYear, closedMonth, _ := closeDate.Date()
		ac.logger.WithFields(log.Fields{
			"func":          "GetCreditRequests",
			"currency code": cr.Request.Currency.Code,
			"closed Month":  int(closedMonth),
			"closed Year":   closedYear,
		}).Infoln("Prior to calling for exchange rate")

		if cr.Request.AdditionalAmount >= 0.0 {
			fmt.Printf("additional >=0.0:  %v + %v = %v\n", cr.Request.Amount, cr.Request.AdditionalAmount, cr.Request.Amount+cr.Request.AdditionalAmount)
			cr.Request.TotalAmount = cr.Request.AdditionalAmount + cr.Request.Amount
		} else {
			fmt.Printf("additional < 0.0:   %v + %v = %v\n", cr.Request.Amount, cr.Request.AdditionalAmount, cr.Request.Amount+cr.Request.AdditionalAmount)
			cr.Request.TotalAmount = cr.Request.Amount
		}

		if cr.Request.Currency.Code != "USD" {
			reqAmount, err := ac.coreodsService.ExchangeRate(cr.Request.Currency.Code, int(closedMonth), closedYear)
			if err != nil {
				ac.logger.WithFields(log.Fields{
					"func":          "GetCreditRequests",
					"currency code": cr.Request.Currency.Code,
					"closed Month":  int(closedMonth),
					"closed Year":   closedYear,
					"error":         err.Error(),
				}).Warnln("error getting exchange rate")
				return nil, err
			}
			currCode := "USD"
			cr.Request.Amount = reqAmount * cr.Request.Amount

			cr.Request.Currency.Code = currCode
		}

	}
	return &creditRequests, nil
}

func getIDForDC(deptMap map[int]string, n string) (int, error) {

	// Filter for dc name passed in
	for k, v := range deptMap {
		words := strings.Split(v, " ")
		fdc := words[len(words)-1]

		if fdc == n {
			return k, nil
		}

		if len(n) > 3 {

			if fdc == n[:3] {

				return k, nil
			}
		}
	}
	return 0, fmt.Errorf("DC Not Found: %s", n)
}
