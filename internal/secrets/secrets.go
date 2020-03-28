package secrets

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func Get(ctx context.Context, secretName string) (map[string]string, error) {
	secretsMgr := secretsmanager.New(session.New())
	secret, err := secretsMgr.GetSecretValueWithContext(ctx, &secretsmanager.GetSecretValueInput{SecretId: &secretName})
	if err != nil {
		return nil, err
	}

	var secretValues map[string]string
	if err := json.Unmarshal([]byte(*secret.SecretString), &secretValues); err != nil {
		return nil, err
	}

	return secretValues, nil
}

func Create(ctx context.Context, secretName string) error {
	secretsMgr := secretsmanager.New(session.New())
	_, err := secretsMgr.CreateSecretWithContext(ctx, &secretsmanager.CreateSecretInput{
		Name:         &secretName,
		SecretString: aws.String("{}"),
	})

	return err
}

func Set(ctx context.Context, secretName, key, value string) error {
	values, err := Get(ctx, secretName)
	if err != nil {
		return err
	}

	values[key] = value
	return update(ctx, secretName, values)
}

func Unset(ctx context.Context, secretName, key string) error {
	values, err := Get(ctx, secretName)
	if err != nil {
		return err
	}

	delete(values, key)
	return update(ctx, secretName, values)
}

func update(ctx context.Context, id string, values map[string]string) error {
	secretValue, err := json.Marshal(values)
	if err != nil {
		return err
	}

	secretsMgr := secretsmanager.New(session.New())
	_, err = secretsMgr.UpdateSecretWithContext(ctx, &secretsmanager.UpdateSecretInput{
		SecretId:     &id,
		SecretString: aws.String(string(secretValue)),
	})

	return err
}
