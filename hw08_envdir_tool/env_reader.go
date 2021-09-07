package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]string

// EnvValue helps to distinguish between empty files and files with the first empty line.
// type EnvValue struct {
// 	Value      string
// 	NeedRemove bool
// }

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func readFirstLineFile(file *os.File) (string, error) {
	reader := bufio.NewReader(file)
	value, _ := reader.ReadString('\n')

	return value, nil
}

func processValue(value string) string {
	value = strings.TrimRight(value, " \t\n")
	value = strings.ReplaceAll(value, "\x00", "\n")
	return value
}

func getValueFromFileDir(dir, fileName string) (string, error) {
	filePath := path.Join(dir, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	value, _ := readFirstLineFile(file)
	value = processValue(value)

	return value, nil
}
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	fileNames, _ := ioutil.ReadDir(dir)

	env := make(Environment, len(fileNames))
	for _, fileName := range fileNames {
		value, _ := getValueFromFileDir(dir, fileName.Name())
		env[fileName.Name()] = value
	}
	return env, nil

}
