# Gin-Framework

## Dev Configuration

### Run Docker Service

```
docker-compose up
```

## SQL init

If you are the first time starting this project, just use below sql commands to login.

```sql
INSERT INTO `Users` (`user_id`, `name`, `account`, `password`, `roles_string`, `status`) VALUES
(1, 'ADMIN', '__pano_admin__', '$2a$10$94v4wlp6ZRanI6Xv1k4hyePZJlTJf.o08fSUqPby/mABlGGgRiRAa', '[\"admin\"]', 1);
INSERT INTO `Clinic` (`clinic_id`, `name`, `start_at`, `end_at`, `quota_per_month`, `token`) VALUES
(1, '測試診所', '2022-10-10', '2099-12-31', 200, 'rHsxKe6qPxxoZJh2oPJPk2mVzNFB5XmfOnkLpCwvbhnOnbU9i3');
```


## Prod Configuration

### 1. Copy config.yml to config-prod.yml & change the dsn, mode of mysql

```yml
server:
  version: v0.1
  addr: :80
  mode: prod
  static_dir: ./static
  public_dir: ./public
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
  datasource: <user>:<password>@tcp(mysql:3306)/pano?charset=utf8mb4&timeout=10s&parseTime=True
  dir: migrations
  table: migrations
  max_idle_conns: 2
  max_open_conns: 16
```


### 2. Copy docker-compose.yml to Copy docker-compose-prod.yml, change the config of mysql, and chang the target of pano-go service.

```yml
version: "3"
services:
  pano-python:
    image: poabob/pano-python:prod-1.0.0
    build: ./dist/python
    container_name: pano-python
    volumes:
      - ./dist/static:/app/go/static
    ports:
      - 5001:5001
  pano-go:
    image: poabob/pano-go:prod-1.0.0
    build: 
      context: .
      dockerfile: Dockerfile
      target: prod
    container_name: pano-go
    volumes:
      - ./dist/log:/app/go/log
      - ./dist/static:/app/go/static
      - ./dist/public:/app/go/public
    depends_on:
      - mysql
    ports:
      - 80:80
  mysql:
    image: mariadb:10.9
    container_name: mysql
    restart: always
    environment:
      MYSQL_DATABASE: "<DB>"
      MYSQL_USER: "<user>"
      MYSQL_PASSWORD: "<password>"
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


## Generate swag document

All comments were written in router/router.go, so you need to find the path.

```
# In Docker
docker exec pano-go swag init -g ./router/router_v1.go -o ./docs

# In real mechine
swag init -g ./router/router_v1.go -o ./docs
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
# In Docker
docker exec pano-go wire ./router/.

# In real mechine
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
# In Docker
docker exec pano-go protoc ./protos/*/*.proto  --go_out=plugins=grpc:. --go_opt=paths=source_relative

# In real mechine
protoc ./protos/*/*.proto  --go_out=plugins=grpc:. --go_opt=paths=source_relative
```

## Furture Table

```sql
CREATE TABLE IF NOT EXISTS `Record` (
  `record_id` int PRIMARY KEY AUTO_INCREMENT COMMENT '紀錄ID',
  `clinic_id` int DEFAULT 0 COMMENT '診所ID',
  `predict_id` int DEFAULT 0 COMMENT '預測ID',
  `score` int NOT NULL DEFAULT 80 COMMENT '準確度',
  `comment` varchar(1024) NOT NULL DEFAULT "" COMMENT '準確度評論'
)ENGINE = InnoDB DEFAULT CHARSET=utf8mb4;
```

## TODO

- Write down the real process
- Finish Predict Unit Tests
- User RBAC control
- Clinic Token Request
- Record API
  - services used per month every clinic
  - services score list