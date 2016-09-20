package vcstokens

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/hpcloud/portal-proxy/repository/tokens"
)

const (
	findVCSToken = `SELECT access_token
                  FROM vcstokens
                  WHERE endpoint = $1 AND user_guid = $2`

	countVCSTokens = `SELECT COUNT(*)
	                  FROM vcstokens
	                  WHERE endpoint = $1 AND user_guid = $2`

	insertVCSToken = `INSERT INTO vcstokens (endpoint, user_guid, access_token)
		                VALUES ($1, $2, $3)`

	updateVCSToken = `UPDATE vcstokens
		                SET access_token = $3
		                WHERE endpoint = $1 AND user_guid = $2`
)

// PgsqlTokenRepository is a PostgreSQL-backed token repository
type PgsqlVCSTokenRepository struct {
	db *sql.DB
}

// NewPgsqlTokenRepository - get a reference to the token data source
func NewPgsqlVCSTokenRepository(dcp *sql.DB) (Repository, error) {
	log.Println("NewPgsqlTokenRepository")
	return &PgsqlVCSTokenRepository{db: dcp}, nil
}

// SaveVCSToken - Save the VCS token to the datastore
func (p *PgsqlVCSTokenRepository) SaveVCSToken(tr VCSTokenRecord, encryptionKey []byte) error {
	log.Println("SaveVCSToken")

	if tr.UserGUID == "" {
		msg := "Unable to save VCS Token without a valid VCS user id."
		log.Println(msg)
		return fmt.Errorf(msg)
	}

	if tr.Endpoint == "" {
		msg := "Unable to save VCS Token without a valid VCS endpoint."
		log.Println(msg)
		return fmt.Errorf(msg)
	}

	if tr.AccessToken == "" {
		msg := "Unable to save VCS Token without a valid VCS access token."
		log.Println(msg)
		return fmt.Errorf(msg)
	}

	log.Println("Encrypting Access Token")
	ciphertextAccessToken, err := EncryptToken(encryptionKey, tr.AccessToken)
	if err != nil {
		return err
	}

	// Is there an existing token?
	var count int
	err = p.db.QueryRow(countVCSTokens, tr.Endpoint, tr.UserGUID).Scan(&count)
	if err != nil {
		log.Printf("Unknown error attempting to find UAA token: %v", err)
	}

	switch count {
	case 0:

		log.Println("Performing INSERT of encrypted tokens")
		if _, err := p.db.Exec(insertVCSToken, tr.Endpoint, tr.UserGUID, ciphertextAccessToken); err != nil {
			msg := "Unable to INSERT VCS token: %v"
			log.Printf(msg, err)
			return fmt.Errorf(msg, err)
		}

		log.Println("UAA token INSERT complete")

	default:

		log.Println("Performing UPDATE of encrypted tokens")
		if _, updateErr := p.db.Exec(updateVCSToken, tr.Endpoint, tr.UserGUID, ciphertextAccessToken); updateErr != nil {
			msg := "Unable to UPDATE VCS token: %v"
			log.Printf(msg, updateErr)
			return fmt.Errorf(msg, updateErr)
		}

		log.Println("VCS token UPDATE complete.")
	}

	return nil
}

// FindVCSToken - return the VCS token from the datastore
func (p *PgsqlVCSTokenRepository) FindVCSToken(endpoint string, userGUID string, encryptionKey []byte) (VCSTokenRecord, error) {
	log.Println("FindUAAToken")
	if endpoint == "" {
		msg := "Unable to find VCS Token without a valid endpoint."
		log.Printf(msg)
		return VCSTokenRecord{}, fmt.Errorf(msg)
	}

	if userGUID == "" {
		msg := "Unable to find VCS Token without a valid UAA user GUID."
		log.Printf(msg)
		return VCSTokenRecord{}, fmt.Errorf(msg)
	}

	// temp vars to retrieve db data
	var (
		plaintextAccessToken  string
		ciphertextAccessToken []byte
	)

	// Get the VCS record from the db
	err := p.db.QueryRow(findVCSToken, endpoint, userGUID).Scan(&ciphertextAccessToken)
	if err != nil {
		msg := "Unable to Find VCS token: %v"
		log.Printf(msg, err)
		return VCSTokenRecord{}, fmt.Errorf(msg, err)
	}

	log.Println("Decrypting Access Token")
	plaintextAccessToken, err = decryptToken(encryptionKey, ciphertextAccessToken)
	if err != nil {
		return VCSTokenRecord{}, err
	}

	// Build a new TokenRecord based on the decrypted tokens
	tr := new(VCSTokenRecord)
	tr.AccessToken = plaintextAccessToken
	tr.Endpoint = endpoint
	tr.UserGUID = userGUID

	return *tr, nil
}

// Note:
// When it's time to store the encrypted token in PostgreSQL, it's gets a bit
// hairy. The encrypted token is binary data, not really text data, which
// typically has a character set, unlike binary data. Generally speaking, it
// comes down to one of two choices: store it in a bytea column, and deal with
// some funkiness; or store it in a text column and make sure to base64 encode
// it going in and decode it coming out.
// https://wiki.postgresql.org/wiki/BinaryFilesInDB
// http://engineering.pivotal.io/post/ByteA_versus_TEXT_in_PostgreSQL/
// I chose option 1.

// encryptToken - TBD
func EncryptToken(key []byte, t string) ([]byte, error) {
	log.Println("encryptToken")
	var plaintextToken = []byte(t)
	ciphertextToken, err := tokens.Encrypt(key, plaintextToken)
	if err != nil {
		msg := "Unable to encrypt token: %v"
		log.Printf(msg, err)
		return nil, fmt.Errorf(msg, err)
	}

	return ciphertextToken, nil
}

// Note:
// When it's time to store the encrypted token in PostgreSQL, it's gets a bit
// hairy. The encrypted token is binary data, not really text data, which
// typically has a character set, unlike binary data. Generally speaking, it
// comes down to one of two choices: store it in a bytea column, and deal with
// some funkiness; or store it in a text column and make sure to base64 encode
// it going in and decode it coming out.
// https://wiki.postgresql.org/wiki/BinaryFilesInDB
// http://engineering.pivotal.io/post/ByteA_versus_TEXT_in_PostgreSQL/
// I chose option 1.

// decryptToken - TBD
func decryptToken(key, t []byte) (string, error) {
	log.Println("decryptToken")
	plaintextToken, err := tokens.Decrypt(key, t)
	if err != nil {
		msg := "Unable to decrypt token: %v"
		log.Printf(msg, err)
		return "", fmt.Errorf(msg, err)
	}

	return string(plaintextToken), nil
}
