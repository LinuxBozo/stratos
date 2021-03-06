package interfaces

import (
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

type V2Info struct {
	AuthorizationEndpoint    string `json:"authorization_endpoint"`
	TokenEndpoint            string `json:"token_endpoint"`
	DopplerLoggingEndpoint   string `json:"doppler_logging_endpoint"`
	AppSSHEndpoint           string `json:"app_ssh_endpoint"`
	AppSSHHostKeyFingerprint string `json:"app_ssh_host_key_fingerprint"`
	AppSSHOauthCLient        string `json:"app_ssh_oauth_client"`
}

type InfoFunc func(apiEndpoint string, skipSSLValidation bool) (CNSIRecord, interface{}, error)

//TODO this could be moved back to cnsis subpackage, and extensions could import it?
type CNSIRecord struct {
	GUID                   string   `json:"guid"`
	Name                   string   `json:"name"`
	CNSIType               string   `json:"cnsi_type"`
	APIEndpoint            *url.URL `json:"api_endpoint"`
	AuthorizationEndpoint  string   `json:"authorization_endpoint"`
	TokenEndpoint          string   `json:"token_endpoint"`
	DopplerLoggingEndpoint string   `json:"doppler_logging_endpoint"`
	SkipSSLValidation      bool     `json:"skip_ssl_validation"`
}

// ConnectedEndpoint
type ConnectedEndpoint struct {
	GUID                   string   `json:"guid"`
	Name                   string   `json:"name"`
	CNSIType               string   `json:"cnsi_type"`
	APIEndpoint            *url.URL `json:"api_endpoint"`
	Account                string   `json:"account"`
	TokenExpiry            int64    `json:"token_expiry"`
	DopplerLoggingEndpoint string   `json:"-"`
	SkipSSLValidation      bool     `json:"skip_ssl_validation"`
	TokenMetadata          string   `json:"-"`
}

const (
	AuthTypeOAuth2    = "OAuth2"
	AuthTypeOIDC      = "OIDC"
	AuthTypeHttpBasic = "HttpBasic"
)

const (
	AuthConnectTypeCreds = "creds"
)

// Token record for an endpoint (includes the Endpoint GUID)
type EndpointTokenRecord struct {
	*TokenRecord
	EndpointGUID    string
	EndpointType    string
	APIEndpint      string
	LoggingEndpoint string
}

//TODO this could be moved back to tokens subpackage, and extensions could import it?
type TokenRecord struct {
	AuthToken    string
	RefreshToken string
	TokenExpiry  int64
	Disconnected bool
	AuthType     string
	Metadata     string
}

type CFInfo struct {
	EndpointGUID string
	SpaceGUID    string
	AppGUID      string
}

// Structure for optional metadata for an OAuth2 Token
type OAuth2Metadata struct {
	ClientID     string
	ClientSecret string
	IssuerURL    string
}

type VCapApplicationData struct {
	API           string `json:"cf_api"`
	ApplicationID string `json:"application_id"`
	SpaceID       string `json:"space_id"`
}

type LoginRes struct {
	Account     string   `json:"account"`
	TokenExpiry int64    `json:"token_expiry"`
	APIEndpoint *url.URL `json:"api_endpoint"`
	Admin       bool     `json:"admin"`
}

type LoginHookFunc func(c echo.Context) error

type ProxyRequestInfo struct {
	EndpointGUID string
	URI          *url.URL
	UserGUID     string
	ResultGUID   string
	Headers      http.Header
	Body         []byte
	Method       string
}

type SessionStorer interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
	Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error
}

// ConnectedUser - details about the user connected to a specific service or UAA
type ConnectedUser struct {
	GUID  string `json:"guid"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

type JWTUserTokenInfo struct {
	UserGUID    string   `json:"user_id"`
	UserName    string   `json:"user_name"`
	TokenExpiry int64    `json:"exp"`
	Scope       []string `json:"scope"`
}

// Info - this represents user specific info
type Info struct {
	Versions     *Versions                             `json:"version"`
	User         *ConnectedUser                        `json:"user"`
	Endpoints    map[string]map[string]*EndpointDetail `json:"endpoints"`
	CloudFoundry *CFInfo                               `json:"cloud-foundry,omitempty"`
	PluginConfig map[string]string                     `json:"plugin-config,omitempty"`
}

// Extends CNSI Record and adds the user
type EndpointDetail struct {
	*CNSIRecord
	User          *ConnectedUser    `json:"user"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	TokenMetadata string            `json:"-"`
}

// Versions - response returned to caller from a getVersions action
type Versions struct {
	ProxyVersion    string `json:"proxy_version"`
	DatabaseVersion int64  `json:"database_version"`
}

type ConsoleConfig struct {
	UAAEndpoint         *url.URL `json:"uaa_endpoint"`
	ConsoleAdminScope   string   `json:"console_admin_scope"`
	ConsoleClient       string   `json:"console_client"`
	ConsoleClientSecret string   `json:"console_client_secret"`
	SkipSSLValidation   bool     `json:"skip_ssl_validation"`
	IsSetupComplete     bool     `json:"is_setup_complete"`
}

// CNSIRequest
type CNSIRequest struct {
	GUID     string `json:"-"`
	UserGUID string `json:"-"`

	Method      string      `json:"-"`
	Body        []byte      `json:"-"`
	Header      http.Header `json:"-"`
	URL         *url.URL    `json:"-"`
	StatusCode  int         `json:"statusCode"`
	Status      string      `json:"status"`
	PassThrough bool        `json:"-"`

	Response     []byte `json:"-"`
	Error        error  `json:"-"`
	ResponseGUID string `json:"-"`
}

type PortalConfig struct {
	HTTPClientTimeoutInSecs     int64    `configName:"HTTP_CLIENT_TIMEOUT_IN_SECS"`
	HTTPConnectionTimeoutInSecs int64    `configName:"HTTP_CONNECTION_TIMEOUT_IN_SECS"`
	TLSAddress                  string   `configName:"CONSOLE_PROXY_TLS_ADDRESS"`
	TLSCert                     string   `configName:"CONSOLE_PROXY_CERT"`
	TLSCertKey                  string   `configName:"CONSOLE_PROXY_CERT_KEY"`
	TLSCertPath                 string   `configName:"CONSOLE_PROXY_CERT_PATH"`
	TLSCertKeyPath              string   `configName:"CONSOLE_PROXY_CERT_KEY_PATH"`
	CFClient                    string   `configName:"CF_CLIENT"`
	CFClientSecret              string   `configName:"CF_CLIENT_SECRET"`
	AllowedOrigins              []string `configName:"ALLOWED_ORIGINS"`
	SessionStoreSecret          string   `configName:"SESSION_STORE_SECRET"`
	EncryptionKeyVolume         string   `configName:"ENCRYPTION_KEY_VOLUME"`
	EncryptionKeyFilename       string   `configName:"ENCRYPTION_KEY_FILENAME"`
	EncryptionKey               string   `configName:"ENCRYPTION_KEY"`
	AutoRegisterCFUrl           string   `configName:"AUTO_REG_CF_URL"`
	CFAdminIdentifier           string
	CloudFoundryInfo            *CFInfo
	HTTPS                       bool
	EncryptionKeyInBytes        []byte
	ConsoleVersion              string
	IsCloudFoundry              bool
	LoginHook                   LoginHookFunc
	SessionStore                SessionStorer
	ConsoleConfig               *ConsoleConfig
	PluginConfig                map[string]string
}
