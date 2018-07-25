package main

import (
	"fmt"
	"log"
	"net/http"

	"k8s.io/apimachinery/pkg/api/resource"

	"io/ioutil"

	"encoding/json"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request handling...")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	ar := v1beta1.AdmissionReview{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		payload, err := json.Marshal(&v1beta1.AdmissionResponse{
			UID:     ar.Request.UID,
			Allowed: false,
			Result: &metav1.Status{
				Message: err.Error(),
			},
		})
		if err != nil {
			fmt.Println(err)
		}
		w.Write(payload)
	}

	admitResponse := &v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:     ar.Request.UID,
			Allowed: false,
			// error status by default
			Result: &metav1.Status{
				Message: "Submitted deployment doesn't contain ressource limits",
			},
		},
	}

	if ar.Request.Kind.Kind == "Pod" {
		pod := v1.Pod{}
		json.Unmarshal(ar.Request.Object.Raw, &pod)
		// look over all the containers: if a container has a resource limitation -
		// break the loop with OK status (you didn't precise all the conditions
		// so I've implemented the simplest case, but it can be easily modified)
		for _, container := range pod.Spec.Containers {
			if *container.Resources.Limits.Cpu() != (resource.Quantity{}) ||
				*container.Resources.Limits.Memory() != (resource.Quantity{}) {
				admitResponse.Response.Allowed = true
				admitResponse.Response.Result = &metav1.Status{
					Message: "OK",
				}
				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	payload, err := json.Marshal(admitResponse)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(payload)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
