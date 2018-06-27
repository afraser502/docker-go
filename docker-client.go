package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
	"gopkg.in/ldap.v2"
)

//Main function used to pull image based upon command line args
//Arg1 = repo name
//Arg2 = tag
func main() {
	//Set main variables
	viper.SetConfigFile("./configs/env.json") //Set config file

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	newRepo := viper.GetString("repo.new") //new-repo name for retagging

	//Search LDAP
	search()

	checkTwistcliExists()

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	//cli, err := client.

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

//LDAP Functions
func search() {

	bindUser := viper.GetString("ldap.username") //bind username
	bindPass := viper.GetString("ldap.password") //bind password
	ldapHost := viper.GetString("ldap.host")     // ldap host address
	ldapPort := viper.GetInt("ldap.port")        // ldap port
	ldapBaseDN := viper.GetString("ldap.baseDN") // baseDN

	//fmt.Println(ldapHost)

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind(bindUser, bindPass)
	if err != nil {
		log.Fatal(err)
	}

	searchRequest := ldap.NewSearchRequest(
		ldapBaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=*))", // The filter to apply
		[]string{"dn", "cn"}, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}
}

func checkTwistcliExists() {
	path, err := exec.LookPath("twistcli")
	if err != nil {
		fmt.Printf("Can't find 'twistcli' executable\n")
	} else {
		fmt.Printf("'twistcli' executable is in '%s'\n", path)
	}
}
