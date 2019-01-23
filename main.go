package main

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"os"
	"os/user"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {
	ctx := context.Background()

	tokenSource := &TokenSource{
		AccessToken: os.Getenv("DO_PAT"),
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	clusters, _, err := client.Kubernetes.List(ctx, &godo.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}

	for _, cluster := range clusters {
		clusterConfig, _, err := client.Kubernetes.GetKubeConfig(ctx, cluster.ID)
		if err != nil {
			panic(err.Error())
		}

		f, err := os.Create(fmt.Sprintf("%s/.kube/%s-kubeconfig.yaml", usr.HomeDir, cluster.Name))
		if err != nil {
			panic(err.Error())
		}

		_, err = f.Write(clusterConfig.KubeconfigYAML)
		if err != nil {
			panic(err.Error())
		}

		err = f.Close()
		if err != nil {
			panic(err.Error())
		}

		_, err = fmt.Fprintf(os.Stdout, "Config for %s written to .kube\n", cluster.Name)
		if err != nil {
			panic(err.Error())
		}
	}
}
