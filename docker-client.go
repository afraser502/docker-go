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
		imageFullname = imageName
	}

	fmt.Println(imageFullname)

	out, err := cli.ImagePull(ctx, imageFullname, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

	reTagImage(imageFullname, imageTag)
}

func reTagImage(imageFullname string, imageTag string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	out, err := cli.ImageTag(ctx, imageFullname, "test"+`:`+imageTag)
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

}
