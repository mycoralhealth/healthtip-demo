package healthtip

// User is stored in DB
type User struct {
	Id        int    `json:"id,string,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Password  string `json:"password,omitempty"`
	LastTip   int64  `json:"lastTipEpoch,omitempty"`
}

// AuthToken is the session token
type AuthToken struct {
	ApiUser int    `json:"apiUser,string,omitempty"`
	ApiKey  string `json:"apiKey,omitempty"`
}

// Record is the lab test results of the user
type Record struct {
	Id                  int  `json:"id,string,omitempty"`
	UserId              int  `json:"userId,string,omitempty"`
	Age                 int  `json:"age,string,omitempty"`
	Height              int  `json:"height,string,omitempty"`
	Weight              int  `json:"weight,string,omitempty"`
	Cholesterol         int  `json:"cholesterol,string,omitempty"`
	BloodPressure       int  `json:"bloodPressure,string,omitempty"`
	NumberOfCysts       int  `json:"numberOfCysts,string,omitempty"`
	Baldness            bool `json:"baldness,omitempty"`
	BaldnessFromDisease bool `json:"baldnessFromDisease,omitempty"`
	TipSent             int  `json:"tipSent,string,omitempty"`
}

// InsuranceCompany is a data representation of an insurance
// company.
type InsuranceCompany struct {
	Id   int    `json:"id,string,omitempty"`
	Name string `json:"name,omitempty"`
}

// Procedure is a medical procedure that a user can request.
type Procedure struct {
	Id   int    `json:"id,string,omitempty"`
	Name string `json:"name,omitempty"`
}

// InsuranceApprovalRequest is request information required for an insurance
// approval.
type InsuranceApprovalRequest struct {
	Procedure Procedure        `json:"procedure"`
	Company   InsuranceCompany `json:"company"`
}

type InsuranceApprovalResponse struct {
	InsuranceApprovalRequest
	Approved bool `json:"approved"`
}

// LoginResult is the model of an authorized login
type LoginResult struct {
	Token     AuthToken `json:"token,omitempty"`
	Email     string    `json:"email,omitempty"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
}
