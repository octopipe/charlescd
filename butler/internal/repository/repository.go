package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
	apiv1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reposytory interface {
	Clone() error
}

type repository struct {
	module charlescdiov1alpha1.Module
	client.Client
}

func NewRepository(client client.Client, module charlescdiov1alpha1.Module) Reposytory {
	return repository{
		module: module,
	}
}

func (r repository) getSecretByModule(secretRef string) (apiv1.Secret, error) {
	secretObjectKey := utils.GetObjectKeyByPath(secretRef)
	secret := apiv1.Secret{}
	err := r.Get(context.Background(), secretObjectKey, &secret)
	if err != nil {
		return apiv1.Secret{}, err
	}

	return secret, nil
}

func (r *repository) getAuthMethodBySecret(secret apiv1.Secret) (transport.AuthMethod, error) {
	if len(secret.Data["sshPrivateKey"]) > 0 {
		return ssh.NewPublicKeys("git", secret.Data["sshPrivateKey"], "")
	}

	if len(secret.Data["username"]) > 0 && len(secret.Data["password"]) > 0 {
		authMethod := &http.BasicAuth{
			Username: string(secret.Data["username"]),
			Password: string(secret.Data["password"]),
		}

		return authMethod, nil
	}

	if len(secret.Data["username"]) > 0 && len(secret.Data["accessToken"]) > 0 {
		authMethod := &http.BasicAuth{
			Username: string(secret.Data["username"]),
			Password: string(secret.Data["accessToken"]),
		}

		return authMethod, nil
	}

	return nil, errors.New("repository auth method is not valid")
}

func (r repository) Clone() error {
	gitCloneConfig := &git.CloneOptions{
		URL:  r.module.Spec.RepositoryPath,
		Auth: nil,
	}

	if r.module.Spec.SecretRef != nil {
		secret, err := r.getSecretByModule(*r.module.Spec.SecretRef)
		if err != nil {
			return err
		}
		authMethod, err := r.getAuthMethodBySecret(secret)
		if err != nil {
			return err
		}
		gitCloneConfig.Auth = authMethod
	}

	_, err := git.PlainClone(fmt.Sprintf("%s/%s", os.Getenv("REPOSITORIES_TMP_DIR"), r.module.Spec.RepositoryPath), false, gitCloneConfig)
	if err != nil && !strings.Contains(err.Error(), "repository already exists") {
		return err
	}

	return nil
}
