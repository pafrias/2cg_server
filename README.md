# 2cgaming API

## Table of Contents
* [**Getting Started**](#getting-started)
  * [Installing Dependencies](#installing-dependencies)
  * [Database Configuration](#database-configuration)
  * [Starting Server](#starting-server)
* [**Trap Compendium API**](#trap-compendium-api)
* [**Epic Spell Codex API**](#epic-api)
* [**Change Log**](#log)
<hr>

## Getting started
### Installing Dependencies
  * go version go1.13 darwin/amd64
  * npm v6.12.1

### Database Configuration
  * Mysql 5.7
  * FILL ME IN


### Starting Server
  For security purposes, the username and password for the database are not stored in the repo but are accessed through environment variables at run time. For example, from the go root directory:
  ```bash
  go build app
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
## Epic Spell Codex API
  * FILL ME IN

## Change Log

|Name |Version |Date |Description |
|:- |:- |:- |:- |
|[@pafrias](https://github.com/code4sac) |0.1.0 |Oct-31-2019 |adksd