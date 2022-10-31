# Restaurant Management API

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)

This repository contains the restaurant management service API

## Clone the project

```
$ git clone https://github.com/minhtran241/restaurant-management.git
$ cd restaurant-management
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

## LICENSE

MIT

## Contributor

Minh Tran (Me)

<!-- |            Endpoints             |            Descriptions            | Methods |
| :------------------------------: | :--------------------------------: | :-----: |
|                /                 |     Show the status of server      |   GET   |
|       /swagger/index.html        |             Swagger UI             |   GET   |
|              /users              |           Get all users            |   GET   |
|         /users/:user_id          |          Get single user           |   GET   |
|         /users/:user_id          |            Update user             |  PATCH  |
|          /users/signup           |              Sign up               |  POST   |
|           /users/login           |               Login                |  POST   |
|              /foods              |           Get all foods            |   GET   |
|              /foods              |            Create food             |  POST   |
|         /foods/:food_id          |          Get single food           |   GET   |
|         /foods/:food_id          |            Update food             |  PATCH  |
|              /menus              |           Get all menus            |   GET   |
|              /menus              |            Create menu             |  POST   |
|         /menus/:menu_id          |          Get single menu           |   GET   |
|         /menus/:menu_id          |            Update menu             |  PATCH  |
|            /invoices             |          Get all invoices          |   GET   |
|            /invoices             |           Create invoice           |  POST   |
|      /invoices/:invoice_id       |         Get single invoice         |   GET   |
|      /invoices/:invoice_id       |           Update invoice           |  PATCH  |
|             /orders              |           Get all orders           |   GET   |
|             /orders              |            Create order            |  POST   |
|        /orders/:order_id         |          Get single order          |   GET   |
|        /orders/:order_id         |            Update order            |  PATCH  |
|           /order-items           |       Get all ordered items        |   GET   |
|           /order-items           |        Create ordered item         |  POST   |
|    /orderItems/:order_item_id    |    Get all single ordered items    |   GET   |
| /orderItems-order/:order_item_id | Get all ordered items in one order |   GET   |
|    /orderItems/:order_item_id    |        Update ordered item         |  PATCH  |
|             /tables              |           Get all tables           |   GET   |
|        /tables/:table_id         |          Get single table          |   GET   |
|             /tables              |            Create table            |  POST   |
|        /tables/:table_id         |            Update table            |  PATCH  | -->

</div>
