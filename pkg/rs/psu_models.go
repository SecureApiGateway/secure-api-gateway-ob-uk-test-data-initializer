package rs

type PSU struct {
	UserId    string `json:"_id"`
	UserName  string `json:"userName"`
	SN        string `json:"sn"`
	GivenName string `json:"givenName"`
	Mail      string `json:"mail"`
	Password  string `json:"password"`
}

type STUDIO struct {
	UserId    string `json:"_id"`
	UserName  string `json:"userName"`
	SN        string `json:"sn"`
	GivenName string `json:"givenName"`
	Mail      string `json:"mail"`
	Password  string `json:"password"`
}

type UserResponse struct {
	UserId        string `json:"_id"`
	Rev           string `json:"_rev"`
	UserName      string `json:"userName"`
	AccountStatus string `json:"accountStatus"`
	GivenName     string `json:"givenName"`
	Sn            string `json:"sn"`
	Mail          string `json:"mail"`
}
