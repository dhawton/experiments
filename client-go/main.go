package main

import (
	"context"
	"fmt"
	"os"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	certmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	"github.com/joho/godotenv"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var namespace = "default"

func main() {
	godotenv.Load()

	delete()
}

func delete() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	certClient, err := cmclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	cert := certClient.CertmanagerV1().Certificates(namespace)
	crt := &certv1.Certificate{
		TypeMeta: meta.TypeMeta{
			Kind:       "Certificate",
			APIVersion: "certmanager.k8s.io/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "client-tls-kzdv-io",
			Namespace: namespace,
		},
	}

	err = cert.Delete(context.TODO(), crt.Name, meta.DeleteOptions{})
	if err != nil {
		//panic(err)
		fmt.Println("Certificate does not exist")
	}

	err = client.CoreV1().Secrets(namespace).Delete(context.TODO(), "client-tls-kzdv-io", meta.DeleteOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Hopefully the certificate is deleted")
}

func issue() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}

	certClient, err := cmclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	cert := certClient.CertmanagerV1().Certificates("")
	crt := &certv1.Certificate{
		TypeMeta: meta.TypeMeta{
			Kind:       "Certificate",
			APIVersion: "certmanager.k8s.io/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "client-tls-kzdv-io",
			Namespace: namespace,
		},
		Spec: certv1.CertificateSpec{
			IssuerRef: certmeta.ObjectReference{
				Name: "addysnip-http01-issuer",
				Kind: "ClusterIssuer",
			},
			SecretName: "client-tls-kzdv-io",
			DNSNames:   []string{"kzdv.io"},
		},
	}

	_, err = cert.Create(context.TODO(), crt, meta.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Hopefully the certificate is applied")
}
