package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"q/rbac-slicer/kubevar"
)

// make some reports showing the general state of roles

func main() {

	var kubeconfig kubevar.Kubeconfig
	flag.Var(&kubeconfig, "kubeconfig", "the path to the kubeconfig")
	labelSelector := flag.String("l", "", "a label selector to filter the results")
	sortByCount := flag.Bool("c", false, "sort the results by count")

	flag.Parse()

	if kubeconfig.String() == "" {
		if err := kubeconfig.Set(kubevar.Default()); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	fmt.Printf("Kubeconfig: %v\n", kubeconfig.String())

	cs := kubeconfig.Clientset
	roles, err := cs.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{LabelSelector: *labelSelector})
	if err != nil {
		panic(err.Error())
	}

	resources := map[string][][]string{}
	verbs := map[string][][]string{}

	for _, role := range roles.Items {
		for _, rule := range role.Rules {

			apigroups := []string{}
			for _, grp := range rule.APIGroups {
				if grp == "" {
					apigroups = append(apigroups, "core")
					continue
				}
				apigroups = append(apigroups, grp)
			}
			apigroupsKey := strings.Join(apigroups, ",")

			if _, found := resources[apigroupsKey]; found {
				resources[apigroupsKey] = append(resources[apigroupsKey], rule.Resources)
			} else {
				resources[apigroupsKey] = [][]string{}
				resources[apigroupsKey] = append(resources[apigroupsKey], rule.Resources)
			}

			if _, found := verbs[apigroupsKey]; found {
				verbs[apigroupsKey] = append(verbs[apigroupsKey], rule.Verbs)
			} else {
				verbs[apigroupsKey] = [][]string{}
				verbs[apigroupsKey] = append(verbs[apigroupsKey], rule.Verbs)
			}
		}
	}

	resourcesCounter := countThem(resources)
	verbsCounter := countThem(verbs)

	ranker := sortKeys
	if *sortByCount {
		ranker = rankByWordCount
	}

	fmt.Println("Resources:")
	for _, k := range ranker(resourcesCounter) {
		fmt.Println(k, resourcesCounter[k])
	}

	fmt.Println("")
	fmt.Println("Verbs:")
	for _, k := range ranker(verbsCounter) {
		fmt.Println(k, verbsCounter[k])
	}

	return
}

func countThem(dat map[string][][]string) map[string]int {
	rv := map[string]int{}
	for k, v := range dat {
		for _, v := range v {
			key := fmt.Sprintf("%s %s", k, strings.Join(v, ","))
			rv[key] = rv[key] + 1
		}
	}

	return rv
}
