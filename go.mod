module github.com/osins/osin-storage

go 1.19

require (
	github.com/google/uuid v1.2.0
	github.com/osins/osin-simple v0.1.9
	github.com/sirupsen/logrus v1.9.3
	gorm.io/driver/mysql v1.5.4
	gorm.io/gorm v1.25.7-0.20240204074919-46816ad31dde
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/osins/osin-simple => /root/my/osin-simple
