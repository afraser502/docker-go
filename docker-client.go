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
	reTagImage(imageFullname, imageTag)

	defer out.Close()

	io.Copy(os.Stdout, out)

}

//Function to retag image
func reTagImage(imageFullname string, imageTag string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	//Set your repo name here
	newRepo := "your-repo"
	err1 := cli.ImageTag(ctx, imageFullname, newRepo+`:`+imageTag)
	if err1 != nil {
		panic(err1)
	}
}
