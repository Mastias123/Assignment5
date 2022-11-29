# Assignment5

To run the program stand in the root of the program and write the following command three different terminals:

'go run server/server.go -sPort 5001'

'go run server/server.go -sPort 5002'

'go run server/server.go -sPort 5002'

This will start tree servers at port 5001, 5002 and 5003.

Then open two new terminals and and write the following command:

'go run client/client.go -cPort 8080 -sPort1 5001 -sPort2 5002 -sPort3 5003 -cId x

Here you have to substitude the x after cId for a unique number.
The command will start a client, and connect it to the three servers.

To bid input in any given client terminal: "bid"

To check result input: "result" in the client terminal