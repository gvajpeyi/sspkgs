 package identity
//
 import (
	 "crypto/tls"
	 "fmt"
	 "github.com/sirupsen/logrus"
	 "net/http"
	 "os"
	 "testing"
	 "time"
 )



func setupTests() (IdentityService, error) {
	ssid := os.Getenv("SSID")
	sspw := os.Getenv("SSPW")

	tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{

			Transport: tr,
		}
		
		
	var logger = logrus.New()

	logger.Out = os.Stdout

	logger.Formatter = new(logrus.TextFormatter) //default
	logger.SetOutput(logger.Writer())

	ll := logrus.DebugLevel

	logger.SetLevel(ll)

	logger.Out = os.Stdout
		
	if ssid == "" || sspw == "" {
		return nil, fmt.Errorf("identity env variables not configured")
	}
	
	id := NewIdentityService(client, "https://identity-internal.api.rackspacecloud.com/v2.0", logger)






	
	
	
	return id, nil
	
}

func TestAuthenticate(t *testing.T) {
	idsrv, _ :=setupTests()
	testCases := []struct {
		user  string
		pass  string
		token string
		want  bool
	}{
		{user: os.Getenv("SSID"), pass: os.Getenv("SSPW"),token: "AAAbVMKVFzPgneC5Yy4p3IAicxDNx__nq9e7snYeLRQavGiryAICCjJDBH9b_d3bsw0povu88c5YCaIf9KzhM7fW3qABttiONaqBaAOjccBPVDkBjMSlU8AK", want: true},
	//	{user: "SSID",pass: "SSPW", token: "", want: false},

	}
	

	for _, tc := range testCases {
		authReq := Request{}
		authReq.Auth.RAXAUTHDomain.Name = "Rackspace"
		authReq.Auth.PasswordCredentials.Username = tc.user
		authReq.Auth.PasswordCredentials.Password = tc.pass

		got, err := idsrv.Authenticate(authReq)

		if err != nil {

			t.Fatalf("Want %v   Got  %v\n", tc.want, err)
		}

		fiveMinutesInFuture := time.Now().UTC().Add(time.Hour*time.Duration(0) +
			time.Minute*time.Duration(1) +
			time.Second*time.Duration(0))

		timeDiff := (fiveMinutesInFuture).Sub(got.Access.Token.RAXAUTHIssued.Add(1 * time.Minute))

		var validIssueTime= false
		if timeDiff <= (1 * time.Minute) {
			validIssueTime = true
		}
		if validIssueTime != tc.want {
			t.Errorf("Want %v   Got %v\n", tc.want, validIssueTime)
		}

	}
}

func TestVerifyToken(t *testing.T){

	idsrv, nil := setupTests()
	authReq := Request{}
	authReq.Auth.RAXAUTHDomain.Name = "Rackspace"
	authReq.Auth.PasswordCredentials.Username = os.Getenv("SSID")
	authReq.Auth.PasswordCredentials.Password = os.Getenv("SSPW")

	authToken := ""
	
	resp, err := idsrv.Authenticate(authReq)

	if err != nil {
		t.Error(err)
	}


		if resp.Access.Token.ID != authToken && authToken != ""{
			t.Errorf("want: %s     got: %s", authToken, resp.Access.Token.ID)
		} else {
		validToken := resp.Access.Token.ID

		testCases := []struct {
			arg  string
			want bool
		}{
			{arg: validToken, want: true},
			{arg: "not valid", want: false},
		}
		for _, tc := range testCases {

			got, _, err := idsrv.VerifyToken(tc.arg)
			if err != nil || got != tc.want {
				t.Errorf("Want: %t  %T Got: %T  %t\n", tc.want, tc.want, got, got)
			}
		}

	}

}