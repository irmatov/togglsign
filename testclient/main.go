// Test client. Uses the same SIGNER_JWT_KEY variable as a server.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/irmatov/togglsign/infra/ports/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClaims struct {
	Email string
	jwt.RegisteredClaims
}

func createToken(email string) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  []string{},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "01H4EKGQSY5637MQP395283JR8",
		},
	})
	s, err := t.SignedString([]byte(os.Getenv("SIGNER_JWT_KEY")))
	if err != nil {
		panic(err)
	}
	return s
}

func printJSON(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: mytest <server-grpc-address>\nExample: mytest localhost:12345")
		os.Exit(1)
	}
	conn, err := grpc.Dial(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := rpc.NewSignerClient(conn)
	ctx := context.Background()
	signResp, err := client.Sign(ctx, &rpc.SignRequest{
		JwtToken: createToken("test@example.org"),
		Responses: []*rpc.Response{
			{Question: "2+2", Answer: "4"},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Sign response:")
	printJSON(signResp)

	verifyResp, err := client.Verify(ctx, &rpc.VerifyRequest{
		Email:     "test@example.org",
		Signature: signResp.Signature,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Verify response:")
	printJSON(verifyResp)
}
