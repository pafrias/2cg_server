# 2cgaming API

## Table of Contents
* [**Introduction**](#introduction)
* [**Getting Started**](#getting-started)
  * [Installing Dependencies](#installing-dependencies)
  * [Database Configuration](#database-configuration)
  * [Starting Server](#starting-server)
* [**Trap Compendium API**](#trap-compendium-api)
* [**Epic Spell Codex API**](#epic-api)
* [**Tutorial**](#tutorial)
* [**Change Log**](#log)
<hr>


## Introduction
This project had a few goals in mind, the foremost of which was to simplify and accelerate the work of future developers, specifically in extending the functionality of the application for future published works of 2CGaming, LLC. **Golang** was chosen for its fast compiling, iteration, and testing speeds, as well as its ease of learning and superior performance to as a backend language. By maintaining the seperation of concerns into discrete **go** packages, descriptive naming, and this documentation, future developers will hopefully find it easy to build upon this project.

`/ - main package` - instantiates a server, connects to the database, and routes requests to their respective handlers. Middleware is applied to route handlers at this level. Static asset requests are handled here, but all other request and response handling (to prefix ***/api***) should be handled in `/app/`.

`/app` - contains the request and response handling for the various APIs. Functions hang as methods off a `Handler` struct for each major service of the API, because, as names across services may become very similar, this will reduce naming conflicts.

`/middleware` - contains logic for form validation, authentication, logging, etc.

`/db` - connects to and operates with the database. All database controllers hang off a `Connection` struct to keep the namespace environment clean. Files in this package should **never** operate on response or request objects.

`/db/models` - contains various struct for easier json scanning and form validation. 

`/web` - contains frontend assets

--**FILL ME IN MORE**--

## Getting started
### Installing Dependencies
  * go version go1.13 darwin/amd64
  * npm v6.12.1

### Database Configuration
  * Mysql 5.7
  * FILL ME IN


### Starting Server
  For security purposes, the username and password for the database are not stored in the repo but are accessed through environment variables at run time.
  For example, from the go root directory:
  ```bash
  go get github.com/pafrias/2cgaming api

  go build
  
  SQL_USER="very_serious_name" SQL_PW='verys3curep@ssword' ./app
  ```
  Failure to connect will result in a fatal error. Check your variables

## Trap Compendium API
  Because of the ubiquity of the word component in many languages, packages, and contexts, the word 'effect' has been used in all Trap Compendium tables, apps, etc.

  Effects have 4 param fields, which are 8 tab seperated values, starting with the name of the field followed by the 7 tier values.

### `GET api/tc/components`
Requests all triggers, targets, and effects.

**Success Response**
  * **Code**: 200
  * **Content**: Array of component objects in the following format:

  |Key |Datatype |
  |:- |:- |
  |_id |int|
  |name |string|
  |type |string|
  |text |string|
  |cost |int |
  |param1|string|
  |param2|string|
  |param3|string|
  |param4|string|

**Error Response**
  * **Code:** 500 INTERNAL SERVER ERROR

### `GET api/tc/components/short`
Requests all triggers, targets, and effects, truncated to key values. This route is mostly for editing and posting new upgrades.

**Success Response**
  * **Code**: 200
  * **Content**: Array of component objects in the following format:

  |Key |Datatype |
  |:- |:- |
  |_id |int|
  |name |string|
  |type |string|

**Error Response**
  * **Code:** 500 INTERNAL SERVER ERROR
<hr>

## Epic Spell Codex API
  * FILL ME IN
<hr>

## Tutorial
If you are new to this project and **Go**, here would be the steps you would take to extend the application for a ***Monster Manual*** application, after deciding on data shape and tables.

1. Create `db/models/monstermanual.go`
    * export a struct type Monster:
    ```Go
    {
      ID          int
      CR          int
      Name        string
      Description string
    }
    ```

2. Create `db/monstermanual.go`
    * create a new method for the Connection struct, such as **GetMonsters()**. Ping the database, and query for monsters if successful
    ```Go
    func (c *Connection) GetComponents(ctx context.Context, queryType string) (r *sql.Rows, err error) {
      if err := c.Client.Ping(); err != nil {
        return nil, err
      }
      
      return c.Client.QueryContext(ctx, "SELECT * FROM mm_monsters")
    }
    ```

3. Create `/app/monstermanual` -> package `monstermanual`
    * In `app/monstermanual/app.go`, export a struct type MonsterManualAPI, and a function that returns an instance of MonsterManualAPI.

    ```Golang
    type MonsterManualAPI struct {
      DB *db.Connection
    }

    func NewHandler(db *db.Connection) Handler {
      return MonsterManualAPI{db}
    }
    ```

    * In `app/monstermanual/routes.go`, create a method for MonsterManualAPI that uses the new database controller, and returns a route Handler.

    ```Golang
    func (mm *MonsterManualAPI) GetMonsters() http.HandlerFunc {

      return func(res http.ResponseWriter, req *http.Request) {
        ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

        rows, err := mm.DB.GetMonsters(ctx)
        if err != nil {
          //handle error
        }
        defer rows.Close()

        var monsters []models.Monster

        for rows.Next() {
          var m models.Monster
          err = rows.Scan(&m.ID, &m.CR, &m.Name, &m.Description)
          if err != nil {
          //handle error
          }
          monsters = append(monsters, m)
        }

        data, err := json.Marshal(monsters)
        if err != nil {
          //handle error
        }

        res.Header().Set("Content-Type", "application/json")
        res.Write(data)
      }
    }
    ```

  4. in `/router.go` connect new route
      * create a new method for the **server** struct, such as **applyMonsterManualRoutes()**
      ```Golang
      func (s *server) applyMonsterManualRoutes(router *mux.Router) {
        MMAPI := monstermanual.NewHandler{s.db}
        router.HandlerFunc("/monsters", MMAPI.GetMonsters()).Methods("GET")
      ```
      * Connect the new router to the main router
      ```Golang
        func (s *server) createMainRouter() {
          r := mux.NewRouter()
          // apply other routes

          subrouter := r.PathPrefix("/api/monstermanual").Subrouter()
          s.applyMonsterManualRoutes(subrouter)

          //...
        }
      ```
      


## Change Log

|Name |Version |Date |Description |
|:- |:- |:- |:- |
|[@pafrias](https://github.com/pafrias) |0.1.0 | Nov-5-2019 | Create tutorial, scrap custom validator for *go-playground/form* validator.
|[@pafrias](https://github.com/pafrias) |0.1.0 |Oct-31-2019 | Create version 0.1.0