package infisicalclient

import (
	"fmt"
	"os"
)

func (client Client) UniversalMachineIdentityAuth() (string, error) {
	if client.Config.ClientId == "" || client.Config.ClientSecret == "" {
		return "", fmt.Errorf("you must set the client secret and client ID for the client before making calls")
	}

	var loginResponse MachineIdentityAuthResponse

	res, err := client.Config.HttpClient.R().SetResult(&loginResponse).SetHeader("User-Agent", USER_AGENT).SetBody(map[string]string{
		"clientId":     client.Config.ClientId,
		"clientSecret": client.Config.ClientSecret,
	}).Post("api/v1/auth/universal-auth/login")

	if err != nil {
		return "", fmt.Errorf("UniversalMachineIdentityAuth: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return "", fmt.Errorf("UniversalMachineIdentityAuth: Unsuccessful response: [response=%s]", res)
	}

	return loginResponse.AccessToken, nil
}

func (client Client) GetServiceTokenDetailsV2() (GetServiceTokenDetailsResponse, error) {
	var tokenDetailsResponse GetServiceTokenDetailsResponse
	response, err := client.Config.HttpClient.
		R().
		SetResult(&tokenDetailsResponse).
		SetHeader("User-Agent", USER_AGENT).
		Get("api/v2/service-token")

	if err != nil {
		return GetServiceTokenDetailsResponse{}, fmt.Errorf("CallGetServiceTokenDetails: Unable to complete api request [err=%s]", err)
	}

	if response.IsError() {
		return GetServiceTokenDetailsResponse{}, fmt.Errorf("CallGetServiceTokenDetails: Unsuccessful response: [response=%s]", response)
	}

	return tokenDetailsResponse, nil
}

func (client Client) OidcMachineIdentityAuth() (string, error) {
	authJwt := os.Getenv(INFISICAL_AUTH_JWT_NAME)

	if client.Config.IdentityId == "" {
		return "", fmt.Errorf("you must set the identity ID for the client before making calls")
	}

	if authJwt == "" {
		return "", fmt.Errorf("%s is not present in the environment", INFISICAL_AUTH_JWT_NAME)
	}

	var loginResponse MachineIdentityAuthResponse

	res, err := client.Config.HttpClient.R().SetResult(&loginResponse).SetHeader("User-Agent", USER_AGENT).SetBody(map[string]string{
		"identityId": client.Config.IdentityId,
		"jwt":        authJwt,
	}).Post("api/v1/auth/oidc-auth/login")

	if err != nil {
		return "", fmt.Errorf("OidcMachineIdentityAuth: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return "", fmt.Errorf("OidcMachineIdentityAuth: Unsuccessful response: [response=%s]", res)
	}

	return loginResponse.AccessToken, nil
}
