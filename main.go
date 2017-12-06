package main

import (
    "fmt"
    "flag"
    "os/exec"
    "os"
    "encoding/json"
    "net/http"
    "io/ioutil"
)

var username = flag.String("username", "", "the user who owns the repositories to be synced")
var rootDir = flag.String("root", ".", "the root directory to sync the repositories to")

func main() {
    flag.Parse()
    fmt.Println("username=" + *username)
    fmt.Println("root=" + *rootDir)

    var repositories = retrieveRepositories(*username)
    if repositories == nil {
        fmt.Printf("Unable to retrieve repositories for user %s, exiting...", username)
    }

    for _, repo := range repositories {
        fmt.Println(repo.Name)
    }
    repoUrl := fmt.Sprintf("https://github.com/%s/MbDotNet.git", *username)
    cmdText := fmt.Sprintf("\nexec: git -C %s clone %s", *rootDir, repoUrl)
    fmt.Println(cmdText)
    cmd := exec.Command("git", "-C", *rootDir, "clone", repoUrl)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Println(err)
    }
}

func retrieveRepositories(username string) []Repo {
    url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("http.NewRequest: ", err)
        return nil 
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("http.Client.Do: ", err)
        return nil
    }

    defer resp.Body.Close()

    var repos []Repo

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("ioutil.ReadAll: ", err)
        return nil
    }

    if err := json.Unmarshal(body, &repos); err != nil {
        fmt.Println(err)
        return nil
    }

    return repos
}

type Repo struct {
    Name string `json:"name"`
}
