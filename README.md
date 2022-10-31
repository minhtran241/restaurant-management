# Restaurant Management API

[![Repo Size](https://img.shields.io/github/repo-size/minhtran241/restaurant-management)](https://github.com/minhtran241/restaurant-management)
[![License](https://img.shields.io/github/license/minhtran241/restaurant-management)](https://github.com/minhtran241/restaurant-management)
[![Commit](https://img.shields.io/github/last-commit/minhtran241/restaurant-management)](https://github.com/minhtran241/restaurant-management)

This repository contains the restaurant management service API

## Package index

List dependencies used in this system by using [github.com/ribice/glice](https://github.com/ribice/glice), a Golang license and dependency checker. The package prints list of all dependencies, their URL, license and saves all the license files in /licenses

```
+----------------------------------------+--------------------------------------------+------------+
|               DEPENDENCY               |                  REPOURL                   |  LICENSE   |
+----------------------------------------+--------------------------------------------+------------+
| github.com/gin-gonic/gin               | https://github.com/gin-gonic/gin           | MIT        |
| github.com/go-playground/validator/v10 | https://github.com/go-playground/validator | MIT        |
| github.com/golang-jwt/jwt/v4           | https://github.com/golang-jwt/jwt          | MIT        |
| github.com/swaggo/swag                 | https://github.com/swaggo/swag             | MIT        |
| go.mongodb.org/mongo-driver            | https://github.com/mongodb/mongo-go-driver | Apache-2.0 |
| golang.org/x/crypto                    | https://go.googlesource.com/crypto         |            |
| github.com/swaggo/files                | https://github.com/swaggo/files            | MIT        |
| github.com/swaggo/gin-swagger          | https://github.com/swaggo/gin-swagger      | MIT        |
+----------------------------------------+--------------------------------------------+------------+
```

## Endpoint Table

<div align="center">

|            Endpoints             |            Descriptions            | Methods |
| :------------------------------: | :--------------------------------: | :-----: |
|                /                 |     Show the status of server      |   GET   |
|       /swagger/index.html        |             Swagger UI             |   GET   |
| /orderItems-order/:order_item_id | Get all ordered items in one order |   GET   |

|    Method    |      User       |      Food       |      Menu       |        Invoice        |       Order       |       Ordered Item        |       Table       |
| :----------: | :-------------: | :-------------: | :-------------: | :-------------------: | :---------------: | :-----------------------: | :---------------: |
|  GET (all)   |     /users      |     /foods      |     /menus      |       /invoices       |      /orders      |        /orderItems        |      /tables      |
| GET (single) | /users/:user_id | /foods/:food_id | /menus/menu_id  | /invoices/:invoice_id | /orders/:order_id | /orderItems/order_item_id | /tables/:table_id |
|     POST     |  /users/signup  |     /foods      |     /menus      |       /invoices       |      /orders      |        /orderItems        |      /tables      |
|     POST     |  /users/login   |                 |                 |                       |                   |                           |                   |
|    PATCH     | /users/:user_id | /foods/:food_id | /menus/:menu_id | /invoices/:invoice_id | /orders/:order_id | /orderItems/order_item_id | /tables/:table_id |

</div>

## License

The project is licensed under the MIT license. Check the [LICENSE](LICENSE) file for details

## Contributor

Minh Tran (Me)
