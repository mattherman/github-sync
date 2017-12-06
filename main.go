package main

import (
    "fmt"
    "flag"
    "os/exec"
    "os"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "path"
)

var username = flag.String("username", "", "the user who owns the repositories to be synced")
var rootDir = flag.String("root", ".", "the root directory to sync the repositories to")

func main() {
    flag.Parse()
    fmt.Println("username=" + *username)
    fmt.Println("root=" + *rootDir)

    var repositories = retrieveRepositories(*username)
    if repositories == nil {
        fmt.Printf("Unable to retrieve repositories for user \"%s\", exiting...", *username)
        return
    }

    fmt.Printf("\nFound %d repositories for user \"%s\".\n", len(repositories), *username)

    for _, repo := range repositories {
        syncRepository(repo, *rootDir)
    }
}

func syncRepository(repo Repo, rootDir string) {
    repoNameWithExtension := fmt.Sprintf("%s.git", repo.Name)
    repoPath := path.Join(rootDir, repoNameWithExtension)
    if directoryExists(repoPath) {
        fmt.Printf("\nRepository \"%s\" is already mirrored. Fetching latest...", repo.Name)
        fetchRepository(repoPath)
    } else {
        fmt.Printf("\nRepository \"%s\" is not mirrored. Cloning...", repo.Name)
        cloneRepository(repo, rootDir)
    }
}

func fetchRepository(repoPath string) {
    cmdText := fmt.Sprintf("\nexec: git -C %s fetch origin", repoPath)
    fmt.Println(cmdText)

    cmd := exec.Command("git", "-C", repoPath, "fetch", "origin")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err := cmd.Run()
    if err != nil {
        fmt.Println(err)
    }
}

func cloneRepository(repo Repo, rootDir string) {
    repoUrl := fmt.Sprintf("https://github.com/%s", repo.FullName)
    cmdText := fmt.Sprintf("\nexec: git -C %s clone --mirror %s", rootDir, repoUrl)
    fmt.Println(cmdText)

    cmd := exec.Command("git", "-C", rootDir, "clone", "--mirror", repoUrl)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err := cmd.Run()
    if err != nil {
        fmt.Println(err)
    }
}

func directoryExists(directory string) bool {
    _, err := os.Stat(directory)
    return !os.IsNotExist(err)
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
    FullName string `json:"full_name"`
}
