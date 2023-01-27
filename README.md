# Gin-Framework


## Dev Configuration

### 1. Run Docker Service

```
docker-compose up
```

### 2. Wait for mysql & migrate the DB schema

**Create the database first**

```
./migrate.sh new DB
```

**Write down your db sql in migrations file**

```

-- +migrate Up
-- 新增pano DB
CREATE DATABASE IF NOT EXISTS `pano` DEFAULT CHARACTER SET utf8mb4;

-- +migrate Down
-- DROP DATABASE `pano`;

```

**Migrates the database to the most recent version available**

```
./migrate.sh up
```

#### Others 

**Undo a database migration**

```
./migrate.sh down
```

**Show migration status**

```
./migrate.sh status
```

**Create a new migration**

```
./migrate.sh new a_new_migration
```



## Prod Configuration

### 1. Copy config.yml to config-prod.yml & change the dsn, mode of mysql

```yml
server:
  version: v0.1
  addr: :80
  mode: prod
  static_dir: ./static
  # view_dir: ./view
  # upload_dir: ./storage
  max_multipart_memory: 50

python:
  dev_host: 127.0.0.1
  test_host: 127.0.0.1
  prod_host: pano-python

database-in-docker:
  dialect: mysql
  datasource: <user>:<password>@tcp(mysql:3306)/pano?charset=utf8mb4&timeout=10s&parseTime=True
  dir: migrations
  table: migrations
  max_idle_conns: 2
  max_open_conns: 16

database:
  dialect: mysql
  datasource: pano:ppaannoo@tcp(localhost:3306)/pano?charset=utf8mb4&timeout=10s&parseTime=True
  dir: migrations
  table: migrations
  max_idle_conns: 2
  max_open_conns: 16
```


### 2. Copy docker-compose.yml to Copy docker-compose-prod.yml & change the config of mysql

```yml
version: "3"
services:
  pano-python:
    build: ./dist/python
    container_name: pano-python
    volumes:
      - ./dist/static:/app/go/static
    ports:
      - 5000:5000
  pano-go:
    build: .
    container_name: pano-go
    volumes:
      - ./dist/python:/python
      - ./dist/log:/app/go/log
      - ./dist/static:/app/go/static
    depends_on:
      - mysql
    ports:
      - 80:80
  mysql:
    image: cap1573/mysql:5.6
    container_name: mysql
    restart: always
    environment:
      MYSQL_DATABASE: "pano"
      MYSQL_USER: <user>
      MYSQL_PASSWORD: <password>
      MYSQL_RANDOM_ROOT_PASSWORD: true
    ports:
      - "3306:3306"
    volumes:
      - ./dist/mysql:/var/lib/mysql

```

### 3. Start the service

```
docker-compose -f docker-compose-prod.yml up
```

### 4. Wait for mysql & migrate the DB schema

```
./migrate.sh up
```

## Generate swag document

All comments were written in router/router.go, so you need to find the path.

```
swag init -g ./router/router.go -o ./docs
```

## Generate wire service

1. Edit `router/wire.go` and write the initialization of dependency injections.

* Example
```go
//go:build wireinject
// +build wireinject

package router_v1

import (
	"github.com/google/wire"

	clinic_repository "go-pano/domain/repository/clinic"

	clinic_service "go-pano/domain/service/clinic"
	"go-pano/utils"
)

func initClinicService() clinic_service.IClinicService {
	wire.Build(
		clinic_service.NewClinicService,
		clinic_repository.NewClinicRepository,
		utils.NewDBInstance,
	)
	return nil
}
```

2. When you added the relationship of each class, run command `wire <relative_path>.`.

```bash
# I put the wire.go files in router folder. 
wire ./router/.
```

3. Add the functions of dependency injections where you want.

```go
package router_v1

import (
	"github.com/gin-gonic/gin"
	"go-pano/domain/delivery/http"
)

func NewRouter(app *gin.Engine) {

	api := app.Group("/api")
	{
		// After
		http.NewClinicHandler(api, initClinicService())

    // Before
    // cr := clinic_repository.NewClinicRepository(db)
		// cs := clinic_service.NewClinicService(cr)
		// http.NewClinicHandler(api, cs)

	}

}

```

## Use Protocol Buffers

Generate all pb files.

```
protoc ./protos/*/*.proto  --go_out=plugins=grpc:. --go_opt=paths=source_relative
```

## Problems

1. If there is a error when you migrated.
```
poabob@gengyingxiangdeMacBook-Pro go-pano % ./migrate.sh up    
/Users/poabob/go/bin/sql-migrate up -config=config.yml -env=database
Migration failed: Error 1146: Table 'pano.clinic' doesn't exist handling 20230119061136-Add-New-Tables.sql
```
> Just restart the container of mysql. It will be fixed.

## TODO

- Write down the real process
- Finish Predict Unit Tests
- User Login API
  - RBAC control
- Clinic Manage API
  - Clinic Token Request
- Record API
  - services used per month every clinic
  - services score list