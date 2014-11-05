// Copyright 2014 li. All rights reserved.

package mux

import (
	"fmt"
	"testing"
)

func Example_route() {

	router := New("simpleRouter")
	router.Add([]string{"GET", "POST"}, `/home/profile(id:^[1-9]*$)`)
	err := router.Start()
	if err != nil {
		//handle
		return
	}

	//route for url
	result := router.Route("POST", "/home/profile1")

	if !result.IsMatch {
		//handle
		return
	}

	fmt.Println(result.IsMatch) // should be true
	fmt.Println(result.Url)     // should be `/home/profile(id:^[1-9]*$)`
	if result != nil {
		//handle params
		fmt.Println(result.Parse()) // should equals map[id:[1]]
	}
}

func TestRouterRoute(t *testing.T) {
	router := New("myRouter")
	router.Add([]string{"GET", "POST"}, "/")
	router.Add([]string{"GET"}, "/home")
	router.Add([]string{"POST"}, "/home/")
	router.Add([]string{"GET", "POST"}, "/home/profile1")
	router.Add([]string{"GET", "POST"}, `/home/profile(id:^[1-9]*$)`)
	router.Add([]string{"GET", "POST"}, `/home/profile(id)`)
	router.Add([]string{"GET", "POST"}, `/home/(all:^[a-z]*$)`)
	router.Add([]string{"GET", "POST"}, "/home/(all)")
	router.Add([]string{"GET", "POST"}, "/home/profile1/view")

	err := router.Start()
	if err != nil {
		t.FailNow()
	}

	result1 := router.Route("GET", "/")
	assertTrue(result1.Url == "/", "case1", t)
	result2 := router.Route("POST", "/")
	assertTrue(result2.Url == "/", "case2", t)
	result3 := router.Route("HEAD", "/")
	assertFalse(result3.Url == "/", "case3", t)
	result4 := router.Route("POST", "/profile")
	assertTrue(result4.Url == "/", "case4", t)
	result5 := router.Route("GET", "/profile")
	assertTrue(result5.Url == "/", "case5", t)
	result6 := router.Route("GET", "")
	assertTrue(result6.Url == "/", "case6", t)

	// /home  /home/
	result7 := router.Route("GET", "/home")
	assertTrue(result7.Url == "/home", "case7", t)
	result8 := router.Route("POST", "/home")
	assertTrue(result8.Url == "/home/", "case8", t)
	result9 := router.Route("POST", "/home/abc/xyz")
	assertTrue(result9.Url == `/home/(all:^[a-z]*$)`, "case9", t)

	// /home/profile1
	result10 := router.Route("POST", "/home/profile1")
	assertTrue(result10.Url == "/home/profile1", "case10", t)
	result11 := router.Route("POST", "/home/profile1")
	assertTrue(result11.Url == "/home/profile1", "case11", t)

	// /home/profile(id:^[1-9]*$)

	result12 := router.Route("POST", "/home/profile123")
	assertTrue(result12.Url == `/home/profile(id:^[1-9]*$)`, "case12", t)
	arr12, _ := result12.Parse()["id"]
	assertTrue(len(arr12) == 1 && arr12[0] == "123", "case12", t)

	// /home/profile(id)
	result13 := router.Route("POST", "/home/profileabc")
	assertTrue(result13.Url == `/home/profile(id)`, "case13", t)
	arr13, _ := result13.Parse()["id"]
	assertTrue(len(arr13) == 1 && arr13[0] == "abc", "case13", t)

	// /home/(all:^[a-z]*$)
	result14 := router.Route("POST", "/home/hello")
	assertTrue(result14.Url == `/home/(all:^[a-z]*$)`, "case14", t)
	arr14, _ := result14.Parse()["all"]
	assertTrue(len(arr14) == 1 && arr14[0] == "hello", "case14", t)

	// /home/(all)
	result15 := router.Route("POST", "/home/hello123")
	assertTrue(result15.Url == `/home/(all)`, "case15", t)
	arr15, _ := result15.Parse()["all"]
	assertTrue(len(arr15) == 1 && arr15[0] == "hello123", "case15", t)

	// /home/profile1/view
	result16 := router.Route("POST", "/home/profile1/view/")
	assertTrue(result16.Url == `/home/profile1/view`, "case16", t)
	assertTrue(result16.Parse() == nil, "case16", t)
}
