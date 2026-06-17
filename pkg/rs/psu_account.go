package rs

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"securebanking-test-data-initializer/pkg/common"
	"securebanking-test-data-initializer/pkg/httprest"
	"securebanking-test-data-initializer/pkg/types"
	"strconv"
)

// CreatePSU - create the psu user if necessary and always return the userId if exist to populate de user data into RS
func CreatePSU() string {
	exist, userId := identityExists(common.Config.Users.PsuUsername)
	if exist {
		zap.S().Infof("Account found, skipping creation of Payment Services User", "userID", userId)
		return userId
	}

	zap.S().Infof("Creating PSU (Payment Services User)")

	user := &PSU{
		UserId:    common.Config.Users.PsuUserId,
		UserName:  common.Config.Users.PsuUsername,
		SN:        "Payment Services User",
		GivenName: "PSU",
		Mail:      "psu@acme.com",
		Password:  common.Config.Users.PsuPassword,
	}
	// TODO: check the managed user object, it's different for cloud
	var managedUserObject = "user"
	if common.Config.Environment.CloudType == types.Platform.Instance().FIDC {
		managedUserObject = common.Config.Identity.AmRealm + "_user"
	}
	path := "/openidm/managed/" + managedUserObject + "/?_action=create"
	body, s := httprest.Client.Post(path, user, map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/json",
		"Connection":   "keep-alive",
	})
	userRes := &UserResponse{}
	err := json.Unmarshal(body, userRes)
	if err != nil {
		panic(err)
	}
	zap.S().Debugw("PSU user created", "Response", userRes, "UserId", userRes.UserId)

	zap.S().Infow("PSU user created", "statusCode", s)
	return userRes.UserId
}

// PSUIdentityExists will check for psu identities in the alpha realm
func identityExists(identity string) (bool, string) {
	filter := "?_queryFilter=uid+eq+%22" + url.QueryEscape(identity) + "%22&_fields=username"
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/users" + filter
	serviceIdentityFilter := &types.ResultFilter{}
	b, _ := httprest.Client.Get(path, map[string]string{
		"Accept":             "application/json",
		"X-Requested-With":   "ForgeRock Identity Cloud Postman Collection",
		"Accept-Api-Version": "protocol=2.1, resource=4.0",
	})

	err := json.Unmarshal(b, serviceIdentityFilter)
	if err != nil {
		panic(err)
	}
	var psuID = ""
	if serviceIdentityFilter.ResultCount > 0 {
		psuID = serviceIdentityFilter.Result[0].ID
	}
	return serviceIdentityFilter.ResultCount > 0, psuID
}

// PopulateRSData -
func PopulateRSData(userId string) {
	if userId == "" {
		zap.L().Info("Skipping populate PSU Data to RS service")
		return
	}

	path := common.Config.Hosts.RsBaseUri + "/admin/data/user/has-data?userId=" + userId
	if mustPopulateUserData(path) {
		zap.S().Infow("Populate with RS Data the Payment Services User with the userId: " + userId)
		params := "userId=" + userId + "&username=" + common.Config.Users.PsuUsername + "&profile=random"
		path := common.Config.Hosts.RsBaseUri + "/admin/fake-data/generate?" + params
		s := httprest.Client.PostRS(path, map[string]string{
			"Accept":     "*/*",
			"Connection": "keep-alive",
		})
		zap.S().Infow("Populate RS Data response", "path", path, "statusCode", s)
	}
	//}
}

// mustPopulateUserData check if the user has data and if the environment is initialised, return true/false
func mustPopulateUserData(path string) bool {
	b, state := httprest.Client.GetRS(path, map[string]string{
		"Accept": "*/*",
	})
	if state != http.StatusOK {
		zap.S().Infow("No environment initialised", "path", path, "request status", state)
		return false
	}
	value := string(b)
	zap.S().Infow("User has data?", "path", path, "result", value)
	result, err := strconv.ParseBool(value)
	if err != nil {
		panic(err.Error())
	}
	return !result
}
