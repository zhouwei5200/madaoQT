package exchange

import (
	"errors"
	Mongo "madaoQT/mongo"
	Utils "madaoQT/utils"
)

type ExchangeInfo struct {
	/* User password encrypted */
	API string
	/* User password encrypted */
	Secret string
}

func GetExchangeKey(mongo *Mongo.ExchangeDB, exchange string, key []byte, nonce []byte) (error, *ExchangeInfo) {

	if mongo == nil {
		return errors.New("Invalid mongo handler"), nil
	}

	if err := mongo.Connect(); err != nil {
		return err, nil
	}

	defer mongo.Close()

	err, record := mongo.FindOne(exchange)
	if err != nil {
		return errors.New("APIKEY not found"), nil
	}

	crypto := Utils.AESCrypto{
		Type:  Utils.AESTypeBuffer,
		Key:   key,
		Nonce: nonce,
	}

	err, plainAPI := crypto.DecryptInMemory(record.API)
	if err != nil {
		return err, nil
	}

	err, plainSecret := crypto.DecryptInMemory(record.Secret)
	if err != nil {
		return err, nil
	}

	return nil, &ExchangeInfo{
		API:    string(plainAPI),
		Secret: string(plainSecret),
	}
}

func AddExchangeKey(mongo *Mongo.ExchangeDB, name string, api string, secret string) error {

	if mongo == nil {
		return errors.New("Invalid mongo handler")
	}

	if err := mongo.Connect(); err != nil {
		return err
	}

	defer mongo.Close()

	crypto := Utils.AESCrypto{
		Type: Utils.AESTypeBuffer,
	}

	err, encryptedAPI := crypto.EncryptInMemory([]byte(api))
	if err != nil {
		return err
	}

	err, encryptedSecret := crypto.EncryptInMemory([]byte(secret))
	if err != nil {
		return err
	}

	return mongo.Insert(&Mongo.ExchangeInfo{
		Name:   name,
		API:    encryptedAPI,
		Secret: encryptedSecret,
	})

}
