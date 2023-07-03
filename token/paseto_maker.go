package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

// var PasetoPayLoad = make(map[string]interface{})
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker(key []byte) (Maker, error) {
	/*tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	PasetoPayLoad["ID"] = tokenID
	PasetoPayLoad["Username"] = username
	PasetoPayLoad["IssuedAt"] = time.Now()
	PasetoPayLoad["ExpiredAt"] = time.Now().Add(duration)

	return PasetoPayLoad, nil */
	symmetricKey, err := paseto.V4SymmetricKeyFromBytes(key)
	if err != nil {
		return nil, err
	}
	/*maker := PasetoMaker{
		symmetricKey: symmetricKey,
	}*/
	return &PasetoMaker{symmetricKey}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	payload := map[string]interface{}{
		"ID":       tokenID,
		"username": username,
		"iat":      time.Now(),
		"exp":      time.Now().Add(duration),
	}
	/*fmt.Println("Create payload.ID: ", payload["ID"])
	fmt.Println("Create payload.Username: ", payload["username"])
	fmt.Println("Create payload.IssuedAT: ", payload["iat"])
	fmt.Println("Create payload.ExpiredAt: ", payload["exp"])
	fmt.Println()*/
	token, err := paseto.MakeToken(payload, nil)
	//key := paseto.V4SymmetricKeyFromHex()

	encrypted := token.V4Encrypt(maker.symmetricKey, nil)
	return encrypted, nil
}

func (maker *PasetoMaker) VerifyToken(encrypted string) (*Payload, error) {
	parser := paseto.NewParser()
	token, err := parser.ParseV4Local(maker.symmetricKey, encrypted, nil) // this will fail if parsing failes, cryptographic checks fail, or validation rules fail
	if err != nil {
		return nil, err
	}
	payload := Payload{}
	token.Get("ID", &payload.ID)
	token.Get("username", &payload.Username)
	token.Get("iat", &payload.IssuedAt)
	token.Get("exp", &payload.ExpiredAt)
	//claims := token.Claims()
	/*fmt.Println("Verify payload: ", payload)
	fmt.Println("Verify payload.ID: ", payload.ID)
	fmt.Println("Verify payload.Username: ", payload.Username)
	fmt.Println("Verify payload.IssuedAt: ", payload.IssuedAt)
	fmt.Println("Verify payload.CreatedAt: ", payload.ExpiredAt)*/

	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
