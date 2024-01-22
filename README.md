## How to run

    (cd cmd; go build -o ../signer)
    export SIGNER_DSN="postgres://postgres:mypass@localhost:5432/toggldb"
    export SIGNER_JWT_KEY="mysecrettoken"
    export SIGNER_SIGN_KEY="mysignkey"
    export SIGNER_GRPC_ADDRESS=":12345"
    ./signer

One can use `database.sql` to create a required database structure.

The applications (and unit tests) uses PostgreSQL, so it must be running somewhere. To run it locally with docker:

    docker container run --detach --name toggldb -e POSTGRES_PASSWORD=mypass -p 5432:5432 postgres:latest

## How to run simplest test client

   cd testclient
   go build
   export SIGNER_JWT_KEY="mysecrettoken"
   ./testclient localhost:12345

## Other notes

I didn't really worked with JWT tokens before, so I spend a lot of time figuring it out and testing it.
Finally found there was an error in variable name that I used.

Also, spec is (intentionally?) vague on how to process JWT tokens (issuer?) or what is means to "sign" it.
Not to drown myself with crypto stuff I have used the simples "signing" method that could possibly work.

Also spent some time setting up other things, like GRPC proto/ server/ client. Although I'm a big
supporter of unit/ other kinds of testing, this time I couldn't even start writing all the needed tests.
(Maybe going with simple REST API could have saved some time for me? I don't know.)

On the app logic side, it is also not spelled out what the "test" is, how "questions" are represented.
