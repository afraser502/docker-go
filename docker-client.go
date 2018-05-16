package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

//Main function used to pull image based upon command line args
//Arg1 = repo name
//Arg2 = tag
func main() {
	//Set main variables
	newRepo := "afraser502/" //new-repo name for retagging

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	imageName := ""
	imageTag := ""

	//fmt.Println(len(os.Args))
	if len(os.Args) <= 1 {
		fmt.Println("You must enter the name of an image to download")
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		imageName = strings.ToLower(os.Args[1])
	}
	if len(os.Args) > 2 {
		imageTag = os.Args[2]
	}

	splitImageName := strings.Split(imageName, "/")
	//fmt.Println(len(splitImageName))

	intSlice := splitImageName
	last := intSlice[len(splitImageName)-1]
	//fmt.Printf("Last element: %v\n", last)

	imageFullname := ""
	if imageTag != "" {
		imageFullname = imageName + `:` + imageTag
	} else {
		imageTag = "latest"
		imageFullname = imageName + `:` + imageTag

	}

	out, err := cli.ImagePull(ctx, imageFullname, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	//If image downloads successfully, retag image
	defer reTagImage(imageFullname, imageTag, newRepo, last)

	defer out.Close()

	io.Copy(os.Stdout, out)

}

//Function to retag image
func reTagImage(imageFullname string, imageTag string, newRepo string, splitImageName string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	newImageFullName := newRepo + splitImageName + `:` + imageTag
	fmt.Println(newImageFullName)
	fmt.Println(imageFullname)
	err1 := cli.ImageTag(ctx, imageFullname, newImageFullName)
	if err1 != nil {
		panic(err1)
	}

	ImagePush(newImageFullName)
}

//ImagePush of retagged image to new registry
func ImagePush(newImageFullName string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()

	authConfig := types.AuthConfig{
		Username: "username",
		Password: "password",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePush(ctx, newImageFullName, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	//Parse the responses of image push
	responses, err := ioutil.ReadAll(out)
	fmt.Println(string(responses))

	defer out.Close()

}
