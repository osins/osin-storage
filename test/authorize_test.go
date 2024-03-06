package test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/osins/osin-simple/simple/config"
	"github.com/osins/osin-simple/simple/practices"
	"github.com/osins/osin-simple/simple/request"
	"github.com/osins/osin-storage/storage"
	"github.com/osins/osin-storage/storage/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
)

var (
	privatekeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----`)

	publickeyPEM = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----`)
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Debug("test init.")
}

func TestDb(t *testing.T) {
	dsn := storage.GetMySQLDSN()

	storage.Init(mysql.Open(dsn))

	storage.DB()
	storage.Migrate()

	NewClient(t)

	t.Logf(`Run Result:  TestDb() = %q`, dsn)
}

func NewClient(t *testing.T) {
	u := &model.User{
		Id:       uuid.UUID(uuid.New()),
		Username: "richard",
		EMail:    "296907@qq.com",
		Password: "123456",
	}

	userStore := storage.NewUserStorage()
	if err := userStore.Create(u); err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error(err)
	}

	client := &model.Client{
		Id:          uuid.MustParse("12465071-a4fa-4a45-b3f1-f0a972fb6875"),
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:14000/appauth",
		NeedLogin:   true,
		NeedRefresh: true,
	}

	// pg.NewClientManage().Delete(client.Id)
	if err := storage.NewClientStorage().Create(client); err != nil {
		logrus.WithFields(logrus.Fields{
			"client": client,
			"error":  err,
		}).Error("create client error")
	}

	c, err := storage.NewClientStorage().Get(client.Id.String())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":     client.Id,
			"client": c,
			"error":  err,
		}).Error("get client error")
	}

	logrus.WithFields(logrus.Fields{
		"client": c,
	}).Debug("get client done")

	NewServer(client, t)
}

func NewServer(client *model.Client, t *testing.T) {
	auhorize, err := practices.NewJwtAuthorize(
		"test server",
		config.Storage{
			Client:    storage.NewClientStorage(),
			User:      storage.NewUserStorage(),
			Authorize: storage.NewAuthorizeStorage(),
			Access:    storage.NewAccessStorage(),
		},
		practices.JwtKeyConfig{
			PrivateKey: privatekeyPEM,
			PublicKey:  publickeyPEM,
		})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error(err)
	}

	logrus.WithFields(logrus.Fields{
		"server": auhorize.Config().Name,
	}).Debug("server create done")

	request := &request.AuthorizeRequest{
		ClientId:     client.Id.String(),
		ClientSecret: client.Secret,
		ResponseType: request.AUTHORIZE_RESPONSE_REGISTER,
		RedirectUri:  "http://localhost:14000/appauth",
		State:        "",
		Username:     "wahahaha",
		Password:     "123456",
		EMail:        "",
		Mobile:       "",
	}

	response, err := auhorize.Register(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"client":  client,
		}).Error(err)

		t.Error(err)
	}

	logrus.WithField("response", response).Debug("register done")
}
