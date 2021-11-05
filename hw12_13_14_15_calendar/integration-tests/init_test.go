package integrationtests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var user User

var _ = BeforeSuite(func() {
	var statusCode int
	var err error
	statusCode, user, err = CreateUser(UsersURL, User{
		Email:     "test@test.com",
		FirstName: "Tom",
		LastName:  "Smith",
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(statusCode).To(Equal(200))
	Expect(user.ID).NotTo(Equal(0))
})
