package healthtip

// Information about the current user.
type UserInfo struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// Database record indicating the last tip timestamp for a given user.
type Tips struct {
	UserId    string `json:"userId,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// Record is the lab test results of the user
type Record struct {
	Id                  int    `json:"id,string,omitempty"`
	UserId              string `json:"userId,omitempty"`
	Age                 int    `json:"age,string,omitempty"`
	Height              int    `json:"height,string,omitempty"`
	Weight              int    `json:"weight,string,omitempty"`
	Cholesterol         int    `json:"cholesterol,string,omitempty"`
	BloodPressure       int    `json:"bloodPressure,string,omitempty"`
	NumberOfCysts       int    `json:"numberOfCysts,string"`
	Baldness            bool   `json:"baldness,omitempty"`
	BaldnessFromDisease bool   `json:"baldnessFromDisease,omitempty"`
	TipSent             int    `json:"tipSent,string,omitempty"`
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
