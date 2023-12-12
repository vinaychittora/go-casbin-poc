package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
)


func main() {
	// create casbin enforcer with policy in mysql and model in conf
	a, err := xormadapter.NewAdapter("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
	    log.Fatalf("error: adapter: %s", err)
	}
	e, err := casbin.NewEnforcer("model.conf", a)
	if err != nil {
		log.Fatalf("unable to create Casbin enforcer: %v", err)
	}

	// Load the policy from DB.
	e.LoadPolicy()

	// set some policy, if it already exists, it wouldn't duplicate.
	_, err = e.AddPolicy("role:admin", "template", "view")
	_, err = e.AddPolicy("role:admin", "template", "edit")
	_, err = e.AddPolicy("role:admin", "template", "rs") // create rs
	_, err = e.AddPolicy("role:admin", "template", "pl") // public links
	_, err = e.AddPolicy("role:accessor", "template", "view")
	_, err = e.AddPolicy("role:accessor", "template", "rs")
	_, err = e.AddPolicy("role:accessor", "template", "pl")
	_, err = e.AddRoleForUser("user1", "role:admin")
	_, err = e.AddRoleForUser("user2", "role:accessor")	


	


	aliceRoles, err := e.GetRolesForUser("user1")
	if err != nil {
		log.Fatalf("could not get roles for user1: %v", err)
	}
	alicePerms, err := e.GetImplicitPermissionsForUser("user1")
	if err != nil {
		log.Fatalf("could not get permissions for user1: %v", err)
	}
	fmt.Printf(
		"user1 is a member of the following roles: %v, and her permissions are: %v\n",
		aliceRoles,
		alicePerms,
	)

	bobRoles, err := e.GetRolesForUser("user2")
	if err != nil {
		log.Fatalf("could not get roles for user2: %v", err)
	}
	bobPerms, err := e.GetImplicitPermissionsForUser("user2")
	if err != nil {
		log.Fatalf("could not get permissions for user2: %v", err)
	}
	fmt.Printf(
		"user2 is a member of the following roles: %v, and his permissions are: %v\n",
		bobRoles,
		bobPerms,
	)

	// Save the policy back to DB.
	e.SavePolicy()
}
