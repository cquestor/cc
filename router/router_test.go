package router_test

import (
	"fmt"
	"testing"

	"github.com/cquestor/cc/router"
)

func TestNode(t *testing.T) {
	router := router.NewCRouter()
	router.AddRoute("GET", "/")
	router.AddRoute("GET", "/chen/*static/nihao")
	router.AddRoute("GET", "/chen/:name")
	router.AddRoute("GET", "/chen/12/nihao")
	router.AddRoute("GET", "/chen/:age/nihao")
	PrintNode(router.Roots["GET"], 0)
	fmt.Println("----------")
	fmt.Println(router.GetRoute("GET", "/chen/12/nihao"))
}

func PrintNode(node *router.CNode, index int) {
	for _, each := range node.Children {
		fmt.Println(index, each.Current)
		PrintNode(each, index+1)
	}
	index++
}
