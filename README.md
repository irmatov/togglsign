## How to run

    (cd cmd; go build -o ../signer)
    export SIGNER_DSN="postgres://postgres:mypass@localhost:5432/toggldb"
    export SIGNER_JWT_KEY="mysecrettoken"
    export SIGNER_SIGN_KEY="mysignkey"
    export SIGNER_GRPC_ADDRESS=":12345"
    ./signer

## How to run simplest test client

   cd testclient
   go build
   export SIGNER_JWT_KEY="mysecrettoken"
   ./testclient localhost:12345
