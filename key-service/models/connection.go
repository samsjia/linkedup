package models

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	infoTableName      = "linkedup-keyservice"
	authTableName      = "linkedup-keyservice-auth"
	emailTableName     = "linkedup-email"
	blacklistTableName = "linkedup-blacklist"
)

var (
	forceS3PathStyle = true
)

// DatabaseContext carries context needed to interact with the database.
type DatabaseContext struct {
	db            *dynamodb.DynamoDB
	s3            *s3.S3
	contentBucket string
}

// NewDatabaseContext will establish a session with the backend db.
//
// `DatabaseContext` effectively acts as a key-value store for a variety of operations
func NewDatabaseContext(localstack bool, contentBucket string) (DatabaseContext, error) {
	return NewDatabaseContextWithCfg(
		session.Must(session.NewSession()),
		localstack,
		contentBucket,
	)
}

// NewDatabaseContextWithCfg constructs a new DatabaseContext, using the given AWS
// session handle.
func NewDatabaseContextWithCfg(cfg client.ConfigProvider, localstack bool,
	bucket string) (context DatabaseContext, err error) {
	if localstack {
		context.db = dynamodb.New(
			cfg,
			&aws.Config{
				Endpoint: aws.String("http://localstack:4569"),
			},
		)
		context.s3 = s3.New(
			cfg,
			&aws.Config{
				Endpoint:         aws.String("http://localstack:4572"),
				S3ForcePathStyle: &forceS3PathStyle,
			},
		)
	} else {
		context.db = dynamodb.New(cfg)
		context.s3 = s3.New(cfg)
	}
	context.contentBucket = bucket

	log.Info("establishing session with dynamo")

	// try create the tables if they haven't already been instantiated
	if err := createTables(context.db); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != dynamodb.ErrCodeResourceInUseException {
				log.Info("dynamo table already created")
				return context, nil
			}
		} else {
			return context, err
		}
	}

	return context, nil
}

func createTables(db *dynamodb.DynamoDB) error {

	/** create table to store attendee information **/
	err := createTable(db, infoTableName)
	if err != nil {
		return err
	}

	/** create table to store auth tokens for key recovery **/
	err = createTable(db, authTableName)
	if err != nil {
		return err
	}

	/** create table to store emails attendees that change it manually **/
	err = createTable(db, emailTableName)
	if err != nil {
		return err
	}

	_, err = db.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(blacklistTableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Email"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Email"),
				KeyType:       aws.String("HASH"),
			},
		},
	})
	
	return err
}

func createTable(db *dynamodb.DynamoDB, tableName string) error {
	_, err := db.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		TableName:   aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
	})
	return err
}

/** Storage **/

// StoreAttendeeInfo -
func (db DatabaseContext) StoreAttendeeInfo(id int, info []byte) bool {
	value := &storedInfo{
		ID:   id,
		Data: info,
	}

	return setInfo(&db, value)
}

// StoreVerificationToken -
func (db DatabaseContext) StoreVerificationToken(id int, token string) bool {
	auth := &storedAuth{
		ID:        id,
		AuthToken: token,
	}

	return setVerificationToken(&db, auth)
}

// StoreEmail sets the email address for that id
func (db DatabaseContext) StoreEmail(id int, address string) bool {
	email := &storeEmail{
		ID:    id,
		Email: address,
	}

	return setEmail(&db, email)
}

/** Retrieval **/

// GetAttendeeInfo -
func (db DatabaseContext) GetAttendeeInfo(id int) ([]byte, error) {
	return getInfoForID(&db, id)
}

// GetVerificationToken -
func (db DatabaseContext) GetVerificationToken(id int) (string, error) {
	return getVerificationTokenForID(&db, id)
}

//GetEmail gets the associated email for that id, expect empty string since few attendees set a new email manually
func (db DatabaseContext) GetEmail(id int) string {
	e := getEmailForID(&db, id)
	if e == nil {
		return ""
	}
	return e.Email
}

// GetImageUploadURL get a URL that an image can be uploaded to
func (db DatabaseContext) GetImageUploadURL(id int) (string, error) {
	req, _ := db.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(db.contentBucket),
		Key:    aws.String(fmt.Sprintf("avatars/%d", id)),
	})

	result, err := req.Presign(60 * time.Minute)
	if err != nil {
		return "", err
	}

	return result, nil
}

// GetBlacklistEntry checks if a particular email is blacklisted
func (db DatabaseContext) GetBlacklistEntry(email string) bool {
	entry := getBlacklistEntry(&db, email)
	if entry != nil {
		return entry.Blacklisted
	}
	return false
}
