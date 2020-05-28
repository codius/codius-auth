package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deductBalance(id *string) error {
	url := fmt.Sprintf("%s/balances/%s:spend", os.Getenv("RECEIPT_VERIFIER_URL"), *id)
	log.Println(url)
	resp, err := http.Post(url, "text/plain", bytes.NewBuffer([]byte(os.Getenv("AUTH_PRICE"))))
	if err != nil {
		fmt.Println("Balance spend error:", err)
		return err
	}
	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println("Balance spend error:", string(b))
		return errors.New(string(b))
	}
	fmt.Println("Balance:", string(b))
	return nil
}

func forwardAuth(rw http.ResponseWriter, req *http.Request) {
	serviceHost := req.Header.Get("x-forwarded-host")
	if serviceHost == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	serviceId := strings.SplitN(serviceHost, ".", 2)[0]
	err := deductBalance(&serviceId)
	if err != nil {
		url402 := fmt.Sprintf("%s/%s/402", os.Getenv("CODIUS_HOST_URL"), serviceId)
		http.Redirect(rw, req, url402, http.StatusSeeOther)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

func tokenAuth(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.NotFound(rw, req)
		return
	}
	body := json.NewDecoder(req.Body)
	tr := &authv1.TokenReview{}
	err := body.Decode(tr)
	if err != nil {
		handleErr(rw, err)
		return
	}

    err = deductBalance(&tr.Spec.Token)
	if err != nil {
		handleErr(rw, err)
		return
	}

	user := os.Getenv("RBAC_USER")
	// groups := []string{
	// 	"testgroup",
	// }
	trResp := &authv1.TokenReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "authentication.k8s.io/v1",
			Kind:       "TokenReview",
		},
		Status: authv1.TokenReviewStatus{
			Authenticated: true,
			User: authv1.UserInfo{
				UID:      user,
				Username: user,
				// Groups:   groups,
			},
		},
	}

	writeResp(rw, trResp)
}

func writeResp(rw http.ResponseWriter, tr *authv1.TokenReview) {
	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	err := enc.Encode(tr)
	if err != nil {
		log.Println("Failed to encode token review response")
	}
}

func handleErr(rw http.ResponseWriter, err error) {
	writeResp(rw, &authv1.TokenReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "authentication.k8s.io/v1",
			Kind:       "TokenReview",
		},
		Status: authv1.TokenReviewStatus{
			Error: err.Error(),
		},
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port) 
	http.HandleFunc("/forward", forwardAuth)
	http.HandleFunc("/token", tokenAuth)
	fmt.Println("Starting server on", addr)
	http.ListenAndServe(addr, nil)
}
