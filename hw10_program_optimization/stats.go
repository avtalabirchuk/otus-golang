package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain), nil
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)

	for i := 0; scanner.Scan(); i++ {
		var user User
		if err = user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) DomainStat {
	result := make(DomainStat)

	for _, user := range u {
		if strings.Contains(user.Email, "."+domain) && strings.Contains(user.Email, "@") {
			splittedString := strings.SplitN(user.Email, "@", 2)
			if len(splittedString) == 2 {
				result[strings.ToLower(splittedString[1])]++
			}
		}
	}
	return result
}
