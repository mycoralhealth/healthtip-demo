package healthtip

// User is stored in DB
type User struct {
	ID         int    `json:"id,string,omitempty"`
	Email      string `json:"email,omitempty"`
	First_name string `json:"firstName,omitempty"`
	Last_name  string `json:"lastName,omitempty"`
	Password   string `json:"password,omitempty"`
	Last_tip   int64  `json:"lastTipEpoch,omitempty"`
}

// AuthToken is the session token
type AuthToken struct {
	Api_user int    `json:"apiUser,string,omitempty"`
	Api_key  string `json:"apiKey,omitempty"`
}

// Record is the lab test results of the user
type Record struct {
	ID                    int  `json:"id,string,omitempty"`
	User_id               int  `json:"userId,string,omitempty"`
	Age                   int  `json:"age,string,omitempty"`
	Height                int  `json:"height,string,omitempty"`
	Weight                int  `json:"weight,string,omitempty"`
	Cholesterol           int  `json:"cholesterol,string,omitempty"`
	Blood_pressure        int  `json:"bloodPressure,string,omitempty"`
	Number_of_cysts       int  `json:"numberOfCysts,string,omitempty"`
	Baldness              bool `json:"baldness,omitempty"`
	Baldness_from_disease bool `json:"baldnessFromDisease,omitempty"`
	Tip_sent              int  `json:"tipSent,string,omitempty"`
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
	Procedure string `json:"procedure"`
	Company   string `json:"company"`
}

type InsuranceApprovalResponse struct {
	InsuranceApprovalRequest
	Approved bool `json:"approved"`
}

// LoginResult is the model of an authorized login
type LoginResult struct {
	Token      AuthToken `json:"token,omitempty"`
	Email      string    `json:"email,omitempty"`
	First_name string    `json:"firstName,omitempty"`
	Last_name  string    `json:"lastName,omitempty"`
}
