package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
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

	fmt.Println(len(os.Args))
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
	fmt.Println(splitImageName[0])
	//fmt.Println(splitImageName[1])

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
	defer reTagImage(imageFullname, imageTag, newRepo, splitImageName[0])

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

	//ImagePush(newImageFullName)
}

//Function to push retagged image to new registry
/*func ImagePush(newImageFullName string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	auth := types.AuthConfig{
		Username: cfg.User,
		Password: cfg.Passwd,
	}
	authBytes, _ := json.Marshal(auth)
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)

	out, err := cli.ImagePush(ctx, newImageFullName, types.ImagePushOptions{})
	if err != nil {
		panic(err)
	}
	out.Close()


	//io.ReadCloser()
}
*/
//(ctx context.Context, image string, options types.ImagePushOptions) (io.ReadCloser, error)
