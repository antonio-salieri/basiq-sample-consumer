# basiq-sample-consumer
Sample Basiq API consumer which is currently able to:
* read/create/delete user using Basiq API
* fetch all user transactions and calculate average amount of debit transactons per transaction category

# Commands
Help about about specific command can be get using one of the following commands
```
basiq-sample-consumer --help
basiq-sample-consumer user [subcommand] --help
basiq-sample-consumer transaction [subcommand] --help
```
## User manipulation
### Create user
```
BASIQ_API_KEY=<YOUR API KEY> basiq-sample-consumer user create --email <user email address> --mobile <mobile phone number>
```
### Fetch user
```
BASIQ_API_KEY=<YOUR API KEY> basiq-sample-consumer user get --userID <user ID received upon user creation>
```
### Delete user
```
BASIQ_API_KEY=<YOUR API KEY> basiq-sample-consumer user delete --userID <user ID received upon user creation>
```
## Transactions

### Calculate average transaction amount per transaction category
```
BASIQ_API_KEY=<YOUR API KEY> basiq-sample-consumer transaction --userID <user ID received upon user creation> [--transactionType <id of transaction type to process>] [--institutionID <institution ID (default "AU00000")>] [--loginID <login ID for institution>] [--loginPassword <login password for institution>]
```

Sample output:
```
[~/go/src/github.com/antonio-salieri/basiq-sample-consumer]> BASIQ_API_KEY="<SOME_KEY>" ./basiq-sample-consumer transaction average-transaction -u bb04fc1c-f3c9-4590-8d33-3044f3de9735

2019/02/03 18:17:31 Requesting: https://au-api.basiq.io/token
2019/02/03 18:17:33 Requesting: https://au-api.basiq.io/users/bb04fc1c-f3c9-4590-8d33-3044f3de9735/connections?filter=institution.id.eq('AU00000')
2019/02/03 18:17:34 Requesting: https://au-api.basiq.io/users/bb04fc1c-f3c9-4590-8d33-3044f3de9735/transactions?filter=connection.id.eq('a15de6a5-cb55-4814-a007-0525ea2a850a')
2019/02/03 18:17:37 Requesting: https://au-api.basiq.io/users/bb04fc1c-f3c9-4590-8d33-3044f3de9735/transactions?next=98e8869a-d4fa-4c32-ad24-a96abdb8bc70&filter=connection.id.eq('a15de6
a5-cb55-4814-a007-0525ea2a850a')
2019/02/03 18:17:40 Requesting: https://au-api.basiq.io/users/bb04fc1c-f3c9-4590-8d33-3044f3de9735/transactions?next=f28476b8-729f-483d-8538-e1111f97443e&filter=connection.id.eq('a15de6
a5-cb55-4814-a007-0525ea2a850a')
2019/02/03 18:17:44 Fetched 926 transactions

Code            | Average               | Total                 | Count                  |Title 
-------------------------------------------------------------------------------------------------------------------------------
412             | -201.13               | -16291.92             | 81                     |Specialised Food Retailing
400             | -385.13               | -67783.19             | 176                    |Fuel Retailing
0               | -92.22                | -4703.03              | 51                     |Unknown
452             | -92.58                | -5832.74              | 63                     |Pubs, Taverns and Bars
451             | -34.03                | -6568.74              | 193                    |Cafes, Restaurants and Takeaway Food Services
411             | -135.09               | -25262.33             | 187                    |Supermarket and Grocery Stores
```

# TODO
- [ ] Add unit tests (DDT :()
- [ ] Refactor app initialization (Client and cmd processor creation in main.go) so Basiq session is not requested before command is validated
- [x] Fix issue with creating connection using current version of github.com/basiqio/basiq-sdk-golang
- [ ] Add `user create` command validation that checks if either `email` or `mobile` is passed
- [ ] Improve formula for calculating average transaction, so not all transactions are stored in memory during calculation
- [ ] Expose commands for manipulating user connections
- [ ] Persist already fetched data and make new request only when forced in command
