# Docs

- migration
- config.env
- swagger

```js
{
  "surname": "Иванов",
  "name": "Иван",
  "patronymic": "Иванович",
  "address": "г. Москва, ул. Ленина, д. 5, кв. 1",
  "passportNumber": "1234 567890" // ser space num
}
```
User

```js
{
    "startTime": "some start time"
    "endTime": "some end time"
}
```
Task

## GET /users

- Filtering
- Pagination

## PUT /users

```js
// DATA
```
Body

## PATCH /users/:userId

```js
// DATA
```
Body

## DELETE /users/:userId

## GET /tasks/time-summary/:userId

- time range sort
- Desc sort

## PATCH /tasks/start/:userId

## PATCH /tasks/end/:userId

---
## GET /info?passportSeries=value&passportNumber=value

req 200, 400, 500

After add ???
