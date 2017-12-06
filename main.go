package main

import (
    "fmt"
    "flag"
    "os/exec"
    "os"
)

var username = flag.String("username", "", "the user who owns the repositories to be synced")
var rootDir = flag.String("root", ".", "the root directory to sync the repositories to")

func main() {
    flag.Parse()
    fmt.Println("username=" + *username)
    fmt.Println("root=" + *rootDir)

    repoUrl := "https://github.com/" + *username + "/MbDotNet.git"
    cmdText := "git -C " + *rootDir + " clone " + repoUrl
    fmt.Println(cmdText)
    cmd := exec.Command("git", "-C", *rootDir, "clone", repoUrl)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Println(err)
    }
}
