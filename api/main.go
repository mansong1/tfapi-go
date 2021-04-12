package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	tf_core_framework "tensorflow/core/framework/tensor_go_proto"
	tf_core_framework_shape "tensorflow/core/framework/tensor_shape_go_proto"
	tf_core_framework_types "tensorflow/core/framework/types_go_proto"

	pb "tensorflow_serving/apis"
)

type Config struct {
	Server struct {
		Port string `envconfig:"SERVER_PORT"`
		Host string `envconfig:"SERVER_HOST"`
	}
}

type payload struct {
	URL string `json:"URL"`
}

type ClassifyResult struct {
	Label string `json:"label"`
}

func main() {
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Homepage!"))
}

func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/classify", classifyHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(router)))
}

func classifyHandler(w http.ResponseWriter, r *http.Request) {

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	var tfServer string = cfg.Server.Host + ":" + cfg.Server.Port

	var imgURL payload
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Didn't receive image url: %s,", err)
	}

	json.Unmarshal(reqBody, &imgURL)
	w.WriteHeader(http.StatusCreated)

	resp, err := http.Get(imgURL.URL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	imageBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	request := &pb.PredictRequest{
		ModelSpec: &pb.ModelSpec{
			Name:          "resnet",
			SignatureName: "serving_default",
		},
		Inputs: map[string]*tf_core_framework.TensorProto{
			"image_bytes": &tf_core_framework.TensorProto{
				Dtype: tf_core_framework_types.DataType_DT_STRING,
				TensorShape: &tf_core_framework_shape.TensorShapeProto{
					Dim: []*tf_core_framework_shape.TensorShapeProto_Dim{
						&tf_core_framework_shape.TensorShapeProto_Dim{
							Size: int64(1),
						},
					},
				},
				StringVal: [][]byte{imageBytes},
			},
		},
	}

	// Connect to server
	conn, err := grpc.Dial(tfServer, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Cannot connect to server %v", err)
		responseError(w, "Cannot connect to tfSever", http.StatusInternalServerError)
	}
	defer conn.Close()

	stub := pb.NewPredictionServiceClient(conn)

	result, err := stub.Predict(context.Background(), request)
	if err != nil {
		responseError(w, "Could not run prediction", http.StatusInternalServerError)
	}

	resultClass := result.Outputs["classes"].Int64Val[0]
	predictidx := int(resultClass) - 1 // return int

	// get the classes
	classes, err := getClassName()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	log.Printf("Classified Image as %v", classes[predictidx][1])
	responseJSON(w, ClassifyResult{Label: classes[predictidx][1]})
}
