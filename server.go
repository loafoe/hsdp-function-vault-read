package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type Client struct {
	*vault.Client
	Endpoint           string `json:"endpoint"`
	OrgSecretPath      string `json:"org_secret_path"`
	RoleID             string `json:"role_id"`
	SecretID           string `json:"secret_id"`
	ServiceSecretPath  string `json:"service_secret_path"`
	ServiceTransitPath string `json:"service_transit_path"`
	SpaceSecretPath    string `json:"space_secret_path"`
}

func stripFirst(path string) string {
	comps := strings.Split(path, "/")
	if len(comps) > 2 {
		return "/" + strings.Join(comps[2:], "/")
	}
	return ""
}

func (v *Client) ReadSpaceSecret(path string) (*vault.Secret, error) {
	return v.Client.Logical().Read(stripFirst(v.SpaceSecretPath) + "/" + path)
}

func (v *Client) ReadServiceSecret(path string) (*vault.Secret, error) {
	return v.Client.Logical().Read(stripFirst(v.ServiceSecretPath) + "/" + path)
}

func (v *Client) ReadServiceTransit(path string) (*vault.Secret, error) {
	return v.Client.Logical().Read(stripFirst(v.ServiceTransitPath) + "/" + path)
}

func main() {
	var err error
	var secret *vault.Secret

	viper.SetEnvPrefix("vread")
	viper.SetDefault("endpoint", "")
	viper.SetDefault("role_id", "")
	viper.SetDefault("secret_id", "")
	viper.SetDefault("org_secret_path", "")
	viper.SetDefault("space_secret_path", "")
	viper.SetDefault("service_secret_path", "")
	viper.SetDefault("service_transit_path", "")

	viper.AutomaticEnv()

	var client = Client{
		Endpoint:           viper.GetString("endpoint"),
		RoleID:             viper.GetString("role_id"),
		SecretID:           viper.GetString("secret_id"),
		OrgSecretPath:      viper.GetString("org_secret_path"),
		SpaceSecretPath:    viper.GetString("space_secret_path"),
		ServiceSecretPath:  viper.GetString("service_secret_path"),
		ServiceTransitPath: viper.GetString("service_transit_path"),
	}
	client.Client, err = vault.NewClient(vault.DefaultConfig())
	if err != nil {
		fmt.Printf("error creating vault client: %v\n", err)
		return
	}
	_ = client.SetAddress(client.Endpoint)
	secret, err = client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   client.RoleID,
		"secret_id": client.SecretID,
	})
	if err != nil {
		fmt.Printf("error authenticating against vault: %v\n", err)
		return
	}
	if secret == nil {
		fmt.Printf("login failed\n")
		return
	}

	client.SetToken(secret.Auth.ClientToken)
	secret, _ = client.Auth().Token().RenewSelf(1800)

	listenString := ":8080"

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Any("/vault/read/:namespace/:key", vaultReader(client))

	if port := os.Getenv("PORT"); port != "" {
		listenString = ":" + port
	}

	e.Logger.Fatal(e.Start(listenString))
}

func vaultReader(client Client) func(ctx echo.Context) error {
	return func(c echo.Context) error {
		var err error
		var secret *vault.Secret

		namespace := c.Param("namespace")
		key := c.Param("key")
		switch namespace {
		case "space":
			secret, err = client.ReadSpaceSecret(key)
		case "service":
			secret, err = client.ReadServiceSecret(key)
		case "transit":
			secret, err = client.ReadServiceTransit(key)
		default:
			err = fmt.Errorf("unsupported namespace: %s", namespace)
		}
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": err.Error()})
		}
		if secret != nil {
			return c.JSON(http.StatusOK, secret.Data)
		}
		return c.JSON(http.StatusNoContent, "{}")
	}
}
