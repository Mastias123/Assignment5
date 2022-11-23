# Assignment5

To run the program stand in the root of the program and write the following command in the terminal:

'go run server/server.go'

This will start tree servers at port 5001, 5002 and 5003.

Then open any number of new terminals and and write the following command:

'go run client/client.go cPort 8080 ... -cId x'

Here you have to substitude the x after cId for a unieque number.
The command will start a client, and connect it to the three servers.