package students

import (
	th "students/testhelpers"
	"testing"
)

func TestUserCRUD(t *testing.T) {

	newUserOpts := CreateNewUserOptions{
		UserName:       "testuser",
		Email:          "testuser@example1.com",
		HashedPassword: "password123",
	}
	newUser, err := newUserOpts.CreateNewUser()
	if err != nil {
		t.Fatalf("Unexpected error creating user: %v", err)
	}

	user, err := GetUser(newUser.OwnerID)
	if err != nil {
		t.Fatalf("Unexpected error getting user: %v", err)
	}

	th.AssertEqual(t, "owner ID", newUser.OwnerID, user.OwnerID)
	th.AssertEqual(t, "user_name", newUser.UserName, user.UserName)
	th.AssertEqual(t, "email", newUser.Email, user.Email)
	//th.AssertEqual(t, "hashed password", newUser.HashedPassword, user.HashedPassword)
	// Update the user's email
	updateUserOpts := UpdateUserOptions{
		OwnerID:        newUser.OwnerID,
		UserName:       "eeeeeeeeeeeels",
		Email:          "newemail@example.com",
		HashedPassword: "2muchcodesendhelp",
	}
	err = updateUserOpts.UpdateUser()
	if err != nil {
		t.Fatalf("Unexpected error updating user: %v", err)
	}

	// Get the updated user
	updatedUser, err := GetUser(newUser.OwnerID)
	if err != nil {
		t.Fatalf("Unexpected error getting user: %v", err)
	}
	th.AssertEqual(t, "user_name", updatedUser.UserName, updateUserOpts.UserName)
	th.AssertEqual(t, "email", updatedUser.Email, updateUserOpts.Email)
	th.AssertEqual(t, "hashed password", updatedUser.HashedPassword, updateUserOpts.HashedPassword)

	// Delete the user
	err = DeleteUser(newUser.OwnerID)
	if err != nil {
		t.Fatalf("Unexpected error deleting user: %v", err)
	}

	// Get the deleted user
	deletedUser, err := GetUser(newUser.OwnerID)
	if err == nil {
		t.Fatalf("Expected error getting deleted user, but got %v", deletedUser)
	}
}
